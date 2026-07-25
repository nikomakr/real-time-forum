package handlers

import (
	"log"
	"net/http"
	"time"

	"real-time-forum/db"
	"real-time-forum/utils"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie present — already logged out
		utils.WriteError(w, http.StatusUnauthorized, "not logged in")
		return
	}

	// Delete the session from the database - security measure to ensure the session is invalidated server-side
	if _, err := db.DB.Exec(
		`DELETE FROM sessions WHERE session_id = ?`, cookie.Value,
	); err != nil {
		log.Printf("[ERROR] [Logout Session Delete Fault]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not log out")
		return
	}

	// We cannot force a browser to "delete" a cookie directly. Instead, we overwrite it with an expired, empty version!
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "", // Clear the cookie value aka insert an empty string
		Expires:  time.Unix(0, 0), // Expires: time.Unix(0, 0) & MaxAge: -1: Forces the browser to immediately delete the cookie because its expiration date is set to January 1, 1970.
		MaxAge:   -1,
		HttpOnly: true, // Prevents JavaScript access to the cookie, mitigating XSS attacks
		Secure:   true, // Ensures the cookie is only sent over HTTPS, protecting it from being intercepted in transit
		SameSite: http.SameSiteStrictMode, // Prevents the cookie from being sent in cross-site requests, mitigating CSRF attacks
		Path:     "/",
	})

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "logged out successfully",
	}); err != nil {
		log.Printf("[ERROR] [Logout Response JSON Write Fault]: %v", err)
	}
}