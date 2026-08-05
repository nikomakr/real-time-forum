package handlers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"

	"real-time-forum/db"
	"real-time-forum/handlers"
	"real-time-forum/ws"
)

func setupWSDB(t *testing.T) {
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
		id TEXT PRIMARY KEY, sender_id TEXT NOT NULL, receiver_id TEXT NOT NULL,
		content TEXT NOT NULL, image_url TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)
	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-2', 'bob', 'bob@test.com', 'hash', 'Bob', 'Jones', 28, 'Male')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	db.DB = mockDB
}

// readEnvelopeOfType reads messages off conn until one of the given type
// arrives (skipping others, e.g. the presence broadcast sent on connect).
func readEnvelopeOfType(t *testing.T, conn *websocket.Conn, wantType string) ws.Envelope {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	for range 5 {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed reading websocket message: %v", err)
		}
		var env ws.Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			t.Fatalf("failed decoding envelope: %v", err)
		}
		if string(env.Type) == wantType {
			return env
		}
	}
	t.Fatalf("did not receive a %q envelope after 5 messages", wantType)
	return ws.Envelope{}
}

func TestServeWSRejectsUnauthenticated(t *testing.T) {
	setupWSDB(t)
	hub := ws.NewHub()

	srv := httptest.NewServer(handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWS(hub, w, r)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for unauthenticated upgrade attempt, got %d", resp.StatusCode)
	}
}

func TestServeWSDeliversAndPersistsChatMessage(t *testing.T) {
	setupWSDB(t)
	hub := ws.NewHub()

	srv := httptest.NewServer(handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWS(hub, w, r)
	}))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	header := http.Header{}
	header.Set("Cookie", "session_id=test-session-id")

	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		t.Fatalf("failed to dial websocket: %v (status %d)", err, resp.StatusCode)
	}
	defer conn.Close()

	outgoing := map[string]string{
		"type":        "chat_message",
		"receiver_id": "u-2",
		"content":     "hello there",
	}
	payload, _ := json.Marshal(outgoing)
	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	env := readEnvelopeOfType(t, conn, "chat_message")

	payloadBytes, _ := json.Marshal(env.Payload)
	var chat ws.ChatPayload
	if err := json.Unmarshal(payloadBytes, &chat); err != nil {
		t.Fatalf("failed decoding chat payload: %v", err)
	}
	if chat.Content != "hello there" {
		t.Errorf("expected echoed content %q, got %q", "hello there", chat.Content)
	}
	if chat.SenderID != "u-1" || chat.ReceiverID != "u-2" {
		t.Errorf("unexpected sender/receiver: %+v", chat)
	}
	if chat.SenderName != "alice" {
		t.Errorf("expected sender name 'alice', got %q", chat.SenderName)
	}

	var count int
	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM messages WHERE sender_id = ? AND receiver_id = ? AND content = ?`,
		"u-1", "u-2", "hello there",
	).Scan(&count); err != nil {
		t.Fatalf("failed to query persisted message: %v", err)
	}
	if count != 1 {
		t.Errorf("expected message to be persisted exactly once, found %d", count)
	}
}
