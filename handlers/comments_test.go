package handlers_test

import (
	"bytes"
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

func setupCommentsDB(t *testing.T) {
	mockDB, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
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

	mockDB.Exec(`CREATE TABLE posts (
		id TEXT PRIMARY KEY, author_id TEXT NOT NULL,
		title TEXT NOT NULL, content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	mockDB.Exec(`CREATE TABLE comments (
		id TEXT PRIMARY KEY, post_id TEXT NOT NULL,
		author_id TEXT NOT NULL, content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	mockDB.Exec(`INSERT INTO posts (id, author_id, title, content)
		VALUES ('p-1', 'u-1', 'Test Post', 'Test content')`)

	mockDB.Exec(`INSERT INTO comments (id, post_id, author_id, content)
		VALUES ('c-1', 'p-1', 'u-1', 'Existing comment')`)

	db.DB = mockDB
}

func TestGetComments(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Returns comments for valid post",
			url:            "/api/posts/p-1/comments",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Returns 404 for unknown post",
			url:            "/api/posts/unknown/comments",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Rejects non-GET method",
			url:            "/api/posts/p-1/comments",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupCommentsDB(t)

			method := http.MethodGet
			if tt.name == "Rejects non-GET method" {
				method = http.MethodPost
			}

			req, _ := http.NewRequest(method, tt.url, nil)
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
			rr := httptest.NewRecorder()
			handlers.RequireAuth(handlers.GetComments).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s",
					tt.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestCreateComment(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Rejects invalid method",
			url:            "/api/posts/p-1/comments",
			payload:        `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Rejects missing content",
			url:            "/api/posts/p-1/comments",
			payload:        `{"content":""}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "content is required",
		},
		{
			name:           "Rejects unknown post",
			url:            "/api/posts/unknown/comments",
			payload:        `{"content":"hello"}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "post not found",
		},
		{
			name:           "Creates comment successfully",
			url:            "/api/posts/p-1/comments",
			payload:        `{"content":"This is a reply"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "comment created successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupCommentsDB(t)

			method := http.MethodPost
			if tt.name == "Rejects invalid method" {
				method = http.MethodGet
			}

			req, _ := http.NewRequest(method, tt.url, bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
			rr := httptest.NewRecorder()
			handlers.RequireAuth(handlers.CreateComment).ServeHTTP(rr, req)

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
