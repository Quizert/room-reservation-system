package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/config"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/controller"
	grpcserver "github.com/Quizert/room-reservation-system/AuthSvc/internal/controller/grpc"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/service"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/storage/postgres"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"github.com/Quizert/room-reservation-system/Libs/metrics"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	server       *http.Server
	metricServer *http.Server
	GRPCServer   *grpcserver.Server
	dbPool       *pgxpool.Pool
	log          *zap.Logger
}

func NewApp() *App {
	return &App{}
}

func NewDatabasePool(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	logger.Info("Connecting to database", zap.String("connection_string", connString))
	return pgxpool.Connect(ctx, connString)
}

func (a *App) ListenGRPCServer() error {
	grpcServer := grpc.NewServer()

	authpb.RegisterAuthServiceServer(grpcServer, a.GRPCServer)

	lis, err := net.Listen("tcp", a.GRPCServer.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return grpcServer.Serve(lis)
}

func (a *App) Init(ctx context.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// Тут можно сделать MustLoad ля-ля
		return fmt.Errorf("myerror initializing zap logger: %v", err)
	}
	a.log = logger

	a.log.Info("Loading configuration")
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("myerror loading config: %v", err)
	}

	dbPool, err := NewDatabasePool(ctx, cfg, a.log)
	if err != nil {
		return fmt.Errorf("failed to initialize database pool: %w", err)
	}
	repo := postgres.NewPostgresRepository(dbPool)
	a.dbPool = dbPool

	tokenTTLString := cfg.TokenTTl
	tokenTTL, err := time.ParseDuration(tokenTTLString)
	if err != nil {
		return fmt.Errorf("myerror parsing duration: %w", err)

	}
	secret := cfg.Secret

	authService := service.NewAuthServiceImpl(repo, tokenTTL, secret, logger)

	authHandler := controller.NewAuthHandler(authService)
	route := controller.SetupRoutes(authHandler)
	metricRoute := metrics.SetupMetricsRoute()

	a.server = &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: route,
	}
	a.metricServer = &http.Server{
		Addr:    ":" + cfg.HTTPMetricPort,
		Handler: metricRoute,
	}
	a.GRPCServer = grpcserver.NewServer(authService, ":"+cfg.GRPCPort)
	a.log.Debug("Initialization complete")

	return nil
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info("Starting HTTP server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if err := a.metricServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("Error in ListenAndServe metricServer", zap.Error(err))
			return fmt.Errorf("failed to serve HTTP metricServer: %w", err)
		}
		a.log.Info("HTTP mainServer stopped")
		return nil
	})

	group.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("Error in ListenAndServe", zap.Error(err))
			return fmt.Errorf("failed to serve HTTP server: %w", err)
		}
		a.log.Info("HTTP server stopped")
		return nil
	})

	group.Go(func() error {
		if err := a.ListenGRPCServer(); err != nil {
			a.log.Error("Error in ListenAndServe", zap.Error(err))
			return fmt.Errorf("failed to serve GRPC server: %w", err)
		}
		a.log.Info("GRPC server stopped")
		return nil
	})

	group.Go(func() error {
		<-groupCtx.Done()
		return a.Stop(context.Background())
	})

	if err := group.Wait(); err != nil {
		a.log.Error("Error after wait", zap.Error(err))
		return err
	}
	a.log.Info("Server shutdown gracefully")
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	a.log.Info("Shutting down HTTP server")
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.log.Error("HTTP server shutdown myerror", zap.Error(err))
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	a.log.Info("HTTP server shutdown gracefully")

	if a.dbPool != nil {
		a.dbPool.Close()
		a.log.Info("Database connection closed")
	} else {
		a.log.Warn("Database pool is nil, skipping close")
	}
	return nil
}
