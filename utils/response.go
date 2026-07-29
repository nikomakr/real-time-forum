package utils

import (
	"bytes"
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
	// NOTE: w.Header().Set() must always be called before w.WriteHeader().
	// Once WriteHeader is called, headers are locked and any further Set() calls are silently ignored by the Go HTTP runtime.
	// Encode into a buffer first — if encoding fails, no headers have been written yet so the caller can still return a proper error response to the client.
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err := buf.WriteTo(w)
	return err
}
