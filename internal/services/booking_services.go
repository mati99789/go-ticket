package services

import (
	"context"
	"encoding/json"
	"log/slog"

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
	tm          domain.TransactionManager
}

func NewBookingService(
	eventRepo *postgres.EventRepository,
	bookingRepo *postgres.BookingRepository,
	outboxRepo *postgres.OutBoxRepository,
	pool domain.TransactionManager,
) *BookingService {
	return &BookingService{
		eventRepo:   eventRepo,
		bookingRepo: bookingRepo,
		outboxRepo:  outboxRepo,
		tm:          pool,
	}
}

func (bs *BookingService) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	err := bs.tm.RunInTx(ctx, func(ctx context.Context) error {
		if err := bs.eventRepo.ReserveSpots(ctx, booking.EventID(), 1); err != nil {
			return err
		}
		if err := bs.bookingRepo.CreateBooking(ctx, booking); err != nil {
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
		if err := bs.outboxRepo.Create(ctx, outboxEvent); err != nil {
			return err
		}
		slog.Info("Crated Booking and Outbox Event", "booking", dto.ToBookingResponse(booking), "outboxEvent", outboxEvent)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
