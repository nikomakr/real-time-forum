package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
)

func setupSessionDB(t *testing.T) {
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
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)

	db.DB = mockDB
}

func TestCreateSession(t *testing.T) {
	setupSessionDB(t)

	rr := httptest.NewRecorder()
	if err := createSession(rr, "u-1"); err != nil {
		t.Fatalf("createSession returned error: %v", err)
	}

	var cookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == "session_id" {
			cookie = c
			break
		}
	}
	if cookie == nil {
		t.Fatal("expected session_id cookie to be set, but none was found")
	}
	if cookie.Value == "" {
		t.Error("expected non-empty session ID")
	}
	if !cookie.HttpOnly {
		t.Error("expected cookie to be HttpOnly")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
	if time.Until(cookie.Expires) < 23*time.Hour || time.Until(cookie.Expires) > 25*time.Hour {
		t.Errorf("expected cookie to expire ~24h from now, got %v", cookie.Expires)
	}

	var userID string
	var expiresAt time.Time
	err := db.DB.QueryRow(
		`SELECT user_id, expires_at FROM sessions WHERE session_id = ?`, cookie.Value,
	).Scan(&userID, &expiresAt)
	if err != nil {
		t.Fatalf("expected session row to be persisted: %v", err)
	}
	if userID != "u-1" {
		t.Errorf("expected session to belong to u-1, got %q", userID)
	}
}
