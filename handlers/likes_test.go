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

func setupLikesDB(t *testing.T) {
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

	mockDB.Exec(`CREATE TABLE post_likes (
		post_id TEXT NOT NULL, user_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	mockDB.Exec(`INSERT INTO posts (id, author_id, title, content)
		VALUES ('p-1', 'u-1', 'Test Post', 'Test content')`)

	db.DB = mockDB
}

func TestToggleLike(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Rejects non-POST method",
			method:         http.MethodGet,
			url:            "/api/posts/p-1/like",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed",
		},
		{
			name:           "Rejects missing post ID",
			method:         http.MethodPost,
			url:            "/api/posts//like",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "post ID is required",
		},
		{
			name:           "Rejects unknown post",
			method:         http.MethodPost,
			url:            "/api/posts/unknown/like",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "post not found",
		},
		{
			name:           "Likes a post successfully",
			method:         http.MethodPost,
			url:            "/api/posts/p-1/like",
			expectedStatus: http.StatusOK,
			expectedBody:   `"liked":true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupLikesDB(t)

			req, _ := http.NewRequest(tt.method, tt.url, nil)
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
			rr := httptest.NewRecorder()
			handlers.RequireAuth(handlers.ToggleLike).ServeHTTP(rr, req)

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

func TestToggleLikeTwiceUnlikes(t *testing.T) {
	setupLikesDB(t)

	req, _ := http.NewRequest(http.MethodPost, "/api/posts/p-1/like", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.ToggleLike).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 on first like, got %d. Body: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"liked":true`) || !strings.Contains(rr.Body.String(), `"like_count":1`) {
		t.Fatalf("expected first toggle to like the post with count 1, got %s", rr.Body.String())
	}

	req2, _ := http.NewRequest(http.MethodPost, "/api/posts/p-1/like", nil)
	req2.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr2 := httptest.NewRecorder()
	handlers.RequireAuth(handlers.ToggleLike).ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200 on second like, got %d. Body: %s", rr2.Code, rr2.Body.String())
	}
	if !strings.Contains(rr2.Body.String(), `"liked":false`) || !strings.Contains(rr2.Body.String(), `"like_count":0`) {
		t.Fatalf("expected second toggle to unlike the post with count 0, got %s", rr2.Body.String())
	}
}
