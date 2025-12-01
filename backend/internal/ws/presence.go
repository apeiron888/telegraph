package ws

import (
	"sync"
	"time"
)

// Presence tracks user online/offline status
type Presence struct {
	UserID   string
	Online   bool
	LastSeen time.Time
}

// PresenceTracker manages user presence
type PresenceTracker struct {
	presence map[string]*Presence
	mu       sync.RWMutex
}

func NewPresenceTracker() *PresenceTracker {
	return &PresenceTracker{
		presence: make(map[string]*Presence),
	}
}

func (pt *PresenceTracker) SetOnline(userID string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.presence[userID] = &Presence{
		UserID:   userID,
		Online:   true,
		LastSeen: time.Now(),
	}
}

func (pt *PresenceTracker) SetOffline(userID string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if p, exists := pt.presence[userID]; exists {
		p.Online = false
		p.LastSeen = time.Now()
	}
}

func (pt *PresenceTracker) GetPresence(userID string) *Presence {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	if p, exists := pt.presence[userID]; exists {
		return p
	}

	return &Presence{
		UserID:   userID,
		Online:   false,
		LastSeen: time.Time{},
	}
}

func (pt *PresenceTracker) GetAllOnline() []string {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	var online []string
	for userID, p := range pt.presence {
		if p.Online {
			online = append(online, userID)
		}
	}

	return online
}
