package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingEventPayload struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"eventID"`
	UserEmail string    `json:"userEmail"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}
