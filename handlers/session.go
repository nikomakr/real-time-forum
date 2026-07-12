package handlers

import (
	"crypto/rand" // For generating secure random bytes for session IDs
	"encoding/hex" // For encoding the random bytes into a hexadecimal string
	"net/http"
	"time"

	"real-time-forum/db"
)

const sessionDuration = 24 * time.Hour

func createSession(w http.ResponseWriter, userID string) error {
	bytes := make([]byte, 32) // 32 bytes = 256 bits, providing a strong level of entropy for the session ID. This length is sufficient to prevent brute-force attacks and ensure that session IDs are unique and hard to guess.
	if _, err := rand.Read(bytes); err != nil { // Use crypto/rand for cryptographically secure random number generation. This is crucial for generating session IDs that are unpredictable and resistant to attacks. The rand.Read function fills the byte slice with random data.
		return err
	}
	sessionID := hex.EncodeToString(bytes) // Convert the random bytes to a hexadecimal string representation. This makes the session ID safe to use in URLs and cookies, as it consists of alphanumeric characters only.

	// Everything in UTC — no local timezone desynchronisation
	expiresAt := time.Now().UTC().Add(sessionDuration)

	// Pass time.Time directly — driver handles the format correctly
	_, err := db.DB.Exec(
		`INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)`,
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiresAt, // UTC — matches database
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return nil
}