package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/mati/go-ticket/docs"
	"github.com/mati/go-ticket/internal/api"
	"github.com/mati/go-ticket/internal/api/middleware"
	"github.com/mati/go-ticket/internal/auth"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/mati/go-ticket/internal/event_handler"
	"github.com/mati/go-ticket/internal/kafka"
	"github.com/mati/go-ticket/internal/postgres"
	"github.com/mati/go-ticket/internal/rabbitmq"
	"github.com/mati/go-ticket/internal/ratelimit"
	"github.com/mati/go-ticket/internal/services"
	"github.com/mati/go-ticket/internal/workers"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Go Ticket API
// @version 1.0
// @description Simple API for booking tickets.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger := setupLogger()
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	if err := godotenv.Load(".env", ".env.local"); err != nil {
		logger.Warn("Failed to load environment variables", "error", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	secretKey := os.Getenv("JWT_SECRET_KEY")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	if dbUrl == "" {
		return errors.New("DATABASE_URL is not set")
	}
	if secretKey == "" {
		return errors.New("JWT_SECRET_KEY is not set")
	}

	authService, err := auth.NewJWTService(secretKey)
	if err != nil {
		return fmt.Errorf("failed to create JWT service: %w", err)
	}

	initialCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(initialCtx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to create database connection pool: %w", err)
	}
	defer pool.Close()

	logger.Info("Running database migrations...")
	if err := postgres.RunMigrations(dbUrl); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("Database migrations completed")

	// Redis
	redisClient := ratelimit.NewClient()
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("failed to close redis client", "error", err)
		}
	}()
	authLimiter := ratelimit.NewRateLimiter(redisClient, 5, 15*time.Minute)
	apiLimiter := ratelimit.NewRateLimiter(redisClient, 100, 1*time.Minute)
	rateLimitAuth := middleware.RateLimiterMiddleware(authLimiter, middleware.IPKey)
	rateLimitAPI := middleware.RateLimiterMiddleware(apiLimiter, middleware.UserKey)

	// === Repositories ===
	eventRepository, bookingRepository, userRepository := setupRepositories(pool)
	// === Services ===
	bookingService, userService, outboxRepository := setupServices(
		eventRepository,
		bookingRepository,
		userRepository,
		authService,
		pool,
	)
	// === Handlers ===
	eventHandler := api.NewHTTPHandler(eventRepository, bookingRepository, bookingService)
	authHandler := api.NewAuthHandler(userService)

	mux := http.NewServeMux()
	setupRoutes(mux, authService, eventHandler, authHandler, rateLimitAuth, rateLimitAPI)

	erChan := make(chan error, 3)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// RabbitMQ
	connection, rabbitMQPublisher, err := setupRabbitMQ(rabbitMQURL)
	if err != nil {
		return err
	}

	defer func() {
		if err := connection.Close(); err != nil {
			logger.Error("failed to close connection", "error", err)
		}
	}()

	defer func() {
		if err := rabbitMQPublisher.Close(); err != nil {
			logger.Error("failed to close rabbitmq publisher", "error", err)
		}
	}()

	// Producer
	workerCtx, relayCancel := context.WithCancel(context.Background())
	defer relayCancel()
	producer, relay, err := setupKafkaRelay(outboxRepository)
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("failed to close kafka producer", "error", err)
		}
	}()

	go func() {
		err := relay.Start(workerCtx)
		if err != nil {
			erChan <- fmt.Errorf("relay error: %w", err)
		}
	}()

	// RabbitMQ Email
	consumer, worker, errSetUpEmail := setupEmailWorker(connection, redisClient, logger)
	if errSetUpEmail != nil {
		return errSetUpEmail
	}

	go func() {
		err := worker.Start(workerCtx)
		if err != nil {
			erChan <- fmt.Errorf("failed to start worker: %w", err)
		}
	}()

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Error("failed to close consumer", "error", err)
		}
	}()

	// Consumer
	bookingEvent := event_handler.NewBookingEventHandler(logger, rabbitMQPublisher)
	consumerGroup, consumerWorker, err := setupKafkaConsumer(logger, bookingEvent)
	if err != nil {
		return err
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			logger.Error("failed to close consumer group", "error", err)
		}
	}()

	go func() {
		err := consumerWorker.Start(workerCtx)
		if err != nil {
			erChan <- fmt.Errorf("consumer error: %w", err)
		}
	}()
	srv := setupServer(mux)

	go func() {
		if err := runServer(srv); err != nil {
			erChan <- fmt.Errorf("failed to start server %w", err)
		}
	}()

	logger.Info("Server started on :8080")

	select {
	case err := <-erChan:
		logger.Error("Worker failed", "error", err)
	case <-stop:
		logger.Info("Shutdown signal received")
	}

	relayCancel()

	if err := gracefulShutdown(srv, logger); err != nil {
		logger.Error("Shutdown error", "error", err)
	}
	return nil
}

