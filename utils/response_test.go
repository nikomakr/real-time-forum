package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_ProveEncoderFailure(t *testing.T) {
	rec := httptest.NewRecorder()
	unmarshallableData := map[string]any{
		"invalid": func() {},
	}
	WriteJSON(rec, http.StatusCreated, unmarshallableData)
	t.Logf("Status Code Sent: %d", rec.Code)
	t.Logf("Body Received: %q", rec.Body.String())
	if rec.Code == http.StatusCreated && rec.Body.String() == "" {
		t.Error("PROOF: client received 201 with empty body — encoder failed silently")
	}
}

func Test_ProveHeaderFreezing(t *testing.T) {
	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)
	rec.Header().Set("Content-Type", "application/json")
	contentType := rec.Result().Header.Get("Content-Type")
	t.Logf("Resulting Content-Type: %q", contentType)
	if contentType != "application/json" {
		t.Error("PROOF: Content-Type dropped — WriteHeader locked the headers before Set() was called")
	}
}

func Test_WriteJSON_HappyPath(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, http.StatusOK, map[string]string{"message": "ok"})
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "ok") {
		t.Errorf("expected body to contain 'ok', got %q", rec.Body.String())
	}
}

func Test_WriteError_HappyPath(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusBadRequest, "nickname is required")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "nickname is required") {
		t.Errorf("expected body to contain error message, got %q", rec.Body.String())
	}
}
