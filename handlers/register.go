package handlers

import (
	"encoding/json"
	"net/http"
	"log"
	"real-time-forum/db"
	"real-time-forum/utils"
)

type registerPayload struct {
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var payload registerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate all required fields are present
	if payload.Nickname == "" || payload.FirstName == "" || payload.LastName == "" ||
		payload.Email == "" || payload.Password == "" ||
		payload.Age == 0 || payload.Gender == "" {
		utils.WriteError(w, http.StatusBadRequest, "all fields are required")
		return
	}

	// Check for duplicate nickname
	var count int
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?", payload.Nickname).Scan(&count)
	if count > 0 {
		utils.WriteError(w, http.StatusConflict, "nickname already taken")
		return
	}

	// Check for duplicate email
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", payload.Email).Scan(&count)
	if count > 0 {
		utils.WriteError(w, http.StatusConflict, "email already registered")
		return
	}

	// Hash password
	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not process password")
		return
	}

	// Generate UUID
	id, err := utils.NewUUID()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not generate user ID")
		return
	}

	// Insert user
	_, err = db.DB.Exec(
		`INSERT INTO users (id, nickname, first_name, last_name, email, password_hash, age, gender)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, payload.Nickname, payload.FirstName, payload.LastName,
		payload.Email, hash, payload.Age, payload.Gender,
	)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "registration successful",
	}); err != nil {
		log.Printf("could not write response: %v", err)
	}
}