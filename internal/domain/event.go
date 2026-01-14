package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventNameEmpty     = errors.New("name is empty")
	ErrEventPriceNegative = errors.New("price is negative")
	ErrEventStartAfterEnd = errors.New("startAt is after endAt")
	ErrEventIDNil         = errors.New("id is nil")
)

type Event struct {
	id        uuid.UUID
	name      string
	price     int64
	startAt   time.Time
	endAt     time.Time
	createdAt time.Time
	updatedAt time.Time
}

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

func (e *Event) UpdateName(name string) error {
	if name == "" {
		return ErrEventNameEmpty
	}
	e.name = name
	e.updatedAt = time.Now()
	return nil
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Price() int64 {
	return e.price
}

func (e *Event) StartAndEndAt() (time.Time, time.Time) {
	return e.startAt, e.endAt
}

func (e *Event) Reschedule(startAt time.Time, endAt time.Time) error {
	if startAt.After(endAt) {
		return ErrEventStartAfterEnd
	}
	e.startAt = startAt
	e.endAt = endAt
	e.updatedAt = time.Now()
	return nil
}
