package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/mati/go-ticket/internal/postgres"
	"github.com/stretchr/testify/assert"
)

func TestBookingService_CreateBooking_Success(t *testing.T) {
	ctx := context.Background()
	pool := postgres.SetupDb(ctx, t)

	// Create  test event
	event := postgres.CreateTestEvent(ctx, t, pool, postgres.WithCapacity(100))

	// Create booking service
	queries := postgres.New(pool)
	eventRepository := postgres.NewEventRepository(queries)
	bookingRepository := postgres.NewBookingRepository(queries)
	bookingService := NewBookingService(eventRepository, bookingRepository, pool)

	booking, err := domain.NewBooking(uuid.New(), event.ID(), "test@example.com", domain.BookingStatusPending)
	assert.NoError(t, err)

	// Create booking
	err = bookingService.CreateBooking(ctx, booking)
	assert.NoError(t, err)

	// Verify booking created
	retrievedBooking, err := bookingRepository.GetBookingByID(ctx, booking.ID())
	assert.NoError(t, err)
	assert.Equal(t, booking.ID(), retrievedBooking.ID())
	assert.Equal(t, booking.EventID(), retrievedBooking.EventID())
	assert.Equal(t, booking.UserEmail(), retrievedBooking.UserEmail())
	assert.Equal(t, booking.Status(), retrievedBooking.Status())

	// Verify spots reserved
	retrievedEvent := postgres.GetEventFromDB(ctx, t, pool, event.ID())
	assert.Equal(t, 99, retrievedEvent.AvailableSpots())
}

func TestBookingService_CreateBooking_EventNotFound(t *testing.T) {
	ctx := context.Background()
	pool := postgres.SetupDb(ctx, t)

	// Create booking service
	queries := postgres.New(pool)
	eventRepository := postgres.NewEventRepository(queries)
	bookingRepository := postgres.NewBookingRepository(queries)
	bookingService := NewBookingService(eventRepository, bookingRepository, pool)

	// Try to create booking for non-existent event
	fakeEventID := uuid.New()
	booking, err := domain.NewBooking(uuid.New(), fakeEventID, "test@example.com", domain.BookingStatusPending)
	assert.NoError(t, err)

	err = bookingService.CreateBooking(ctx, booking)

	// Verify error is domain.ErrEventNotFound
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestBookingService_CreateBooking_EventFull(t *testing.T) {
	ctx := context.Background()
	pool := postgres.SetupDb(ctx, t)

	// Create event with 0 capacity
	event := postgres.CreateTestEvent(ctx, t, pool, postgres.WithCapacity(0))

	// Create booking service
	queries := postgres.New(pool)
	eventRepository := postgres.NewEventRepository(queries)
	bookingRepository := postgres.NewBookingRepository(queries)
	bookingService := NewBookingService(eventRepository, bookingRepository, pool)

	// Try to create booking
	booking, err := domain.NewBooking(uuid.New(), event.ID(), "test@example.com", domain.BookingStatusPending)
	assert.NoError(t, err)

	err = bookingService.CreateBooking(ctx, booking)

	// Verify error is domain.ErrEventIsFull
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEventIsFull)
}

func TestBookingService_CreateBooking_TransactionRollback(t *testing.T) {
	ctx := context.Background()
	pool := postgres.SetupDb(ctx, t)

	// Create event with 1 spot
	event := postgres.CreateTestEvent(ctx, t, pool, postgres.WithCapacity(1))

	// Create booking service
	queries := postgres.New(pool)
	eventRepository := postgres.NewEventRepository(queries)
	bookingRepository := postgres.NewBookingRepository(queries)
	bookingService := NewBookingService(eventRepository, bookingRepository, pool)

	// Create first booking (should succeed)
	booking1, err := domain.NewBooking(uuid.New(), event.ID(), "test1@example.com", domain.BookingStatusPending)
	assert.NoError(t, err)

	err = bookingService.CreateBooking(ctx, booking1)
	assert.NoError(t, err)

	// Verify event spots were reserved
	retrievedEvent := postgres.GetEventFromDB(ctx, t, pool, event.ID())
	assert.Equal(t, 0, retrievedEvent.AvailableSpots(), "Event spots should be 0 after first booking")

	// Create second booking (should fail - event full)
	booking2, err := domain.NewBooking(uuid.New(), event.ID(), "test2@example.com", domain.BookingStatusPending)
	assert.NoError(t, err)

	err = bookingService.CreateBooking(ctx, booking2)

	// Verify error is domain.ErrEventIsFull
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEventIsFull, "Second booking should fail because event is full")

	// Verify transaction rolled back - event spots should still be 0
	retrievedEventAfterFailure := postgres.GetEventFromDB(ctx, t, pool, event.ID())
	assert.Equal(t, 0, retrievedEventAfterFailure.AvailableSpots(), "Event spots should remain 0 after failed booking attempt")

	// Verify second booking was not created
	_, err = bookingRepository.GetBookingByID(ctx, booking2.ID())
	assert.Error(t, err, "Second booking should not exist in database")
	assert.ErrorIs(t, err, domain.ErrBookingNotFound, "Second booking should not exist in database")
}
