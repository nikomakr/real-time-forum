package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
)

func setupUploadDB(t *testing.T) {
	mockDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}
	mockDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		db.DB = nil
		mockDB.Close()
	})

	mockDB.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY, nickname TEXT UNIQUE,
		email TEXT UNIQUE, password_hash TEXT,
		first_name TEXT, last_name TEXT,
		age INTEGER, gender TEXT
	)`)

	mockDB.Exec(`CREATE TABLE sessions (
		session_id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	db.DB = mockDB

	// UploadImage writes to a relative "static/uploads" dir off the process
	// working directory (the package dir under `go test`) — clean it up.
	t.Cleanup(func() {
		os.RemoveAll("static")
	})
}

// pngBytes is the minimal 8-byte PNG signature — enough for http.DetectContentType
// to sniff it as image/png without needing a fully valid image.
var pngBytes = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

func multipartImageRequest(t *testing.T, fieldName, filename string, content []byte) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if fieldName != "" {
		part, err := writer.CreateFormFile(fieldName, filename)
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		if _, err := part.Write(content); err != nil {
			t.Fatalf("failed to write form file content: %v", err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/messages/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	return req
}

func TestUploadImage_RejectsNonPost(t *testing.T) {
	setupUploadDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/messages/upload", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.UploadImage).ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestUploadImage_RejectsUnauthenticated(t *testing.T) {
	setupUploadDB(t)

	req := multipartImageRequest(t, "image", "photo.png", pngBytes)
	req.Header.Del("Cookie")
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.UploadImage).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestUploadImage_RejectsMissingFile(t *testing.T) {
	setupUploadDB(t)

	req := multipartImageRequest(t, "", "", nil)
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.UploadImage).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d. Body: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "image file is required") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
}

func TestUploadImage_RejectsNonImageFile(t *testing.T) {
	setupUploadDB(t)

	req := multipartImageRequest(t, "image", "notes.txt", []byte("just plain text, not an image"))
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.UploadImage).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d. Body: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "unsupported image type") {
		t.Errorf("unexpected body: %s", rr.Body.String())
	}
}

func TestUploadImage_Success(t *testing.T) {
	setupUploadDB(t)

	req := multipartImageRequest(t, "image", "photo.png", pngBytes)
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.UploadImage).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp struct {
		ImageURL string `json:"image_url"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !strings.HasPrefix(resp.ImageURL, "/uploads/") || !strings.HasSuffix(resp.ImageURL, ".png") {
		t.Errorf("expected image_url like /uploads/<id>.png, got %q", resp.ImageURL)
	}

	savedPath := "static/uploads/" + strings.TrimPrefix(resp.ImageURL, "/uploads/")
	if _, err := os.Stat(savedPath); err != nil {
		t.Errorf("expected uploaded file to exist at %q: %v", savedPath, err)
	}
}
