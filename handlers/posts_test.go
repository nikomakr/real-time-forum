package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
)

func setupPostsDB(t *testing.T) {
	mockDB, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
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
			id TEXT PRIMARY KEY, nickname TEXT UNIQUE,
			email TEXT UNIQUE, password_hash TEXT,
			first_name TEXT, last_name TEXT,
			age INTEGER, gender TEXT
		);
	`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	_, err = mockDB.Exec(`
		CREATE TABLE category_groups (
			id TEXT PRIMARY KEY, name TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		t.Fatalf("failed to create category_groups table: %v", err)
	}

	_, err = mockDB.Exec(`
		CREATE TABLE categories (
			id TEXT PRIMARY KEY, group_id TEXT NOT NULL,
			name TEXT NOT NULL, sub_group TEXT
		);
	`)
	if err != nil {
		t.Fatalf("failed to create categories table: %v", err)
	}

	_, err = mockDB.Exec(`
		CREATE TABLE posts (
			id TEXT PRIMARY KEY, author_id TEXT NOT NULL,
			title TEXT NOT NULL, content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("failed to create posts table: %v", err)
	}

	_, err = mockDB.Exec(`
		CREATE TABLE post_categories (
			post_id TEXT NOT NULL, category_id TEXT NOT NULL,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
		);
	`)
	if err != nil {
		t.Fatalf("failed to create post_categories table: %v", err)
	}

	_, err = mockDB.Exec(`
		CREATE TABLE comments (
			id TEXT PRIMARY KEY, post_id TEXT NOT NULL,
			author_id TEXT NOT NULL, content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("failed to create comments table: %v", err)
	}

	_, err = mockDB.Exec(`
		CREATE TABLE sessions (
			session_id TEXT PRIMARY KEY, user_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create sessions table: %v", err)
	}

	// Seed data
	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)

	mockDB.Exec(`INSERT INTO category_groups (id, name) VALUES ('grp-1', 'Country')`)

	mockDB.Exec(`INSERT INTO categories (id, group_id, name) VALUES ('cat-1', 'grp-1', 'England')`)

	mockDB.Exec(`INSERT INTO posts (id, author_id, title, content, created_at)
		VALUES ('p-1', 'u-1', 'Test Post', 'Test content', ?)`,
		time.Now().UTC())

	mockDB.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES ('p-1', 'cat-1')`)

	mockDB.Exec(`INSERT INTO comments (id, post_id, author_id, content)
		VALUES ('c-1', 'p-1', 'u-1', 'Test comment')`)

	mockDB.Exec(
		`INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)`,
		"test-session-id", "u-1", time.Now().UTC().Add(24*time.Hour),
	)

	db.DB = mockDB
}

func TestGetPost(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Returns post by valid ID",
			url:            "/api/posts/p-1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Returns 404 for unknown ID",
			url:            "/api/posts/unknown-id",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Rejects non-GET method",
			url:            "/api/posts/p-1",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupPostsDB(t)

			method := http.MethodGet
			if tt.name == "Rejects non-GET method" {
				method = http.MethodPost
			}

			req, _ := http.NewRequest(method, tt.url, nil)
			rr := httptest.NewRecorder()
			http.HandlerFunc(handlers.GetPost).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s",
					tt.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestGetPosts_RejectsNonGet(t *testing.T) {
	setupPostsDB(t)

	req, _ := http.NewRequest(http.MethodPost, "/api/posts", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.GetPosts).ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestGetPosts_IncludesCommentCount(t *testing.T) {
	setupPostsDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/posts", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(handlers.GetPosts).ServeHTTP(rr, req)

	var posts []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &posts)

	if len(posts) == 0 {
		t.Fatal("expected at least one post")
	}

	count, ok := posts[0]["comment_count"]
	if !ok {
		t.Error("expected comment_count field in response")
	}
	if count.(float64) != 1 {
		t.Errorf("expected comment_count 1, got %v", count)
	}
}

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Reject invalid method",
			payload:        `{}`,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Reject missing title",
			payload:        `{"content":"some content","categories":["cat-1"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "title is required",
		},
		{
			name:           "Reject missing content",
			payload:        `{"title":"Test","categories":["cat-1"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "content is required",
		},
		{
			name:           "Reject missing categories",
			payload:        `{"title":"Test","content":"content","categories":[]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "at least one category is required",
		},
		{
			name:           "Reject invalid category ID",
			payload:        `{"title":"Test","content":"content","categories":["invalid-cat"]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid category",
		},
		{
			name:           "Create post successfully",
			payload:        `{"title":"Test Post","content":"Test content","categories":["cat-1"]}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "post created successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupPostsDB(t)

			method := http.MethodPost
			if tt.name == "Reject invalid method" {
				method = http.MethodGet
			}

			req, _ := http.NewRequest(method, "/api/posts", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			// Authenticate via session cookie
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})

			rr := httptest.NewRecorder()
			handlers.RequireAuth(handlers.CreatePost).ServeHTTP(rr, req)

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
