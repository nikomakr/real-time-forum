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

func setupProfileDB(t *testing.T) {
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
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	mockDB.Exec(`CREATE TABLE post_likes (
		post_id TEXT NOT NULL, user_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (post_id, user_id)
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)
	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-2', 'bob', 'bob@test.com', 'hash', 'Bob', 'Jones', 28, 'Male')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	mockDB.Exec(`INSERT INTO posts (id, author_id, title, content) VALUES ('p-1', 'u-1', 'Post 1', 'Content 1')`)
	mockDB.Exec(`INSERT INTO posts (id, author_id, title, content) VALUES ('p-2', 'u-1', 'Post 2', 'Content 2')`)
	mockDB.Exec(`INSERT INTO posts (id, author_id, title, content) VALUES ('p-3', 'u-2', 'Post 3', 'Content 3')`)

	mockDB.Exec(`INSERT INTO comments (id, post_id, author_id, content) VALUES ('c-1', 'p-3', 'u-1', 'A comment')`)

	mockDB.Exec(`INSERT INTO post_likes (post_id, user_id) VALUES ('p-3', 'u-1')`)

	db.DB = mockDB
}

func TestGetProfile(t *testing.T) {
	setupProfileDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetProfile).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var profile struct {
		UserID       string `json:"user_id"`
		Nickname     string `json:"nickname"`
		Email        string `json:"email"`
		PostCount    int    `json:"post_count"`
		CommentCount int    `json:"comment_count"`
		LikeCount    int    `json:"like_count"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &profile); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if profile.Nickname != "alice" {
		t.Errorf("expected nickname 'alice', got %q", profile.Nickname)
	}
	if profile.Email != "alice@test.com" {
		t.Errorf("expected email 'alice@test.com', got %q", profile.Email)
	}
	if profile.PostCount != 2 {
		t.Errorf("expected post_count 2, got %d", profile.PostCount)
	}
	if profile.CommentCount != 1 {
		t.Errorf("expected comment_count 1, got %d", profile.CommentCount)
	}
	if profile.LikeCount != 1 {
		t.Errorf("expected like_count 1, got %d", profile.LikeCount)
	}
}

func TestGetProfile_RejectsNonGet(t *testing.T) {
	setupProfileDB(t)

	req, _ := http.NewRequest(http.MethodPost, "/api/profile", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetProfile).ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestGetProfile_RejectsUnauthenticated(t *testing.T) {
	setupProfileDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/profile", nil)
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetProfile).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
