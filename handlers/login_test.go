package handlers_test

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
	"real-time-forum/utils"
)

func setupLoginDB(t *testing.T) {
	mockDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open mock DB: %v", err)
	}

	// SQLite in-memory databases are per connection. Setting max open connections to 1 ensures the same connection is used throughout the test, preserving the in-memory state.
	mockDB.SetMaxOpenConns(1)

	t.Cleanup(func() {
		db.DB = nil
		mockDB.Close()
	})

	// Users must be created before sessions — sessions has a foreign key to users
	_, err = mockDB.Exec(`
		CREATE TABLE users (
			id            TEXT PRIMARY KEY,
			nickname      TEXT UNIQUE,
			first_name    TEXT,
			last_name     TEXT,
			email         TEXT UNIQUE,
			password_hash TEXT,
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

	hash, err := utils.HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	_, err = mockDB.Exec(
		`INSERT INTO users (id, nickname, email, password_hash) VALUES (?, ?, ?, ?)`,
		"uuid-123", "niko", "niko@test.com", hash,
	)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	db.DB = mockDB
}

func TestLoginHandler(t *testing.T) {
	setupLoginDB(t)

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedBody   string
		expectedCookie string
	}{
		{
			name:           "Reject invalid HTTP method",
			payload:        `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Reject missing identifier",
			payload:        `{"password":"secret123"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "nickname or email is required",
		},
		{
			name:           "Reject missing password",
			payload:        `{"identifier":"niko"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "password is required",
		},
		{
			name:           "Reject unregistered user",
			payload:        `{"identifier":"ghost","password":"secret123"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid credentials",
		},
		{
			name:           "Reject wrong password",
			payload:        `{"identifier":"niko","password":"wrongpassword"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid credentials",
		},
		{
			name:           "Login with nickname",
			payload:        `{"identifier":"niko","password":"secret123"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "login successful",
			expectedCookie: "session_id",
		},
		{
			name:           "Login with email",
			payload:        `{"identifier":"niko@test.com","password":"secret123"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "login successful",
			expectedCookie: "session_id",
		},
		{
			name:           "Reject oversized body (DoS Protection)",
			payload:        `{"identifier":"` + strings.Repeat("A", 1024*1025) + `","password":"123"}`,
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectedBody:   "request body too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := http.MethodPost
			if tt.name == "Reject invalid HTTP method" {
				method = http.MethodGet
			}

			req, err := http.NewRequest(method, "/api/login", bytes.NewBufferString(tt.payload))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			http.HandlerFunc(handlers.Login).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s",
					tt.expectedStatus, rr.Code, rr.Body.String())
			}

			if tt.expectedBody != "" && !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q",
					tt.expectedBody, rr.Body.String())
			}
			if tt.expectedCookie != "" {
				found := false
				for _, c := range rr.Result().Cookies() {
					if c.Name == tt.expectedCookie {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected cookie %q to be set but it was not", tt.expectedCookie)
				}
			}
		})
	}
}
