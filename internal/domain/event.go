// Package domain holds the business entities and rules.
package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrEventNameEmpty is returned when the name is empty.
	ErrEventNameEmpty = errors.New("name is empty")
	// ErrEventPriceNegative is returned when the price is negative.
	ErrEventPriceNegative = errors.New("price is negative")
	// ErrEventStartAfterEnd is returned when the start time is after the end time.
	ErrEventStartAfterEnd = errors.New("startAt is after endAt")
	// ErrEventIDNil is returned when the id is nil.
	ErrEventIDNil = errors.New("id is nil")
)

// Event represents an event in the system.
type Event struct {
	id        uuid.UUID
	name      string
	price     int64
	startAt   time.Time
	endAt     time.Time
	createdAt time.Time
	updatedAt time.Time
}

// NewEvent creates a new validated Event.
func NewEvent(id uuid.UUID, name string, price int64, startAt time.Time, endAt time.Time) (*Event, error) {
	if id == uuid.Nil {
		return nil, ErrEventIDNil
	}
	if name == "" {
		return nil, ErrEventNameEmpty
	}
	if price < 0 {
		return nil, ErrEventPriceNegative
	}
	if startAt.After(endAt) {
		return nil, ErrEventStartAfterEnd
	}
	return &Event{
		id:        id,
		name:      name,
		price:     price,
		startAt:   startAt,
		endAt:     endAt,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// UpdateName updates the event's name.
func (e *Event) UpdateName(name string) error {
	if name == "" {
		return ErrEventNameEmpty
	}
	e.name = name
	e.updatedAt = time.Now()
	return nil
}

// Name returns the event's name.
func (e *Event) Name() string {
	return e.name
}

// Price returns the event's price.
func (e *Event) Price() int64 {
	return e.price
}

// StartAndEndAt returns the event's start and end times.
func (e *Event) StartAndEndAt() (time.Time, time.Time) {
	return e.startAt, e.endAt
}

// Reschedule changes the event's start and end times.
func (e *Event) Reschedule(startAt time.Time, endAt time.Time) error {
	if startAt.After(endAt) {
		return ErrEventStartAfterEnd
	}
	e.startAt = startAt
	e.endAt = endAt
	e.updatedAt = time.Now()
	return nil
}

// ID returns the event's ID.
func (e *Event) ID() uuid.UUID {
	return e.id
}

// UnmarshalEvent creates an Event from the given parameters.
func UnmarshalEvent(id uuid.UUID,
	name string,
	price int64,
	startAt, endAt, createdAt, updatedAt time.Time) *Event {
	return &Event{id, name, price, startAt, endAt, createdAt, updatedAt}
}

// EventRepository defines the interface for event persistence.
type EventRepository interface {
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (*Event, error)
	ListEvents(ctx context.Context) ([]*Event, error)
}
