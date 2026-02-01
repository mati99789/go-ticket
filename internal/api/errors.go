package api

import (
	"errors"
	"net/http"

	"github.com/mati/go-ticket/internal/domain"
)

// MapDomainError maps domain errors to HTTP status codes and user-friendly messages.
// Returns 500 Internal Server Error for unknown errors to avoid exposing sensitive information.
func MapDomainError(err error) (statusCode int, message string) {
	switch {
	// Event errors - Not Found (404)
	case errors.Is(err, domain.ErrEventNotFound):
		return http.StatusNotFound, "Event not found"

	// Event errors - Bad Request (400)
	case errors.Is(err, domain.ErrEventIsFull):
		return http.StatusBadRequest, "Event is full, no available spots"

	case errors.Is(err, domain.ErrEventNameEmpty):
		return http.StatusBadRequest, "Event name cannot be empty"

	case errors.Is(err, domain.ErrEventPriceNegative):
		return http.StatusBadRequest, "Event price must be positive"

	case errors.Is(err, domain.ErrEventStartAfterEnd):
		return http.StatusBadRequest, "Event start time must be before end time"

	case errors.Is(err, domain.ErrEventIDNil):
		return http.StatusBadRequest, "Invalid event ID"

	// Booking errors - Not Found (404)
	case errors.Is(err, domain.ErrBookingNotFound):
		return http.StatusNotFound, "Booking not found"

	// Booking errors - Bad Request (400)
	case errors.Is(err, domain.ErrBookingIDNil):
		return http.StatusBadRequest, "Invalid booking ID"

	case errors.Is(err, domain.ErrBookingEventIDInvalid):
		return http.StatusBadRequest, "Invalid event ID for booking"

	case errors.Is(err, domain.ErrBookingUserEmailEmpty):
		return http.StatusBadRequest, "User email is required"

	case errors.Is(err, domain.ErrBookingStatusInvalid):
		return http.StatusBadRequest, "Invalid booking status"

	// Default: Internal Server Error (500)
	// DO NOT expose internal error details to clients for security reasons
	default:
		return http.StatusInternalServerError, "An unexpected error occurred"
	}
}
