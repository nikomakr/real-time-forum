package handlers_test

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"real-time-forum/db"
	"real-time-forum/handlers" // Adjust to your actual package path
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupMockDB initialises an in-memory SQLite database for test isolation aka "test doubles" and sets it as the global db.DB for the handlers to use.
func setupMockDB(t *testing.T) {
	mockDB, err := sql.Open("sqlite3", "file::memory:?cache=shared&mode=memory")
	mockDB.SetMaxOpenConns(1)
	if err != nil {
		t.Fatalf("Failed to open mock DB: %v", err)
	}

	statement := `
	CREATE TABLE users (
		id TEXT PRIMARY KEY,
		nickname TEXT UNIQUE,
		first_name TEXT,
		last_name TEXT,
		email TEXT UNIQUE,
		password_hash TEXT,
		age INTEGER,
		gender TEXT
	);`
	if _, err := mockDB.Exec(statement); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Seed data for duplicate test constraints
	seedQuery := `INSERT INTO users (id, nickname, email) VALUES ('uuid-123', 'existing_user', 'taken@email.com')`
	if _, err := mockDB.Exec(seedQuery); err != nil {
		t.Fatalf("Failed to seed mock data: %v", err)
	}

	db.DB = mockDB
}

func TestRegisterHandlerVulnerabilities(t *testing.T) {
	setupMockDB(t)

	tests := []struct {
		name           string
		method         string
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Bug 1: Reject Invalid HTTP Method",
			method:         http.MethodGet,
			payload:        `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Bug 2: Reject Malformed JSON Syntax",
			method:         http.MethodPost,
			payload:        `{"nickname": "test", "age": }`, // Broken JSON syntax
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Bug 3: Enforce All Mandatory Fields (Missing Email)",
			method:         http.MethodPost,
			payload:        `{"nickname":"bobby","first_name":"Bob","last_name":"Doe","password":"password123","age":25,"gender":"male"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Bug 4: Enforce Non-Zero Integer Fields (Age is 0)",
			method:         http.MethodPost,
			payload:        `{"nickname":"bobby","first_name":"Bob","last_name":"Doe","email":"bob@test.com","password":"password123","age":0,"gender":"male"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Bug 5: Collision Catch - Nickname Already Taken",
			method:         http.MethodPost,
			payload:        `{"nickname":"existing_user","first_name":"Bob","last_name":"Doe","email":"unique@email.com","password":"password123","age":25,"gender":"male"}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   "nickname already taken",
		},
		{
			name:           "Bug 6: Collision Catch - Email Already Registered",
			method:         http.MethodPost,
			payload:        `{"nickname":"unique_user","first_name":"Bob","last_name":"Doe","email":"taken@email.com","password":"password123","age":25,"gender":"male"}`,
			expectedStatus: http.StatusConflict,
			expectedBody:   "email already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/register", bytes.NewBufferString(tt.payload))
			if err != nil {
				t.Fatalf("Could not create HTTP request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.Register)
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("[%s] Expected status %d, but got %d. Response Body: %s",
					tt.name, tt.expectedStatus, rr.Code, rr.Body.String())
			}

			// Check body content if specified
			if tt.expectedBody != "" && !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("[%s] Expected body to contain %q, got %q",
					tt.name, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
