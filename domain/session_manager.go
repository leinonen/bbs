package domain

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"golang.org/x/term"
)

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) CreateSession(user *User, term *term.Terminal) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sessionID := generateSessionID()
	session := &Session{
		ID:           sessionID,
		User:         user,
		Terminal:     term,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	sm.sessions[sessionID] = session
	return session
}

func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[id]
	return session, exists
}

func (sm *SessionManager) RemoveSession(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, id)
}

func (sm *SessionManager) GetActiveSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}