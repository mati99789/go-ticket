package api

import (
	"errors"
	"net/http"

	"github.com/mati/go-ticket/internal/domain"
)

type errorMapping struct {
	StatusCode int
	Message    string
}

var errorsMap = map[error]errorMapping{
	domain.ErrEventNotFound:          {http.StatusNotFound, "Event not found"},
	domain.ErrEventIsFull:            {http.StatusConflict, "Event is full, no available spots"},
	domain.ErrEventNameEmpty:         {http.StatusBadRequest, "Event name cannot be empty"},
	domain.ErrEventPriceNegative:     {http.StatusBadRequest, "Event price must be positive"},
	domain.ErrEventStartAfterEnd:     {http.StatusBadRequest, "Event start time must be before end time"},
	domain.ErrEventIDNil:             {http.StatusBadRequest, "Invalid event ID"},
	domain.ErrBookingNotFound:        {http.StatusNotFound, "Booking not found"},
	domain.ErrBookingIDNil:           {http.StatusBadRequest, "Invalid booking ID"},
	domain.ErrBookingEventIDInvalid:  {http.StatusBadRequest, "Invalid event ID for booking"},
	domain.ErrBookingUserEmailEmpty:  {http.StatusBadRequest, "User email is required"},
	domain.ErrBookingStatusInvalid:   {http.StatusBadRequest, "Invalid booking status"},
	domain.ErrUserNotFound:           {http.StatusNotFound, "User not found"},
	domain.ErrInvalidCredentials:     {http.StatusUnauthorized, "Invalid credentials"},
	domain.ErrUserPasswordTooShort:   {http.StatusBadRequest, "Password is too short"},
	domain.ErrUserEmailEmpty:         {http.StatusBadRequest, "Email is required"},
	domain.ErrUserEmailAlreadyExists: {http.StatusConflict, "User already exists"},
}

// MapDomainError maps domain errors to HTTP status codes and user-friendly messages.
// Returns 500 Internal Server Error for unknown errors to avoid exposing sensitive information.
func MapDomainError(err error) (statusCode int, message string) {
	for domainErr, mapping := range errorsMap {
		if errors.Is(err, domainErr) {
			return mapping.StatusCode, mapping.Message
		}
	}

	return http.StatusInternalServerError, "An unexpected error occurred"
}
