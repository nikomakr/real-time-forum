package handlers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest" // For testing HTTP servers, provides utilities to create mock HTTP requests and record responses, allowing to test handlers without having a live server.
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
)

func setupMiddlewareDB(t *testing.T) {
	mockDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}

	mockDB.SetMaxOpenConns(1)

	t.Cleanup(func() {
		db.DB = nil
		mockDB.Close()
	})

	_, err = mockDB.Exec(`
		CREATE TABLE users (
			id            TEXT PRIMARY KEY,
			nickname      TEXT UNIQUE,
			email         TEXT UNIQUE,
			password_hash TEXT,
			first_name    TEXT,
			last_name     TEXT,
			age           INTEGER,
			gender        TEXT
		);
	`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	_, err = mockDB.Exec(`
		CREATE TABLE sessions (
			session_id TEXT PRIMARY KEY,
			user_id    TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatalf("failed to create sessions table: %v", err)
	}

	_, err = mockDB.Exec(
		`INSERT INTO users (id, nickname, email, password_hash) VALUES (?, ?, ?, ?)`,
		"uuid-123", "niko", "niko@test.com", "hash",
	)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	_, err = mockDB.Exec(
		`INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)`,
		"valid-session-id", "uuid-123", time.Now().UTC().Add(24*time.Hour),
	)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	_, err = mockDB.Exec(
		`INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)`,
		"expired-session-id", "uuid-123", time.Now().UTC().Add(-1*time.Hour),
	)
	if err != nil {
		t.Fatalf("failed to seed expired session: %v", err)
	}

	db.DB = mockDB
}

func TestRequireAuth(t *testing.T) {
	setupMiddlewareDB(t)

	protected := handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		cookie         string
		expectedStatus int
	}{
		{
			name:           "Reject request with no cookie",
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Reject request with invalid session",
			cookie:         "invalid-session-id",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Reject request with expired session",
			cookie:         "expired-session-id",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Allow request with valid session",
			cookie:         "valid-session-id",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/me", nil)
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: tt.cookie})
			}

			rr := httptest.NewRecorder()
			protected.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}