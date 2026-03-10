package services

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mati/go-ticket/internal/api/dto"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/mati/go-ticket/internal/postgres"
)

type CreateBookingService interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) error
}

type BookingService struct {
	eventRepo   *postgres.EventRepository
	bookingRepo *postgres.BookingRepository
	outboxRepo  *postgres.OutBoxRepository
	pool        *pgxpool.Pool
}

func NewBookingService(
	eventRepo *postgres.EventRepository,
	bookingRepo *postgres.BookingRepository,
	outboxRepo *postgres.OutBoxRepository,
	pool *pgxpool.Pool,
) *BookingService {
	return &BookingService{
		eventRepo:   eventRepo,
		bookingRepo: bookingRepo,
		outboxRepo:  outboxRepo,
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

	eventData, err := json.Marshal(dto.ToBookingResponse(booking))
	if err != nil {
		return err
	}
	outboxEvent, err := domain.CreateOutboxEvent(
		"CreateBooking",
		eventData,
		"booking_events_topic",
		booking.ID(),
	)
	if err != nil {
		return err
	}
	if err := bs.outboxRepo.WithTx(tx).Create(ctx, outboxEvent); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	slog.Info("Crated Booking and Outbox Event", "booking", dto.ToBookingResponse(booking), "outboxEvent", outboxEvent)

	return nil
}
