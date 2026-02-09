// Package domain holds the business entities and rules.
package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Event represents an event in the system.
type Event struct {
	id             uuid.UUID
	name           string
	price          int64
	startAt        time.Time
	endAt          time.Time
	createdAt      time.Time
	updatedAt      time.Time
	capacity       int
	availableSpots int
}

// NewEvent creates a new validated Event.
func NewEvent(id uuid.UUID, name string, price int64, startAt time.Time, endAt time.Time, capacity int) (*Event, error) {
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
		id:             id,
		name:           name,
		price:          price,
		startAt:        startAt,
		endAt:          endAt,
		createdAt:      time.Now(),
		updatedAt:      time.Now(),
		capacity:       capacity,
		availableSpots: capacity,
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

// Capacity returns the event's capacity.
func (e *Event) Capacity() int {
	return e.capacity
}

// AvailableSpots returns the event's available spots.
func (e *Event) AvailableSpots() int {
	return e.availableSpots
}

// NewEventFromPersistence creates an Event from the given parameters.
func NewEventFromPersistence(id uuid.UUID,
	name string,
	price int64,
	startAt, endAt, createdAt, updatedAt time.Time, capacity int, availableSpots int) *Event {
	return &Event{id, name, price, startAt, endAt, createdAt, updatedAt, capacity, availableSpots}
}

// EventRepository defines the interface for event persistence.
type EventRepository interface {
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (*Event, error)
	ListEvents(ctx context.Context) ([]*Event, error)
	ReserveSpots(ctx context.Context, eventID uuid.UUID, spots int) error
}
