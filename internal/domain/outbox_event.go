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
	aggregateID uuid.UUID
}

func CreateOutboxEvent(
	eventName string,
	eventData []byte,
	destination string,
	aggregateID uuid.UUID,
) (*OutboxEvent, error) {
	if eventName == "" || eventData == nil || destination == "" || aggregateID == uuid.Nil {
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
		aggregateID: aggregateID,
	}, nil
}

func ReconstructOutboxEvent(id uuid.UUID,
	eventName string,
	eventData []byte,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
	destination string,
	aggregateID uuid.UUID) *OutboxEvent {
	return &OutboxEvent{
		id:          id,
		eventName:   eventName,
		eventData:   eventData,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		destination: destination,
		aggregateID: aggregateID,
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

func (o *OutboxEvent) AggregateID() uuid.UUID {
	return o.aggregateID
}

type OutboxRepository interface {
	Create(ctx context.Context, event *OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error)
	MarkAsProcessed(ctx context.Context, id string) error
}
