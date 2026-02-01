package domain

import "errors"

// Event errors
var (
	ErrEventNotFound = errors.New("event not found")
	ErrEventIsFull   = errors.New("event is full")
	// ErrEventNameEmpty is returned when the name is empty.
	ErrEventNameEmpty = errors.New("name is empty")
	// ErrEventPriceNegative is returned when the price is negative.
	ErrEventPriceNegative = errors.New("price is negative")
	// ErrEventStartAfterEnd is returned when the start time is after the end time.
	ErrEventStartAfterEnd = errors.New("startAt is after endAt")
	// ErrEventIDNil is returned when the id is nil.
	ErrEventIDNil = errors.New("id is nil")
)

// Booking errors
var (
	// ErrBookingNotFound is returned when the booking is not found.
	ErrBookingNotFound = errors.New("booking not found")
	// ErrBookingIDNil is returned when the id is nil.
	ErrBookingIDNil = errors.New("id is nil")
	// ErrBookingEventIDInvalid is returned when the eventID is invalid.
	ErrBookingEventIDInvalid = errors.New("eventID is invalid")
	// ErrBookingUserEmailEmpty is returned when the userEmail is empty.
	ErrBookingUserEmailEmpty = errors.New("userEmail is empty")
	// ErrBookingStatusInvalid is returned when the status is invalid.
	ErrBookingStatusInvalid = errors.New("invalid status")
)
