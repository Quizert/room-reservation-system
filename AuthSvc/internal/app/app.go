package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/config"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/controller"
	grpcserver "github.com/Quizert/room-reservation-system/AuthSvc/internal/controller/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"

	"github.com/Quizert/room-reservation-system/AuthSvc/internal/service"
	"github.com/Quizert/room-reservation-system/AuthSvc/internal/storage/postgres"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
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
	server         *http.Server
	GRPCServer     *grpcserver.Server
	dbPool         *pgxpool.Pool
	log            *zap.Logger
	tracerProvider *trace.TracerProvider // (1) Храним TracerProvider здесь
}

func NewApp() *App {
	return &App{}
}

func NewDatabasePool(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	logger.Info("Connecting to database", zap.String("connection_string", connString))
	return pgxpool.Connect(ctx, connString)
}

func InitTracerProvider(serviceName, endpoint string) (*trace.TracerProvider, error) {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	fmt.Println("Jaeger exporter initialized successfully")

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Устанавливаем провайдер глобально
	otel.SetTracerProvider(tp)

	// Устанавливаем пропагатор контекста
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

func (a *App) ListenGRPCServer() error {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

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
		return fmt.Errorf("error initializing zap logger: %v", err)
	}
	a.log = logger

	a.log.Info("Loading configuration")
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	dbPool, err := NewDatabasePool(ctx, cfg, a.log)
	if err != nil {
		return fmt.Errorf("failed to initialize database pool: %w", err)
	}
	a.dbPool = dbPool

	tokenTTL, err := time.ParseDuration(cfg.TokenTTl)
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}

	// (1) Инициализируем Jaeger-трейсинг и сохраняем в a.tracerProvider
	tp, err := InitTracerProvider("AuthSvc", "http://jaeger:14268/api/traces")
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	a.tracerProvider = tp // сохраняем, чтобы закрыть позже

	tracer := a.tracerProvider.Tracer("AuthSvc")
	authService := service.NewAuthServiceImpl(
		postgres.NewPostgresRepository(dbPool, tracer),
		tokenTTL,
		cfg.Secret,
		tracer,
		logger,
	)
	authHandler := controller.NewAuthHandler(authService, tracer)
	route := controller.SetupRoutes(authHandler)

	a.server = &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: route,
	}

	// gRPC сервер
	a.GRPCServer = grpcserver.NewServer(authService, ":"+cfg.GRPCPort, tracer)

	return nil
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info("Starting HTTP server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	group, groupCtx := errgroup.WithContext(ctx)

	// Запуск HTTP
	group.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("Error in ListenAndServe (HTTP)", zap.Error(err))
			return fmt.Errorf("failed to serve HTTP server: %w", err)
		}
		a.log.Info("HTTP server stopped")
		return nil
	})

	// Запуск gRPC
	group.Go(func() error {
		if err := a.ListenGRPCServer(); err != nil {
			a.log.Error("Error in ListenAndServe (gRPC)", zap.Error(err))
			return fmt.Errorf("failed to serve GRPC server: %w", err)
		}
		a.log.Info("GRPC server stopped")
		return nil
	})

	// Ожидаем сигналов
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
		a.log.Error("HTTP server shutdown error", zap.Error(err))
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	a.log.Info("HTTP server shutdown gracefully")

	// Закрываем соединение с БД
	if a.dbPool != nil {
		a.dbPool.Close()
		a.log.Info("Database connection closed")
	}

	// (1) Останавливаем tracer provider
	if a.tracerProvider != nil {
		if err := a.tracerProvider.Shutdown(ctx); err != nil {
			a.log.Error("Error shutting down tracer provider", zap.Error(err))
		}
	}

	return nil
}
