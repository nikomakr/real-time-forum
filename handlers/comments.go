package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path"
	"strings"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type commentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type createCommentPayload struct {
	Content string `json:"content"`
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract post ID from /api/posts/{id}/comments
	postID := extractPostID(r.URL.Path)
	if postID == "" {
		utils.WriteError(w, http.StatusBadRequest, "post ID is required")
		return
	}

	// Verify post exists
	var exists int
	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM posts WHERE id = ?`, postID,
	).Scan(&exists); err != nil || exists == 0 {
		utils.WriteError(w, http.StatusNotFound, "post not found")
		return
	}

	rows, err := db.DB.Query(`
		SELECT
			c.id,
			c.post_id,
			u.nickname,
			c.content,
			c.created_at
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		log.Printf("[ERROR] [GetComments DB Query]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch comments")
		return
	}
	defer rows.Close()

	comments := []commentResponse{}
	for rows.Next() {
		var c commentResponse
		if err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.Author,
			&c.Content,
			&c.CreatedAt,
		); err != nil {
			log.Printf("[ERROR] [GetComments Scan]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "error processing comments")
			return
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] [GetComments Rows Loop]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "error reading comments stream")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, comments); err != nil {
		log.Printf("[ERROR] [GetComments Response JSON]: %v", err)
	}
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	postID := extractPostID(r.URL.Path)
	if postID == "" {
		utils.WriteError(w, http.StatusBadRequest, "post ID is required")
		return
	}

	// Verify post exists
	var exists int
	if err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM posts WHERE id = ?`, postID,
	).Scan(&exists); err != nil || exists == 0 {
		utils.WriteError(w, http.StatusNotFound, "post not found")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	var payload createCommentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			utils.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	payload.Content = strings.TrimSpace(payload.Content)
	if payload.Content == "" {
		utils.WriteError(w, http.StatusBadRequest, "content is required")
		return
	}

	authorID := GetUserID(r)
	if authorID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	commentID, err := utils.NewUUID()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not generate comment ID")
		return
	}

	_, err = db.DB.Exec(
		`INSERT INTO comments (id, post_id, author_id, content) VALUES (?, ?, ?, ?)`,
		commentID, postID, authorID, payload.Content,
	)
	if err != nil {
		log.Printf("[ERROR] [CreateComment Insert]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not create comment")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "comment created successfully",
		"id":      commentID,
	}); err != nil {
		log.Printf("[ERROR] [CreateComment Response JSON]: %v", err)
	}
}

// extractPostID extracts the post ID from paths like /api/posts/{id}/comments
func extractPostID(urlPath string) string {
	cleaned := path.Clean(urlPath)
	parts := strings.Split(cleaned, "/")
	// Expected: ["", "api", "posts", "{id}", "comments"]
	if len(parts) < 5 {
		return ""
	}
	id := parts[len(parts)-2]
	if id == "posts" || id == "" {
		return ""
	}
	return id
}

// postIDFromPath extracts post ID from /api/posts/{id} (no trailing segment)
func postIDFromPath(urlPath string) string {
	return path.Base(urlPath)
}
