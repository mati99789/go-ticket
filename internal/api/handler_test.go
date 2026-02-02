package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/api/dto"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

type MockCreateBookingService struct {
	OnCreateBooking func(ctx context.Context, booking *domain.Booking) error
}

func (m *MockCreateBookingService) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	if m.OnCreateBooking != nil {
		return m.OnCreateBooking(ctx, booking)
	}
	return nil
}

func TestCreateBooking_Success(t *testing.T) {
	validEventID := uuid.New()
	validEmail := "user@example.com"

	mockCreateBookingService := &MockCreateBookingService{
		OnCreateBooking: func(ctx context.Context, booking *domain.Booking) error {
			if booking.EventID() != validEventID {
				assert.Equal(t, validEventID, booking.EventID())
			}

			if booking.UserEmail() != validEmail {
				assert.Equal(t, validEmail, booking.UserEmail())
			}

			if booking.Status() != domain.BookingStatusPending {
				assert.Equal(t, domain.BookingStatusPending, booking.Status())
			}

			return nil
		},
	}

	handler := NewHTTPHandler(nil, nil, mockCreateBookingService)

	reqBody := dto.CreateBookingRequest{
		UserEmail: validEmail,
	}

	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", fmt.Sprintf("/events/%s/bookings", validEventID), bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	req.SetPathValue("event_id", validEventID.String())

	recorder := httptest.NewRecorder()

	handler.CreateBooking(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreate_Booking_EventFull(t *testing.T) {
	validEventID := uuid.New()
	validEmail := "user@example.com"

	mockCreateBookingService := &MockCreateBookingService{
		OnCreateBooking: func(ctx context.Context, booking *domain.Booking) error {
			return domain.ErrEventIsFull
		},
	}

	handler := NewHTTPHandler(nil, nil, mockCreateBookingService)

	reqBody := dto.CreateBookingRequest{
		UserEmail: validEmail,
	}

	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", fmt.Sprintf("/events/%s/bookings", validEventID), bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	req.SetPathValue("event_id", validEventID.String())

	recorder := httptest.NewRecorder()

	handler.CreateBooking(recorder, req)

	assert.Equal(t, http.StatusConflict, recorder.Code)
}

func TestCreateBooking_EventNotFound(t *testing.T) {
	validEventID := uuid.New()
	validEmail := "user@example.com"

	mockCreateBookingService := &MockCreateBookingService{
		OnCreateBooking: func(ctx context.Context, booking *domain.Booking) error {
			return domain.ErrEventNotFound
		},
	}

	handler := NewHTTPHandler(nil, nil, mockCreateBookingService)

	reqBody := dto.CreateBookingRequest{
		UserEmail: validEmail,
	}

	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", fmt.Sprintf("/events/%s/bookings", validEventID), bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	req.SetPathValue("event_id", validEventID.String())

	recorder := httptest.NewRecorder()

	handler.CreateBooking(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}
