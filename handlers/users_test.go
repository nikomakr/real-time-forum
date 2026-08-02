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

type mockHub struct {
	onlineIDs []string
}

func (m *mockHub) OnlineUserIDs() []string {
	return m.onlineIDs
}

func setupUsersDB(t *testing.T) {
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
	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-3', 'carol', 'carol@test.com', 'hash', 'Carol', 'White', 28, 'Female')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	// u-2 has a message, u-3 does not
	mockDB.Exec(`INSERT INTO messages (id, sender_id, receiver_id, content)
		VALUES ('m-1', 'u-2', 'u-1', 'Hello')`)

	db.DB = mockDB
}

func TestGetUsers_ReturnsAllExceptCurrentUser(t *testing.T) {
	setupUsersDB(t)

	hub := &mockHub{onlineIDs: []string{"u-2"}}
	req, _ := http.NewRequest(http.MethodGet, "/api/users", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetUsers(hub)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var users []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &users)

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGetUsers_OrderedByLastMessageThenAlpha(t *testing.T) {
	setupUsersDB(t)

	hub := &mockHub{}
	req, _ := http.NewRequest(http.MethodGet, "/api/users", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetUsers(hub)).ServeHTTP(rr, req)

	var users []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &users)

	// bob has a message — should be first
	if users[0]["nickname"] != "bob" {
		t.Errorf("expected bob first (has message), got %v", users[0]["nickname"])
	}
	// carol has no message — should be second alphabetically
	if users[1]["nickname"] != "carol" {
		t.Errorf("expected carol second (no message, alphabetical), got %v", users[1]["nickname"])
	}
}

func TestGetUsers_OnlineStatusSetCorrectly(t *testing.T) {
	setupUsersDB(t)

	hub := &mockHub{onlineIDs: []string{"u-2"}}
	req, _ := http.NewRequest(http.MethodGet, "/api/users", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetUsers(hub)).ServeHTTP(rr, req)

	var users []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &users)

	for _, u := range users {
		if u["nickname"] == "bob" && u["is_online"] != true {
			t.Errorf("expected bob to be online")
		}
		if u["nickname"] == "carol" && u["is_online"] != false {
			t.Errorf("expected carol to be offline")
		}
	}
}

func TestGetUsers_RejectsNonGet(t *testing.T) {
	setupUsersDB(t)

	hub := &mockHub{}
	req, _ := http.NewRequest(http.MethodPost, "/api/users", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetUsers(hub)).ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
