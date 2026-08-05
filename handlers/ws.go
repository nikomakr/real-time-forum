package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"real-time-forum/db"
	"real-time-forum/utils"
	"real-time-forum/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins during development — restrict in production
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type incomingMessage struct {
	Type       string `json:"type"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	ImageURL   string `json:"image_url"`
}

func ServeWS(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	if userID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] [ServeWS Upgrade]: %v", err)
		return
	}

	client := &ws.Client{
		UserID: userID,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	hub.Register(client)

	// Run write pump in a separate goroutine
	go client.WritePump(conn)

	// Run read pump — blocks until connection closes
	go client.ReadPump(conn, func(msg []byte) {
		handleIncoming(hub, client, msg)
	})
}

func handleIncoming(hub *ws.Hub, sender *ws.Client, raw []byte) {
	var msg incomingMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		log.Printf("[WARN] [ServeWS handleIncoming parse]: %v", err)
		return
	}

	if msg.Type != "chat_message" || msg.ReceiverID == "" || (msg.Content == "" && msg.ImageURL == "") {
		return
	}

	// Only ever trust image URLs that point at our own uploads — reject
	// anything else rather than letting a raw WS client inject arbitrary URLs.
	imageURL := ""
	if msg.ImageURL != "" && strings.HasPrefix(msg.ImageURL, "/uploads/") && !strings.Contains(msg.ImageURL, "..") {
		imageURL = msg.ImageURL
	}

	// Fetch sender nickname
	var senderName string
	if err := db.DB.QueryRow(
		`SELECT nickname FROM users WHERE id = ?`, sender.UserID,
	).Scan(&senderName); err != nil {
		log.Printf("[ERROR] [ServeWS nickname lookup]: %v", err)
		return
	}

	// Persist message to database
	msgID, err := utils.NewUUID()
	if err != nil {
		log.Printf("[ERROR] [ServeWS NewUUID]: %v", err)
		return
	}

	now := time.Now().UTC()

	var imageURLColumn any
	if imageURL != "" {
		imageURLColumn = imageURL
	}

	_, err = db.DB.Exec(
		`INSERT INTO messages (id, sender_id, receiver_id, content, image_url, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		msgID, sender.UserID, msg.ReceiverID, msg.Content, imageURLColumn, now,
	)
	if err != nil {
		log.Printf("[ERROR] [ServeWS insert message]: %v", err)
		return
	}

	// Build the outbound envelope
	payload := ws.ChatPayload{
		SenderID:   sender.UserID,
		SenderName: senderName,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		ImageURL:   imageURL,
		CreatedAt:  now.Format(time.RFC3339),
	}

	envelope, err := json.Marshal(ws.Envelope{
		Type:    ws.TypeChatMessage,
		Payload: payload,
	})
	if err != nil {
		log.Printf("[ERROR] [ServeWS marshal envelope]: %v", err)
		return
	}

	// Deliver to receiver
	hub.SendToUser(msg.ReceiverID, envelope)

	// Deliver notification to receiver
	notification, _ := json.Marshal(ws.Envelope{
		Type: ws.TypeNotification,
		Payload: ws.NotificationPayload{
			SenderID:   sender.UserID,
			SenderName: senderName,
			Preview:    msg.Content,
		},
	})
	hub.SendToUser(msg.ReceiverID, notification)

	// Echo back to sender so their own UI updates
	hub.SendToUser(sender.UserID, envelope)
}
