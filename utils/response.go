package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, message string) {
	// Good practice to include charset. Ignores error here to avoid an infinite loop
	_ = WriteJSON(w, status, ErrorResponse{Error: message})
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	// Specifying UTF-8 prevents unexpected encoding issues
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// Returning the error allows the caller to log serialisation failures
	return json.NewEncoder(w).Encode(data)
}
