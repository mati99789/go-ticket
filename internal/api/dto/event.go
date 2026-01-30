package dto

import (
	"time"

	"github.com/mati/go-ticket/internal/domain"
)

// Request DTOs
type CreateEventRequest struct {
	Name     string    `json:"name"`
	StartAt  time.Time `json:"startAt"`
	EndAt    time.Time `json:"endAt"`
	Price    int64     `json:"price"`
	Capacity int       `json:"capacity"`
}

type UpdateEventRequest struct {
	Name    string    `json:"name"`
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
	Price   int64     `json:"price"`
}

// Response DTOs
type EventResponse struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Price          int64     `json:"price"`
	StartAt        time.Time `json:"startAt"`
	EndAt          time.Time `json:"endAt"`
	Capacity       int       `json:"capacity"`
	AvailableSpots int       `json:"availableSpots"`
}

func ToEventResponse(event *domain.Event) EventResponse {
	startAt, endAt := event.StartAndEndAt()
	return EventResponse{
		ID:             event.ID().String(),
		Name:           event.Name(),
		Price:          event.Price(),
		StartAt:        startAt,
		EndAt:          endAt,
		Capacity:       event.Capacity(),
		AvailableSpots: event.AvailableSpots(),
	}
}

func ToEventListResponse(events []*domain.Event) []EventResponse {
	responses := make([]EventResponse, len(events))
	for i, event := range events {
		responses[i] = ToEventResponse(event)
	}
	return responses
}
