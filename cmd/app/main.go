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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mati/go-ticket/internal/api"
	"github.com/mati/go-ticket/internal/api/middleware"
	"github.com/mati/go-ticket/internal/auth"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/mati/go-ticket/internal/postgres"
	"github.com/mati/go-ticket/internal/services"
)

func main() {
	logger := setupLogger()
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	if err := godotenv.Load(); err != nil {
		logger.Warn("Failed to load environment variables", "error", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	secretKey := os.Getenv("JWT_SECRET_KEY")

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to create database connection pool: %w", err)
	}
	defer pool.Close()

	logger.Info("Running database migrations...")
	if err := postgres.RunMigrations(dbUrl); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("Database migrations completed")

	// === Repositories ===
	eventRepository, bookingRepository, userRepository := setupRepositories(pool)
	// === Services ===
	bookingService, userService := setupServices(eventRepository, bookingRepository, userRepository, authService, pool)
	// === Handlers ===
	eventHandler := api.NewHTTPHandler(eventRepository, bookingRepository, bookingService)
	authHandler := api.NewAuthHandler(userService)

	mux := http.NewServeMux()
	setupRoutes(mux, authService, eventHandler, authHandler)

	srv := setupServer(mux)

	go func() {
		if err := runServer(srv); err != nil {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	logger.Info("Server started on :8080")

	// Wait for a signal and then shutdown the server
	gracefulShutdown(srv, logger)
	return nil
}

func setupRoutes(
	mux *http.ServeMux,
	authService *auth.JWTService,
	eventHandler *api.HTTPHandler,
	authHandler *api.AuthHandler,
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

	// === Public endpoints ===
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// === Protected endpoints ===
	mux.HandleFunc("POST /events", auth(requireOrganizer(eventHandler.CreateEvent)))
	mux.HandleFunc("PUT /events/{id}", auth(requireOrganizer(eventHandler.UpdateEvent)))
	mux.HandleFunc("DELETE /events/{id}", auth(requireAdmin(eventHandler.DeleteEvent)))
	mux.HandleFunc("GET /events/{id}", auth(requireAll(eventHandler.GetEvent)))
	mux.HandleFunc("GET /events", auth(requireAll(eventHandler.ListEvents)))
	mux.HandleFunc("POST /events/{event_id}/bookings", auth(requireAll(eventHandler.CreateBooking)))
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

func gracefulShutdown(srv *http.Server, logger *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown server", "error", err)
	}

	logger.Info("Server shutdown successfully")
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
) (*services.BookingService, *services.UserService) {
	bookingService := services.NewBookingService(eventRepository, bookingRepository, pool)
	userService := services.NewUserService(userRepository, authService)
	return bookingService, userService
}
