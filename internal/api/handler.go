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

	id := uuid.New()
	event, err := domain.NewEvent(id, req.Name, req.Price, req.StartAt, req.EndAt)

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

}

func (h *HTTPHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) GetEvent(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) ListEvents(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) responseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
