package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func ResponseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}

func ResponseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
	}
}

func ResponseCreated(w http.ResponseWriter, data any) {
	ResponseJSON(w, http.StatusCreated, data)
}

func ResponseOK(w http.ResponseWriter, data any) {
	ResponseJSON(w, http.StatusOK, data)
}

func ResponseNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
