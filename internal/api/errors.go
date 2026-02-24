package api

import (
	"errors"
	"net/http"

	"github.com/mati/go-ticket/internal/domain"
)

// MapDomainError maps domain errors to HTTP status codes and user-friendly messages.
// Returns 500 Internal Server Error for unknown errors to avoid exposing sensitive information.
func MapDomainError(err error) (statusCode int, message string) { //nolint:gocyclo
	// TODO: extract error mappings to separate functions
	switch {
	// Event errors - Not Found (404)
	case errors.Is(err, domain.ErrEventNotFound):
		return http.StatusNotFound, "Event not found"

	// Event errors - Bad Request (400)
	case errors.Is(err, domain.ErrEventIsFull):
		return http.StatusConflict, "Event is full, no available spots"

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

	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "User not found"

	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "Invalid credentials"

	case errors.Is(err, domain.ErrUserPasswordTooShort):
		return http.StatusBadRequest, "Password is too short"

	case errors.Is(err, domain.ErrUserEmailEmpty):
		return http.StatusBadRequest, "Email is required"

	case errors.Is(err, domain.ErrUserEmailAlreadyExists):
		return http.StatusConflict, "User already exists"

	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "Invalid credentials"

	// Default: Internal Server Error (500)
	// DO NOT expose internal error details to clients for security reasons
	default:
		return http.StatusInternalServerError, "An unexpected error occurred"
	}
}
