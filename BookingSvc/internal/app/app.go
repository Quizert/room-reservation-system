package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/grpc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/controller"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	paymentClient "github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/http/paymentsvc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/kafka"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/config"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/service"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/storage/postgres"

	"github.com/Quizert/room-reservation-system/Libs/metrics"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	mainServer     *http.Server
	metricServer   *http.Server
	dbPool         *pgxpool.Pool
	tracerProvider *trace.TracerProvider // TracerProvider для управления жизненным циклом
	log            *zap.Logger
}

func NewApp() *App {
	return &App{}
}

func NewDatabasePool(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	logger.Info("Connecting to database", zap.String("connection_string", connString))
	return pgxpool.Connect(ctx, connString)
}

func NewKafkaProducer(cfg *config.Config, logger *zap.Logger) *kafka.Producer {
	logger.Info("Initializing Kafka producer", zap.String("broker", cfg.KafkaBroker))
	return kafka.NewProducer([]string{cfg.KafkaBroker}, cfg.KafkaTopicClient, cfg.KafkaTopicHotel)
}

func NewAuthClient(cfg *config.Config, logger *zap.Logger) (*grpc.AuthSvcClient, error) {
	logger.Info("Initializing Auth service client", zap.String("host", cfg.GRPCAuthHost), zap.String("port", cfg.GRPCAuthPort))
	return grpc.NewAuthClient(cfg.GRPCAuthHost, cfg.GRPCAuthPort)
}

func NewHotelClient(cfg *config.Config, logger *zap.Logger) (*grpc.HotelSvcClient, error) {
	logger.Info("Initializing Hotel service client", zap.String("host", cfg.GRPCHotelHost), zap.String("port", cfg.GRPCHotelPort))
	return grpc.NewHotelClient(cfg.GRPCHotelHost, cfg.GRPCHotelPort)
}

func NewPaymentClient(cfg *config.Config, logger *zap.Logger) *paymentClient.Client {
	logger.Info("Initializing Payment service client", zap.String("url", cfg.PaymentSvcURL))
	return paymentClient.NewPaymentSvcClient(cfg.PaymentSvcURL)
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

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}
func (a *App) Init(ctx context.Context) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("myerror initializing zap logger: %v", err)
	}
	a.log = logger

	a.log.Info("Loading configuration")
	// Тут можно сделать MustLoad ля-ля
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.log.Info("Initializing dependencies")
	kafkaProducer := NewKafkaProducer(cfg, a.log)

	hotelClient, err := NewHotelClient(cfg, a.log)
	if err != nil {
		return fmt.Errorf("failed to initialize hotel client: %w", err)
	}

	authClient, err := NewAuthClient(cfg, a.log)
	if err != nil {
		return fmt.Errorf("failed to initialize auth client: %w", err)
	}

	paymentSvcClient := NewPaymentClient(cfg, a.log)

	dbPool, err := NewDatabasePool(ctx, cfg, a.log)
	if err != nil {
		return fmt.Errorf("failed to initialize database pool: %w", err)
	}

	tracerProvider, err := InitTracerProvider("BookingSvc", "http://jaeger:14268/api/traces")
	if err != nil {
		return fmt.Errorf("failed to initialize tracer provider: %w", err)
	}

	a.tracerProvider = tracerProvider
	tracer := a.tracerProvider.Tracer("BookingSvc")
	repo := postgres.NewPostgresRepository(dbPool, tracer)
	a.dbPool = dbPool
	mainService := service.NewBookingServiceImpl(repo, kafkaProducer, hotelClient, authClient, paymentSvcClient, tracer, a.log)
	bookingHandler := controller.NewBookingHandler(mainService, tracer)

	mainRoute := controller.SetupRoutes(bookingHandler)
	metricRoute := metrics.SetupMetricsRoute()
	a.mainServer = &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mainRoute,
	}
	a.metricServer = &http.Server{
		Addr:    ":" + cfg.HTTPMetricPort,
		Handler: metricRoute,
	}

	a.log.Debug("Initialization complete")
	return nil
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info("Starting HTTP mainServer")

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
		if err := a.mainServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("Error in ListenAndServe", zap.Error(err))
			return fmt.Errorf("failed to serve HTTP mainServer: %w", err)
		}
		a.log.Info("HTTP mainServer stopped")
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
	if a.mainServer != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		if err := a.mainServer.Shutdown(shutdownCtx); err != nil {
			a.log.Error("HTTP mainServer shutdown error", zap.Error(err))
			return fmt.Errorf("failed to shutdown HTTP mainServer: %w", err)
		}
	}

	if a.tracerProvider != nil {
		if err := a.tracerProvider.Shutdown(ctx); err != nil {
			a.log.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}
	return nil
}
