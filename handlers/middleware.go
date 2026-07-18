package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"real-time-forum/db"
	"real-time-forum/utils"
)

// ContextKey is a custom type to avoid collisions with other packages using context values
// WHY: Using a custom type for context keys helps prevent collisions with other packages that might use the same key names. By defining a unique type, we ensure that our context values are distinct and won't accidentally overwrite or be overwritten by values from other packages.
type contextKey string

const contextUserID contextKey = "user_id"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		// Look up the session in the database. I have used QueryRow as I have found it optimises connection performance by fetching only the very first matching row and automatically closing the database connection. So, that's in line with the best practices of database connection management and I kmow the way I approached it in my code is inline with it. I have session_id in sessions table as a primary key, so it will always return a single row!
		var userID string
		var expiresAt time.Time
		err = db.DB.QueryRow(
			`SELECT user_id, expires_at FROM sessions WHERE session_id = ?`,
			cookie.Value,
		).Scan(&userID, &expiresAt)

		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid session")
			return
		}

		// Check session has not expired
		if time.Now().UTC().After(expiresAt) {
			// Clean up expired session from database
			db.DB.Exec(`DELETE FROM sessions WHERE session_id = ?`, cookie.Value)
			utils.WriteError(w, http.StatusUnauthorized, "session expired")
			return
		}

		// Session sliding — extend session if less than 12 hours remaining
		if time.Until(expiresAt) < 12*time.Hour {
			newExpiry := time.Now().UTC().Add(sessionDuration)
			if _, err := db.DB.Exec(
				`UPDATE sessions SET expires_at = ? WHERE session_id = ?`,
				newExpiry, cookie.Value,
			); err != nil {
				log.Printf("[WARN] [Session Slide DB Fault]: %v", err)
			} else {
				// Only update the cookie if the DB update succeeded
				// keeps browser expiry in sync with the database
				// HTTPS only to prevent cookie highjacking over insecure connections
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    cookie.Value,
					Expires:  newExpiry,
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
			}
		}

		// Attach user ID to request context so handlers can access it
		ctx := context.WithValue(r.Context(), contextUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserID retrieves the authenticated user ID from the request context
func GetUserID(r *http.Request) string {
	userID, _ := r.Context().Value(contextUserID).(string)
	return userID
}
