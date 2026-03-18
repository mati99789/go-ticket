package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type BookingNotification struct {
	ID        uuid.UUID
	EventID   uuid.UUID
	UserEmail string
	CreatedAt time.Time
	Status    string
}
type NotificationPublisher interface {
	Publish(ctx context.Context, payload *BookingNotification) error
}
