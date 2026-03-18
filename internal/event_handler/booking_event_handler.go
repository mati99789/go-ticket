package event_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

type BookingEventHandler struct {
	logger *slog.Logger
}

func NewBookingEventHandler(logger *slog.Logger) *BookingEventHandler {
	return &BookingEventHandler{
		logger: logger,
	}
}

func (eh *BookingEventHandler) Handle(_ context.Context, payload []byte) error {
	var booking BookingEventPayload
	err := json.Unmarshal(payload, &booking)
	if err != nil {
		return fmt.Errorf("failed unmarshal booking event: %w", err)
	}

	eh.logger.Info("booking event received",
		"booking_id", booking.ID,
		"user_email", booking.UserEmail,
		"status", booking.Status,
	)
	return nil
}
