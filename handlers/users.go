package handlers

import (
	"log"
	"net/http"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type userResponse struct {
	ID            string  `json:"id"`
	Nickname      string  `json:"nickname"`
	LastMessageAt *string `json:"last_message_at"`
	IsOnline      bool    `json:"is_online"`
}

func GetUsers(hub interface{ OnlineUserIDs() []string }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		currentUserID := GetUserID(r)
		if currentUserID == "" {
			utils.WriteError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		// Get online user IDs from hub
		onlineIDs := hub.OnlineUserIDs()
		onlineMap := make(map[string]bool, len(onlineIDs))
		for _, id := range onlineIDs {
			onlineMap[id] = true
		}

		// FIXED QUERY: Isolates the last message timestamp strictly to the
		// dialogue shared between the target user and the current active user.
		rows, err := db.DB.Query(`
			SELECT
				u.id,
				u.nickname,
				MAX(m.created_at) AS last_message_at
			FROM users u
			LEFT JOIN messages m ON 
				(m.sender_id = u.id AND m.receiver_id = ?) OR 
				(m.sender_id = ? AND m.receiver_id = u.id)
			WHERE u.id != ?
			GROUP BY u.id, u.nickname
			ORDER BY last_message_at DESC, u.nickname ASC
		`, currentUserID, currentUserID, currentUserID)
		if err != nil {
			log.Printf("[ERROR] [GetUsers DB Query]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "could not fetch users")
			return
		}
		defer rows.Close()

		// Initialized safely to prevent 'null' rendering in JSON
		users := make([]userResponse, 0)
		for rows.Next() {
			var u userResponse
			if err := rows.Scan(
				&u.ID,
				&u.Nickname,
				&u.LastMessageAt,
			); err != nil {
				log.Printf("[ERROR] [GetUsers Scan]: %v", err)
				utils.WriteError(w, http.StatusInternalServerError, "error processing users")
				return
			}
			u.IsOnline = onlineMap[u.ID]
			users = append(users, u)
		}

		if err := rows.Err(); err != nil {
			log.Printf("[ERROR] [GetUsers Rows Loop]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "error reading users stream")
			return
		}

		if err := utils.WriteJSON(w, http.StatusOK, users); err != nil {
			log.Printf("[ERROR] [GetUsers Response JSON]: %v", err)
		}
	}
}
