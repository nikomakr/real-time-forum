package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"path"
	"strconv"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type messageResponse struct {
	ID         string  `json:"id"`
	SenderID   string  `json:"sender_id"`
	SenderName string  `json:"sender_name"`
	ReceiverID string  `json:"receiver_id"`
	Content    string  `json:"content"`
	ImageURL   *string `json:"image_url,omitempty"`
	CreatedAt  string  `json:"created_at"`
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// r.PathValue works when routed via /api/messages/{id}
	// path.Base fallback covers direct handler calls in tests
	otherUserID := r.PathValue("id")
	if otherUserID == "" {
		otherUserID = path.Base(r.URL.Path)
	}
	if otherUserID == "" || otherUserID == "messages" {
		utils.WriteError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	currentUserID := GetUserID(r)
	if currentUserID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	// Pagination — default offset 0, page size 10
	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		parsed, err := strconv.Atoi(o)
		if err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	rows, err := db.DB.Query(`
		SELECT
			m.id,
			m.sender_id,
			s.nickname,
			m.receiver_id,
			m.content,
			m.image_url,
			m.created_at
		FROM messages m
		JOIN users s ON m.sender_id = s.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?)
		   OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT 10 OFFSET ?
	`, currentUserID, otherUserID, otherUserID, currentUserID, offset)
	if err != nil {
		log.Printf("[ERROR] [GetMessages DB Query]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch messages")
		return
	}
	defer rows.Close()

	messages := []messageResponse{}
	for rows.Next() {
		var m messageResponse
		var imageURL sql.NullString
		if err := rows.Scan(
			&m.ID,
			&m.SenderID,
			&m.SenderName,
			&m.ReceiverID,
			&m.Content,
			&imageURL,
			&m.CreatedAt,
		); err != nil {
			log.Printf("[ERROR] [GetMessages Scan]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "error processing messages")
			return
		}
		if imageURL.Valid {
			m.ImageURL = &imageURL.String
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] [GetMessages Rows Loop]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "error reading messages stream")
		return
	}

	// Reverse so oldest message is first — query returns newest first for pagination
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	if err := utils.WriteJSON(w, http.StatusOK, messages); err != nil {
		log.Printf("[ERROR] [GetMessages Response JSON]: %v", err)
	}
}
