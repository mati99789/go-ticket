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
	"github.com/mati/go-ticket/internal/postgres"
	"github.com/mati/go-ticket/internal/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	// Load environment variables from .env file
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

	// Create database connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to create database connection pool: %w", err)
	}
	defer pool.Close()

	// Run database migrations
	logger.Info("Running database migrations...")
	if err := postgres.RunMigrations(dbUrl); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	logger.Info("Database migrations completed")

	// Create repositories and services
	queries := postgres.New(pool)
	eventRepository := postgres.NewEventRepository(queries)
	bookingRepository := postgres.NewBookingRepository(queries)
	bookingService := services.NewBookingService(eventRepository, bookingRepository, pool)
	eventHandler := api.NewHTTPHandler(eventRepository, bookingRepository, bookingService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /events", eventHandler.CreateEvent)
	mux.HandleFunc("PUT /events/{id}", eventHandler.UpdateEvent)
	mux.HandleFunc("DELETE /events/{id}", eventHandler.DeleteEvent)
	mux.HandleFunc("GET /events/{id}", eventHandler.GetEvent)
	mux.HandleFunc("GET /events", eventHandler.ListEvents)
	mux.HandleFunc("POST /events/{event_id}/bookings", eventHandler.CreateBooking)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      middleware.LoggingMiddleware(middleware.RecoveryMiddleware(mux)),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
		}
	}()

	logger.Info("Server started on :8080")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Wait for signal
	<-stop
	logger.Info("Shutting down server...")

	// Shutdown server
	ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	logger.Info("Server shutdown successfully")
	return nil
}
