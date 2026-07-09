package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/mattn/go-sqlite3"

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

	defer r.Body.Close() // Ensure the request body is closed after reading
// Limit the size of the request body to prevent abuse
	var payload registerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if payload.Nickname == "" || payload.FirstName == "" || payload.LastName == "" ||
		payload.Email == "" || payload.Password == "" ||
		payload.Age <= 0 || payload.Gender == "" {
		utils.WriteError(w, http.StatusBadRequest, "all fields are required and age must be a positive number")
		return
	}

	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not process password")
		return
	}

	id, err := utils.NewUUID()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not generate user ID")
		return
	}

	_, err = db.DB.Exec(
		`INSERT INTO users (id, nickname, first_name, last_name, email, password_hash, age, gender)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, payload.Nickname, payload.FirstName, payload.LastName,
		payload.Email, hash, payload.Age, payload.Gender,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			msg := sqliteErr.Error()
			if strings.Contains(msg, "users.nickname") || strings.Contains(msg, "\"nickname\"") {
				utils.WriteError(w, http.StatusConflict, "nickname already taken")
				return
			}
			if strings.Contains(msg, "users.email") || strings.Contains(msg, "\"email\"") {
				utils.WriteError(w, http.StatusConflict, "email already registered")
				return
			}
			utils.WriteError(w, http.StatusConflict, "nickname or email already registered")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "registration successful",
	}); err != nil {
		log.Printf("could not write response: %v", err)
	}
}
