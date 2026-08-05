package handlers

import (
	"log"
	"net/http"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type profileResponse struct {
	UserID       string `json:"user_id"`
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	PostCount    int    `json:"post_count"`
	CommentCount int    `json:"comment_count"`
	LikeCount    int    `json:"like_count"`
}

// GetProfile returns the authenticated user's profile info and activity stats — GET /api/profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := GetUserID(r)
	if userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	profile := profileResponse{UserID: userID}

	if err := db.DB.QueryRow(
		`SELECT nickname, email FROM users WHERE id = ?`, userID,
	).Scan(&profile.Nickname, &profile.Email); err != nil {
		log.Printf("[ERROR] [GetProfile User Lookup]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch profile")
		return
	}

	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM posts WHERE author_id = ?`, userID,
	).Scan(&profile.PostCount); err != nil {
		log.Printf("[ERROR] [GetProfile Post Count]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch profile")
		return
	}

	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM comments WHERE author_id = ?`, userID,
	).Scan(&profile.CommentCount); err != nil {
		log.Printf("[ERROR] [GetProfile Comment Count]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch profile")
		return
	}

	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM post_likes WHERE user_id = ?`, userID,
	).Scan(&profile.LikeCount); err != nil {
		log.Printf("[ERROR] [GetProfile Like Count]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch profile")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, profile); err != nil {
		log.Printf("[ERROR] [GetProfile Response JSON]: %v", err)
	}
}
