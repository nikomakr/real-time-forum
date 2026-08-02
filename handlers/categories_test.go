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

func setupCategoriesDB(t *testing.T) {
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

	mockDB.Exec(`CREATE TABLE category_groups (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	)`)

	mockDB.Exec(`CREATE TABLE categories (
		id TEXT PRIMARY KEY,
		group_id TEXT NOT NULL,
		name TEXT NOT NULL,
		sub_group TEXT
	)`)

	mockDB.Exec(`INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES ('u-1', 'alice', 'alice@test.com', 'hash', 'Alice', 'Smith', 30, 'Female')`)

	mockDB.Exec(`INSERT INTO sessions (session_id, user_id, expires_at)
		VALUES ('test-session-id', 'u-1', ?)`,
		time.Now().UTC().Add(24*time.Hour))

	mockDB.Exec(`INSERT INTO category_groups (id, name) VALUES ('grp-1', 'Country')`)
	mockDB.Exec(`INSERT INTO category_groups (id, name) VALUES ('grp-2', 'Benefits')`)

	mockDB.Exec(`INSERT INTO categories (id, group_id, name, sub_group) VALUES ('cat-1', 'grp-1', 'England', NULL)`)
	mockDB.Exec(`INSERT INTO categories (id, group_id, name, sub_group) VALUES ('cat-2', 'grp-1', 'Wales', NULL)`)
	mockDB.Exec(`INSERT INTO categories (id, group_id, name, sub_group) VALUES ('cat-3', 'grp-2', 'Universal Credit', NULL)`)

	db.DB = mockDB
}

func TestGetCategories_ReturnsGroupedCategories(t *testing.T) {
	setupCategoriesDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/categories", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetCategories).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var groups []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &groups)

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestGetCategories_GroupsContainCategories(t *testing.T) {
	setupCategoriesDB(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/categories", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetCategories).ServeHTTP(rr, req)

	var groups []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &groups)

	cats := groups[0]["categories"].([]interface{})
	if len(cats) != 2 {
		t.Errorf("expected 2 categories in Country group, got %d", len(cats))
	}
}

func TestGetCategories_RejectsNonGet(t *testing.T) {
	setupCategoriesDB(t)

	req, _ := http.NewRequest(http.MethodPost, "/api/categories", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session-id"})
	rr := httptest.NewRecorder()
	handlers.RequireAuth(handlers.GetCategories).ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
