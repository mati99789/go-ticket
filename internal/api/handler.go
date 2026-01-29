package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/domain"
)

type HTTPHandler struct {
	eventRepository domain.EventRepository
}

func NewHTTPHandler(eventRepository domain.EventRepository) *HTTPHandler {
	return &HTTPHandler{eventRepository: eventRepository}
}

func (h *HTTPHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Name     string    `json:"name"`
		Price    int64     `json:"price"`
		StartAt  time.Time `json:"startAt"`
		EndAt    time.Time `json:"endAt"`
		Capacity int       `json:"capacity"`
	}

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	id := uuid.New()
	event, err := domain.NewEvent(id, req.Name, req.Price, req.StartAt, req.EndAt, req.Capacity)

	if err != nil {
		h.responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.eventRepository.CreateEvent(r.Context(), event)
	if err != nil {
		slog.Error("Failed to create event", "error", err)
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
}

func (h *HTTPHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.responseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.responseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	type request struct {
		Name    string    `json:"name"`
		Price   int64     `json:"price"`
		StartAt time.Time `json:"startAt"`
		EndAt   time.Time `json:"endAt"`
	}

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responseError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.eventRepository.GetEvent(r.Context(), parsedId)
	if err != nil {
		slog.Error("Failed to get event", "error", err)
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	err = event.UpdateName(req.Name)
	if err != nil {
		h.responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = event.Reschedule(req.StartAt, req.EndAt)
	if err != nil {
		h.responseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.eventRepository.UpdateEvent(r.Context(), event)
	if err != nil {
		slog.Error("Failed to update event", "error", err)
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.responseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.responseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.eventRepository.DeleteEvent(r.Context(), parsedId)
	if err != nil {
		slog.Error("Failed to delete event", "error", err)
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.responseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		h.responseError(w, http.StatusBadRequest, "invalid id")
		return
	}

	type responseStruct struct {
		ID      string    `json:"id"`
		Name    string    `json:"name"`
		Price   int64     `json:"price"`
		StartAt time.Time `json:"startAt"`
		EndAt   time.Time `json:"endAt"`
	}

	event, err := h.eventRepository.GetEvent(r.Context(), parsedId)
	if err != nil {
		slog.Error("Failed to get event", "error", err)
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	start, end := event.StartAndEndAt()

	resp := responseStruct{
		ID:      event.ID().String(),
		Name:    event.Name(),
		Price:   event.Price(),
		StartAt: start,
		EndAt:   end,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

}

func (h *HTTPHandler) ListEvents(w http.ResponseWriter, r *http.Request) {

	type responseStruct struct {
		ID      string    `json:"id"`
		Name    string    `json:"name"`
		Price   int64     `json:"price"`
		StartAt time.Time `json:"startAt"`
		EndAt   time.Time `json:"endAt"`
	}

	events, err := h.eventRepository.ListEvents(r.Context())
	if err != nil {
		slog.Error("Failed to list events", "error", err)
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	resp := []responseStruct{}
	for _, event := range events {
		start, end := event.StartAndEndAt()
		resp = append(resp, responseStruct{
			ID:      event.ID().String(),
			Name:    event.Name(),
			Price:   event.Price(),
			StartAt: start,
			EndAt:   end,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		h.responseError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

}

func (h *HTTPHandler) responseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
