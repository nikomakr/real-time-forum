package handlers_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
)

func setupLogoutDB(t *testing.T) {
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

	db.DB = mockDB
}

func TestLogoutHandler(t *testing.T) {

	tests := []struct {
		name           string
		method         string
		cookie         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Reject invalid HTTP method",
			method:         http.MethodGet,
			cookie:         "valid-session-id",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Reject missing cookie",
			method:         http.MethodPost,
			cookie:         "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "not logged in",
		},
		{
			name:           "Successful logout clears session",
			method:         http.MethodPost,
			cookie:         "valid-session-id",
			expectedStatus: http.StatusOK,
			expectedBody:   "logged out successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupLogoutDB(t) // Reset DB state before each test
			req, _ := http.NewRequest(tt.method, "/api/logout", nil)
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "session_id", Value: tt.cookie})
			}

			rr := httptest.NewRecorder()
			http.HandlerFunc(handlers.Logout).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s",
					tt.expectedStatus, rr.Code, rr.Body.String())
			}

			if tt.expectedBody != "" && !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q",
					tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestLogoutClearsCookie(t *testing.T) {
	setupLogoutDB(t)

	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid-session-id"})

	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.Logout).ServeHTTP(rr, req)

	// Verify cookie is cleared
	found := false
	for _, c := range rr.Result().Cookies() {
		if c.Name == "session_id" {
			found = true
			if c.MaxAge != -1 {
				t.Errorf("expected MaxAge -1 to clear cookie, got %d", c.MaxAge)
			}
			break
		}
	}
	if !found {
		t.Error("expected session_id cookie in response to clear it")
	}

	// Verify session deleted from DB
	var count int
	db.DB.QueryRow(
		"SELECT COUNT(*) FROM sessions WHERE session_id = ?", "valid-session-id",
	).Scan(&count)
	if count != 0 {
		t.Error("expected session to be deleted from database but it still exists")
	}
}
