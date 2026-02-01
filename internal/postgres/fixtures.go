package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mati/go-ticket/internal/domain"
)

type EventOptions func(*EventConfig)

type EventConfig struct {
	Name     string
	Price    int64
	StartAt  time.Time
	EndAt    time.Time
	Capacity int
}

func WithName(name string) EventOptions {
	return func(config *EventConfig) {
		config.Name = name
	}
}

func WithPrice(price int64) EventOptions {
	return func(config *EventConfig) {
		config.Price = price
	}
}

func WithStartAt(startAt time.Time) EventOptions {
	return func(config *EventConfig) {
		config.StartAt = startAt
	}
}

func WithEndAt(endAt time.Time) EventOptions {
	return func(config *EventConfig) {
		config.EndAt = endAt
	}
}

func WithCapacity(capacity int) EventOptions {
	return func(config *EventConfig) {
		config.Capacity = capacity
	}
}

func CreateTestEvent(ctx context.Context, t *testing.T, pool *pgxpool.Pool, options ...EventOptions) *domain.Event {
	t.Helper()

	config := &EventConfig{
		Name:     "Test Event",
		Price:    1000,
		StartAt:  time.Now().Add(1 * time.Hour),
		EndAt:    time.Now().Add(2 * time.Hour),
		Capacity: 10,
	}

	for _, option := range options {
		option(config)
	}

	newEvent, err := domain.NewEvent(uuid.New(), config.Name, config.Price, config.StartAt, config.EndAt, config.Capacity)

	if err != nil {
		t.Fatalf("failed to create test event: %v", err)
	}

	queries := New(pool)

	eventRepositry := NewEventRepository(queries)

	err = eventRepositry.CreateEvent(ctx, newEvent)
	if err != nil {
		t.Fatalf("failed to create test event: %v", err)
	}

	return newEvent
}

func GetBookingFromDB(ctx context.Context, t *testing.T, pool *pgxpool.Pool, eventID uuid.UUID) *domain.Booking {
	t.Helper()

	queries := New(pool)

	bookingRepository := NewBookingRepository(queries)

	booking, err := bookingRepository.GetBookingByID(ctx, eventID)
	if err != nil {
		t.Fatalf("failed to get booking from db: %v", err)
	}

	return booking
}

func GetEventFromDB(ctx context.Context, t *testing.T, pool *pgxpool.Pool, eventID uuid.UUID) *domain.Event {
	t.Helper()

	queries := New(pool)

	eventRepository := NewEventRepository(queries)

	event, err := eventRepository.GetEvent(ctx, eventID)
	if err != nil {
		t.Fatalf("failed to get event from db: %v", err)
	}

	return event
}
