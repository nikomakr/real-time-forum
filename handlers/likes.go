package handlers

import (
	"log"
	"net/http"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type likeResponse struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}

// ToggleLike likes a post for the current user if they haven't liked it yet,
// or unlikes it if they have — POST /api/posts/{id}/like
func ToggleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	postID := extractPostID(r.URL.Path)
	if postID == "" {
		utils.WriteError(w, http.StatusBadRequest, "post ID is required")
		return
	}

	var exists int
	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM posts WHERE id = ?`, postID,
	).Scan(&exists); err != nil || exists == 0 {
		utils.WriteError(w, http.StatusNotFound, "post not found")
		return
	}

	userID := GetUserID(r)
	if userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var alreadyLiked int
	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND user_id = ?`, postID, userID,
	).Scan(&alreadyLiked); err != nil {
		log.Printf("[ERROR] [ToggleLike Check]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not toggle like")
		return
	}

	liked := alreadyLiked == 0

	if liked {
		if _, err := db.DB.Exec(
			`INSERT INTO post_likes (post_id, user_id) VALUES (?, ?)`, postID, userID,
		); err != nil {
			log.Printf("[ERROR] [ToggleLike Insert]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "could not toggle like")
			return
		}
	} else {
		if _, err := db.DB.Exec(
			`DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`, postID, userID,
		); err != nil {
			log.Printf("[ERROR] [ToggleLike Delete]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "could not toggle like")
			return
		}
	}

	var count int
	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM post_likes WHERE post_id = ?`, postID,
	).Scan(&count); err != nil {
		log.Printf("[ERROR] [ToggleLike Count]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not toggle like")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, likeResponse{Liked: liked, LikeCount: count}); err != nil {
		log.Printf("[ERROR] [ToggleLike Response JSON]: %v", err)
	}
}
