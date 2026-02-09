package api

import (
	"encoding/json"
	"net/http"
)

func ResponseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func ResponseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
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
