package event_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mati/go-ticket/internal/domain"
)

type BookingEventHandler struct {
	logger                *slog.Logger
	notificationPublisher domain.NotificationPublisher
}

func NewBookingEventHandler(
	logger *slog.Logger,
	notificationPublisher domain.NotificationPublisher) *BookingEventHandler {
	return &BookingEventHandler{
		logger:                logger,
		notificationPublisher: notificationPublisher,
	}
}

func (eh *BookingEventHandler) Handle(ctx context.Context, payload []byte) error {
	var booking domain.BookingEventPayload
	err := json.Unmarshal(payload, &booking)
	if err != nil {
		return fmt.Errorf("failed unmarshal booking event: %w", err)
	}

	bookingNotification := &domain.BookingNotification{
		ID:        booking.ID,
		EventID:   booking.EventID,
		UserEmail: booking.UserEmail,
		CreatedAt: booking.CreatedAt,
		Status:    booking.Status,
	}
	err = eh.notificationPublisher.Publish(ctx, bookingNotification)
	if err != nil {
		return fmt.Errorf("failed publish event: %w", err)
	}
	eh.logger.Info("booking event received",
		"booking_id", booking.ID,
		"user_email", booking.UserEmail,
		"status", booking.Status,
	)
	return nil
}
