package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

var (
	ErrBookingIDNil          = errors.New("id is nil")
	ErrBookingEventIDInvalid = errors.New("eventID is invalid")
	ErrBookingUserEmailEmpty = errors.New("userEmail is empty")
	ErrBookingStatusInvalid  = errors.New("invalid status")
)

type Booking struct {
	id        uuid.UUID
	eventID   uuid.UUID
	userEmail string
	createdAt time.Time
	updatedAt time.Time
	status    BookingStatus
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *Booking) error
	GetBookingByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	ListBookings(ctx context.Context) ([]Booking, error)
	UpdateBooking(ctx context.Context, booking *Booking) error
	DeleteBooking(ctx context.Context, id uuid.UUID) error
	ConfirmBooking(ctx context.Context, id uuid.UUID) error
	CancelBooking(ctx context.Context, id uuid.UUID) error
}

func NewBooking(id uuid.UUID, eventID uuid.UUID, userEmail string, status BookingStatus) (*Booking, error) {
	if id == uuid.Nil {
		return nil, ErrBookingIDNil
	}
	if eventID == uuid.Nil {
		return nil, ErrBookingEventIDInvalid
	}
	if userEmail == "" {
		return nil, ErrBookingUserEmailEmpty
	}
	if status != BookingStatusPending && status != BookingStatusConfirmed && status != BookingStatusCancelled {
		return nil, ErrBookingStatusInvalid
	}
	return &Booking{
		id:        id,
		eventID:   eventID,
		userEmail: userEmail,
		status:    status,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func (b *Booking) Confirm() error {
	if b.status == BookingStatusCancelled {
		return errors.New("cannot confirm a cancelled booking")
	}
	b.status = BookingStatusConfirmed
	b.updatedAt = time.Now()
	return nil
}

func (b *Booking) Cancel() error {
	if b.status == BookingStatusConfirmed {
		return errors.New("cannot cancel a confirmed booking")
	}
	b.status = BookingStatusCancelled
	b.updatedAt = time.Now()
	return nil
}

func (b *Booking) ID() uuid.UUID {
	return b.id
}

func (b *Booking) EventID() uuid.UUID {
	return b.eventID
}

func (b *Booking) UserEmail() string {
	return b.userEmail
}

func (b *Booking) Status() BookingStatus {
	return b.status
}

func (b *Booking) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Booking) UpdatedAt() time.Time {
	return b.updatedAt
}

func UnmarshalBooking(id uuid.UUID, eventID uuid.UUID, userEmail string, status BookingStatus, createdAt, updatedAt time.Time) *Booking {
	return &Booking{
		id:        id,
		eventID:   eventID,
		userEmail: userEmail,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
