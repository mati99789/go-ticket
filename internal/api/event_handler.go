package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/api/dto"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/mati/go-ticket/internal/services"
)

type HTTPHandler struct {
	eventRepository   domain.EventRepository
	bookingRepository domain.BookingRepository
	bookingService    services.CreateBookingService
}

func NewHTTPHandler(
	eventRepository domain.EventRepository,
	bookingRepository domain.BookingRepository,
	bookingService services.CreateBookingService,
) *HTTPHandler {
	return &HTTPHandler{
		eventRepository:   eventRepository,
		bookingRepository: bookingRepository,
		bookingService:    bookingService,
	}
}

func (h *HTTPHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id := uuid.New()
	event, err := domain.NewEvent(id, req.Name, req.Price, req.StartAt, req.EndAt, req.Capacity)

	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	err = h.eventRepository.CreateEvent(r.Context(), event)
	if err != nil {
		slog.Error("Failed to create event", "error", err)
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	ResponseCreated(w, map[string]string{"id": id.String()})
}

func (h *HTTPHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		ResponseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	event, err := h.eventRepository.GetEvent(r.Context(), parsedId)
	if err != nil {
		slog.Error("Failed to get event", "error", err)
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	err = event.UpdateName(req.Name)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	err = event.Reschedule(req.StartAt, req.EndAt)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	err = h.eventRepository.UpdateEvent(r.Context(), event)
	if err != nil {
		slog.Error("Failed to update event", "error", err)
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		ResponseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.eventRepository.DeleteEvent(r.Context(), parsedId)
	if err != nil {
		slog.Error("Failed to delete event", "error", err)
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		ResponseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	event, err := h.eventRepository.GetEvent(r.Context(), parsedId)
	if err != nil {
		slog.Error("Failed to get event", "error", err)
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	resp := dto.ToEventResponse(event)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}
}

func (h *HTTPHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventRepository.ListEvents(r.Context())
	if err != nil {
		slog.Error("Failed to list events", "error", err)
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	resp := dto.ToEventListResponse(events)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}
}

func (h *HTTPHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	var req dto.CreateBookingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	id := uuid.New()

	booking, err := domain.NewBooking(id, eventID, req.UserEmail, domain.BookingStatusPending)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	err = h.bookingService.CreateBooking(r.Context(), booking)
	if err != nil {
		code, message := MapDomainError(err)
		ResponseError(w, code, message)
		return
	}

	ResponseCreated(w, dto.ToBookingResponse(booking))
}
