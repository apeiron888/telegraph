package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Map of UserID to Client(s) for targeted messaging
	userClients map[string]map[*Client]bool
	userLock    sync.RWMutex
	
	// Presence tracker
	presence *PresenceTracker
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[string]map[*Client]bool),
		presence:    NewPresenceTracker(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.registerUserClient(client)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.unregisterUserClient(client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					h.unregisterUserClient(client)
				}
			}
		}
	}
}

func (h *Hub) registerUserClient(client *Client) {
	h.userLock.Lock()
	defer h.userLock.Unlock()

	if _, ok := h.userClients[client.userID]; !ok {
		h.userClients[client.userID] = make(map[*Client]bool)
	}
	h.userClients[client.userID][client] = true
	
	// Set user as online
	h.presence.SetOnline(client.userID)
	
	// Broadcast presence update
	h.BroadcastPresence(client.userID, true)
	
	log.Printf("User %s connected. Total clients: %d", client.userID, len(h.clients))
}

func (h *Hub) unregisterUserClient(client *Client) {
	h.userLock.Lock()
	defer h.userLock.Unlock()

	if clients, ok := h.userClients[client.userID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(h.userClients, client.userID)
			
			// Set user as offline if no more connections
			h.presence.SetOffline(client.userID)
			
			// Broadcast presence update
			h.BroadcastPresence(client.userID, false)
		}
	}
	log.Printf("User %s disconnected. Total clients: %d", client.userID, len(h.clients))
}

// SendToUser sends a message to a specific user's connected clients
func (h *Hub) SendToUser(userID string, message interface{}) {
	h.userLock.RLock()
	defer h.userLock.RUnlock()

	clients, ok := h.userClients[userID]
	if !ok {
		return
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.send <- msgBytes:
		default:
			// If channel is full/closed, the main loop will handle cleanup
			log.Printf("Failed to send to client for user %s", userID)
		}
	}
}

// BroadcastPresence broadcasts user online/offline status to all connected users
func (h *Hub) BroadcastPresence(userID string, online bool) {
	presenceMsg := map[string]interface{}{
		"type":    "PRESENCE",
		"user_id": userID,
		"online":  online,
	}
	
	msgBytes, _ := json.Marshal(presenceMsg)
	
	// Broadcast to all connected clients
	for client := range h.clients {
		select {
		case client.send <- msgBytes:
		default:
		}
	}
}

// BroadcastTyping broadcasts typing indicator to channel members
func (h *Hub) BroadcastTyping(userID, channelID string, typing bool) {
	typingMsg := map[string]interface{}{
		"type":       "TYPING",
		"user_id":    userID,
		"channel_id": channelID,
		"typing":     typing,
	}
	
	msgBytes, _ := json.Marshal(typingMsg)
	
	// In a production app, you'd get channel members and broadcast only to them
	// For now, broadcast to all
	for client := range h.clients {
		if client.userID != userID { // Don't send to self
			select {
			case client.send <- msgBytes:
			default:
			}
		}
	}
}

// GetPresence returns presence data for a user
func (h *Hub) GetPresence(userID string) *Presence {
	return h.presence.GetPresence(userID)
}
