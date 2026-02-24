package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/mati/go-ticket/internal/postgres"
)

type CreateBookingService interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) error
}

type BookingService struct {
	eventRepo   *postgres.EventRepository
	bookingRepo *postgres.BookingRepository
	pool        *pgxpool.Pool
}

func NewBookingService(
	eventRepo *postgres.EventRepository,
	bookingRepo *postgres.BookingRepository,
	pool *pgxpool.Pool,
) *BookingService {
	return &BookingService{
		eventRepo:   eventRepo,
		bookingRepo: bookingRepo,
		pool:        pool,
	}
}

func (bs *BookingService) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	tx, err := bs.pool.Begin(ctx)
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("Error rolling back transaction", "error", err)
		}
	}()

	if err != nil {
		return err
	}

	if err := bs.eventRepo.WithTx(tx).ReserveSpots(ctx, booking.EventID(), 1); err != nil {
		return err
	}
	if err := bs.bookingRepo.WithTx(tx).CreateBooking(ctx, booking); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
