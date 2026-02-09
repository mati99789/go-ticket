package api

import (
	"encoding/json"
	"net/http"
)

func responseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func responseCreated(w http.ResponseWriter, data any) {
	responseJSON(w, http.StatusCreated, data)
}

func responseOK(w http.ResponseWriter, data any) {
	responseJSON(w, http.StatusOK, data)
}

func responseNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