func setupRoutes(
	mux *http.ServeMux,
	authService *auth.JWTService,
	eventHandler *api.HTTPHandler,
	authHandler *api.AuthHandler,
	rateLimitAuth func(http.HandlerFunc) http.HandlerFunc,
	rateLimitAPI func(http.HandlerFunc) http.HandlerFunc,
) {
	auth := func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(authService, handler)
	}

	requireOrganizer := func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.RequireRole([]domain.UserRole{domain.UserRoleOrganizer}, handler)
	}

	requireAdmin := func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.RequireRole([]domain.UserRole{domain.UserRoleAdmin}, handler)
	}

	requireAll := func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.RequireRole(
			[]domain.UserRole{domain.UserRoleUser, domain.UserRoleAdmin, domain.UserRoleOrganizer},
			handler,
		)
	}

	// Swagger
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// === Public endpoints ===
	mux.HandleFunc("POST /auth/register", rateLimitAuth(authHandler.Register))
	mux.HandleFunc("POST /auth/login", rateLimitAuth(authHandler.Login))

	// === Protected endpoints ===
	mux.HandleFunc("POST /events", auth(requireOrganizer(rateLimitAPI(eventHandler.CreateEvent))))
	mux.HandleFunc("PUT /events/{id}", auth(requireOrganizer(rateLimitAPI(eventHandler.UpdateEvent))))
	mux.HandleFunc("DELETE /events/{id}", auth(requireAdmin(rateLimitAPI(eventHandler.DeleteEvent))))
	mux.HandleFunc("GET /events/{id}", auth(requireAll(rateLimitAPI(eventHandler.GetEvent))))
	mux.HandleFunc("GET /events", auth(requireAll(rateLimitAPI(eventHandler.ListEvents))))
	mux.HandleFunc("POST /events/{event_id}/bookings", auth(requireAll(rateLimitAPI(eventHandler.CreateBooking))))
}

func setupServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         ":8080",
		Handler:      middleware.LoggingMiddleware(middleware.RecoveryMiddleware(mux)),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func gracefulShutdown(srv *http.Server, logger *slog.Logger) error {
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	logger.Info("Server shutdown successfully")

	return nil
}

func runServer(srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func setupRepositories(pool *pgxpool.Pool) (
	*postgres.EventRepository,
	*postgres.BookingRepository,
	*postgres.UserRepository,
) {
	queries := postgres.New(pool)
	return postgres.NewEventRepository(queries),
		postgres.NewBookingRepository(queries),
		postgres.NewUserRepository(queries)
}

func setupServices(
	eventRepository *postgres.EventRepository,
	bookingRepository *postgres.BookingRepository,
	userRepository *postgres.UserRepository,
	authService *auth.JWTService,
	pool *pgxpool.Pool,
) (*services.BookingService, *services.UserService, *postgres.OutBoxRepository) {
	transactionManager := postgres.NewPgxTxManager(pool)
	outboxRepository := postgres.NewOutBoxRepository(postgres.New(pool))
	bookingService := services.NewBookingService(eventRepository, bookingRepository, outboxRepository, transactionManager)
	userService := services.NewUserService(userRepository, authService)
	return bookingService, userService, outboxRepository
}

func setupKafkaRelay(outboxRepository *postgres.OutBoxRepository) (sarama.SyncProducer, *workers.OutboxRelay, error) {
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	if kafkaAddr == "" {
		return nil, nil, errors.New("KAFKA_ADDR is not set")
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	kafkaBroker := kafka.NewKafkaBroker(producer)
	relay := workers.NewOutboxRelay(outboxRepository, kafkaBroker)

	return producer, relay, nil
}

func setupKafkaConsumer(
	logger *slog.Logger,
	event domain.EventHandler) (sarama.ConsumerGroup, *workers.KafkaConsumerWorker, error) {
	const (
		topic = "booking_events_topic"
	)
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	if kafkaAddr == "" {
		return nil, nil, errors.New("KAFKA_ADDR is not set")
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup([]string{kafkaAddr}, "bookingEvent-group", config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	kafkaConsumer := kafka.NewKafkaConsumer(consumerGroup, []string{topic}, logger, event)
	kafkaConsumerWorker := workers.NewKafkaConsumerWorker(kafkaConsumer, logger)
	return consumerGroup, kafkaConsumerWorker, nil
}

func setupRabbitMQ(url string) (*amqp.Connection, *rabbitmq.RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	mqPublisher, err := rabbitmq.NewRabbitMQPublisher(conn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, mqPublisher, nil
}

func setupEmailWorker(
	conn *amqp.Connection,
	redisClient *redis.Client,
	logger *slog.Logger) (consumer *rabbitmq.RabbitMQConsumer, worker *workers.EmailWorker, err error) {
	consumer, err = rabbitmq.NewRabbitMqConsumer(conn, "booking.notifications")
	if err != nil {
		return nil, nil, err
	}

	worker = workers.NewEmailWorker(logger, redisClient, consumer)

	return consumer, worker, nil
}
