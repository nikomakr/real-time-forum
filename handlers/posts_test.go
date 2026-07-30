package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
)

func setupPostsDB(t *testing.T) {
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
			PRIMARY KEY (post_id, category_id)
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

	db.DB = mockDB
}

func TestGetPosts(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Returns all posts with no filter",
			url:            "/api/posts",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "Returns filtered posts by category",
			url:            "/api/posts?category=England",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "Returns empty list for unknown category",
			url:            "/api/posts?category=Unknown",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupPostsDB(t)

			req, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			rr := httptest.NewRecorder()
			http.HandlerFunc(handlers.GetPosts).ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			var posts []map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &posts); err != nil {
				t.Fatalf("could not parse response: %v", err)
			}

			if len(posts) != tt.expectedCount {
				t.Errorf("expected %d posts, got %d", tt.expectedCount, len(posts))
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