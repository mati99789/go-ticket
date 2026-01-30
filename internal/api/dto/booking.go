package dto

import (
	"time"

	"github.com/mati/go-ticket/internal/domain"
)

type CreateBookingRequest struct {
	UserEmail string `json:"userEmail"`
}

type BookingResponse struct {
	ID        string    `json:"id"`
	EventID   string    `json:"eventID"`
	UserEmail string    `json:"userEmail"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}

func ToBookingResponse(booking *domain.Booking) BookingResponse {
	return BookingResponse{
		ID:        booking.ID().String(),
		EventID:   booking.EventID().String(),
		UserEmail: booking.UserEmail(),
		CreatedAt: booking.CreatedAt(),
		Status:    string(booking.Status()),
	}
}

func ToBookingListResponse(bookings []*domain.Booking) []BookingResponse {
	responses := make([]BookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = ToBookingResponse(booking)
	}
	return responses
}
