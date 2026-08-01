package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

// MessageType identifies the kind of WebSocket message
type MessageType string

const (
	TypeChatMessage  MessageType = "chat_message"
	TypePresence     MessageType = "presence"
	TypeNotification MessageType = "notification"
)

// Envelope is the standard wrapper for all WebSocket messages
type Envelope struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// ChatPayload carries a private message between two users
type ChatPayload struct {
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

// PresencePayload carries the current online user list
type PresencePayload struct {
	OnlineUserIDs []string `json:"online_user_ids"`
}

// NotificationPayload carries a notification for a new message
type NotificationPayload struct {
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Preview    string `json:"preview"`
}

func buildPresenceMessage(ids []string) []byte {
	msg, _ := json.Marshal(Envelope{
		Type:    TypePresence,
		Payload: PresencePayload{OnlineUserIDs: ids},
	})
	return msg
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump(conn *websocket.Conn, onMessage func([]byte)) {
	defer func() {
		c.Hub.Unregister(c)
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("[WARN] [WS ReadPump]: %v", err)
			}
			break
		}
		if onMessage != nil {
			go onMessage(message)
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Flush any queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
