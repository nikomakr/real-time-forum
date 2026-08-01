package ws

import "sync"

// Client represents a single connected WebSocket user
type Client struct {
	UserID string
	Send   chan []byte
	Hub    *Hub
	closed bool
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client // userID → client
}

// NewHub creates and returns a new Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

// Register adds a client to the hub and broadcasts updated presence
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c.UserID] = c
	h.mu.Unlock()

	h.BroadcastPresence()
}

// Unregister removes a client from the hub and broadcasts updated presence
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if existing, ok := h.clients[c.UserID]; ok && existing == c {
		delete(h.clients, c.UserID)
		if !c.closed {
			c.closed = true
			close(c.Send)
		}
	}
}

// SendToUser delivers a message to a specific user if they are online
func (h *Hub) SendToUser(userID string, message []byte) bool {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	select {
	case client.Send <- message:
		return true
	default:
		go h.Unregister(client)
		return false
	}
}

// OnlineUserIDs returns the list of currently connected user IDs
func (h *Hub) OnlineUserIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ids := make([]string, 0, len(h.clients))
	for id := range h.clients {
		ids = append(ids, id)
	}
	return ids
}

// BroadcastPresence sends the current online user list to all connected clients
func (h *Hub) BroadcastPresence() {
	ids := h.OnlineUserIDs()
	msg := buildPresenceMessage(ids)

	h.mu.RLock()
	activeClients := make([]*Client, 0, len(h.clients))
	for _, client := range h.clients {
		activeClients = append(activeClients, client)
	}
	h.mu.RUnlock()

	for _, client := range activeClients {
		select {
		case client.Send <- msg:
		default:
			go h.Unregister(client)
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
	h.mu.RLock()
	activeClients := make([]*Client, 0, len(h.clients))
	for _, client := range h.clients {
		activeClients = append(activeClients, client)
	}
	h.mu.RUnlock()

	for _, client := range activeClients {
		select {
		case client.Send <- message:
		default:
			go h.Unregister(client)
		}
	}
}
