package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	Pending   = "pending"
	Processed = "processed"
)

type OutboxEvent struct {
	id          uuid.UUID
	eventName   string
	eventData   []byte
	status      string
	createdAt   time.Time
	updatedAt   time.Time
	destination string
}

func CreateOutboxEvent(eventName string, eventData []byte, destination string) (*OutboxEvent, error) {
	if eventName == "" || eventData == nil || destination == "" {
		return nil, ErrOutboxEventInvalid
	}
	id := uuid.New()
	status := Pending
	createdAt := time.Now()
	updatedAt := time.Now()
	return &OutboxEvent{
		id:          id,
		eventName:   eventName,
		eventData:   eventData,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		destination: destination,
	}, nil
}

func ReconstructOutboxEvent(id uuid.UUID,
	eventName string,
	eventData []byte,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
	destination string) *OutboxEvent {
	return &OutboxEvent{
		id:          id,
		eventName:   eventName,
		eventData:   eventData,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		destination: destination,
	}
}

func (o *OutboxEvent) MarkAsProcessed() error {
	if o.status != Pending {
		return ErrOutboxEventInvalid
	}
	o.status = Processed
	o.updatedAt = time.Now()
	return nil
}

func (o *OutboxEvent) EventName() string {
	return o.eventName
}

func (o *OutboxEvent) EventData() []byte {
	return o.eventData
}

func (o *OutboxEvent) Destination() string {
	return o.destination
}

func (o *OutboxEvent) ID() uuid.UUID {
	return o.id
}

type OutboxRepository interface {
	Create(ctx context.Context, event *OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error)
	MarkAsProcessed(ctx context.Context, id string) error
}
