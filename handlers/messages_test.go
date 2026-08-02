package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
)

func setupMessagesDB(t *testing.T) {
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

	mockDB.Exec(`CREATE TABLE messages (
		id TEXT PRIMARY KEY,
		sender_id TEXT NOT NULL,
		receiver_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-2', 'bob', 'bob@test.com', 'hash', 'Bob', 'Jones', 35, 'Male')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	// Seed 12 messages between u-1 and u-2
	for i := 1; i <= 12; i++ {
		senderID := "u-1"
		receiverID := "u-2"
		if i%2 == 0 {
			senderID = "u-2"
			receiverID = "u-1"
		}
		mockDB.Exec(
			`INSERT INTO messages (id, sender_id, receiver_id, content, created_at)
			 VALUES (?, ?, ?, ?, ?)`,
			"m-"+strconv.Itoa(i), senderID, receiverID,
			"Message number "+strconv.Itoa(i),
			time.Now().UTC().Add(time.Duration(i)*time.Minute),
		)
	}

	db.DB = mockDB
}

func TestGetMessages_ReturnsLast10(t *testing.T) {
	setupMessagesDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/messages/u-2", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetMessages).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var messages []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &messages)

	if len(messages) != 10 {
		t.Errorf("expected 10 messages, got %d", len(messages))
	}
}

func TestGetMessages_PaginationLoads10More(t *testing.T) {
	setupMessagesDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/messages/u-2?offset=10", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetMessages).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var messages []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &messages)

	if len(messages) != 2 {
		t.Errorf("expected 2 messages at offset 10, got %d", len(messages))
	}
}

func TestGetMessages_RejectsNonGet(t *testing.T) {
	setupMessagesDB(t)

	req, _ := http.NewRequest(http.MethodPost, "/api/messages/u-2", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetMessages).ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestGetMessages_RequiresMissingUserID(t *testing.T) {
	setupMessagesDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/messages/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetMessages).ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
