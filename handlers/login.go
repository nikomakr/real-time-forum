package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type loginPayload struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type loginResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // Limit the size of the request body to 1MB to prevent denial-of-service attacks and resource exhaustion. This ensures that the server does not process excessively large requests, which could lead to performance degradation or crashes.
	defer r.Body.Close()                               // Ensure the request body is closed after processing to free up resources because the request body is a stream that needs to be closed after reading. This prevents resource leaks and ensures proper cleanup of resources associated with the request.

	var payload loginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			utils.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		utils.WriteError(w, http.StatusBadRequest, "invalid request body format")
		return
	}

	if payload.Identifier == "" {
		utils.WriteError(w, http.StatusBadRequest, "nickname or email is required")
		return
	}
	if payload.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "password is required")
		return
	}

	// Single query — eliminates timing attack vector and double round-trip
	var id, passwordHash string
	err := db.DB.QueryRow(
		`SELECT id, password_hash FROM users WHERE nickname = ? OR email = ? LIMIT 1`,
		payload.Identifier, payload.Identifier,
	).Scan(&id, &passwordHash)

	/*
	   PROTECTION AGAINST USER ENUMERATION (ACCOUNT HARVESTING):

	   SCENARIO 1: The User doesn't exist.
	   If the database returns sql.ErrNoRows, we must NOT tell the client "User not found".
	   Doing so allows malicious scripts to guess and harvest valid forum emails or nicknames.
	   Instead, we deliberately return a vague "invalid credentials" error to mask the result.
	*/
	if errors.Is(err, sql.ErrNoRows) {
		utils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err != nil {
		log.Printf("[ERROR] [Login DB Lookup Fault]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	/*
	   PROTECTION AGAINST USER ENUMERATION (ACCOUNT HARVESTING):

	   SCENARIO 2: The User exists, but they typed a wrong password.
	   If the password check fails, we must NOT tell the client "Incorrect password".
	   By returning the EXACT same generic message ("invalid credentials") as Scenario 1,
	   an attacker cannot differentiate between an existing account and a non-existent one.
	*/
	if err := utils.CheckPassword(payload.Password, passwordHash); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := createSession(w, id); err != nil {
		log.Printf("[ERROR] [Login Session Creation Fault]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not create session")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, loginResponse{
		Message: "login successful",
		UserID:  id,
	}); err != nil {
		log.Printf("[ERROR] [Login Response JSON Write Fault]: %v", err)
	}
}
