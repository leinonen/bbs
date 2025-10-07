package domain

import (
	"testing"
	"time"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()

	if sm == nil {
		t.Error("NewSessionManager should return a non-nil manager")
	}

	if sm.sessions == nil {
		t.Error("SessionManager should have initialized sessions map")
	}

	if len(sm.sessions) != 0 {
		t.Error("New SessionManager should have empty sessions map")
	}
}

func TestCreateSession(t *testing.T) {
	sm := NewSessionManager()
	user := &User{
		ID:       1,
		Username: "testuser",
	}

	session := sm.CreateSession(user, nil) // nil terminal for testing

	if session == nil {
		t.Error("CreateSession should return a non-nil session")
	}

	if session.ID == "" {
		t.Error("Session should have a non-empty ID")
	}

	if session.User != user {
		t.Error("Session should reference the provided user")
	}

	// Check that timestamps are set and recent
	now := time.Now()
	if session.CreatedAt.After(now) || session.CreatedAt.Before(now.Add(-time.Second)) {
		t.Error("CreatedAt timestamp should be recent")
	}

	if session.LastActivity.After(now) || session.LastActivity.Before(now.Add(-time.Second)) {
		t.Error("LastActivity timestamp should be recent")
	}

	// Verify session is stored in manager
	storedSession, exists := sm.GetSession(session.ID)
	if !exists {
		t.Error("Session should be stored in manager")
	}

	if storedSession != session {
		t.Error("Stored session should be the same as created session")
	}
}

func TestGetSession(t *testing.T) {
	sm := NewSessionManager()
	user := &User{ID: 1, Username: "testuser"}

	// Test getting non-existent session
	_, exists := sm.GetSession("nonexistent")
	if exists {
		t.Error("GetSession should return false for non-existent session")
	}

	// Create a session and test retrieval
	session := sm.CreateSession(user, nil)
	retrievedSession, exists := sm.GetSession(session.ID)

	if !exists {
		t.Error("GetSession should return true for existing session")
	}

	if retrievedSession.ID != session.ID {
		t.Error("Retrieved session should have same ID as created session")
	}

	if retrievedSession.User.Username != user.Username {
		t.Error("Retrieved session should have same user as created session")
	}
}

func TestRemoveSession(t *testing.T) {
	sm := NewSessionManager()
	user := &User{ID: 1, Username: "testuser"}

	// Create a session
	session := sm.CreateSession(user, nil)

	// Verify it exists
	_, exists := sm.GetSession(session.ID)
	if !exists {
		t.Error("Session should exist before removal")
	}

	// Remove it
	sm.RemoveSession(session.ID)

	// Verify it's gone
	_, exists = sm.GetSession(session.ID)
	if exists {
		t.Error("Session should not exist after removal")
	}

	// Test removing non-existent session (should not panic)
	sm.RemoveSession("nonexistent")
}

func TestGetActiveSessions(t *testing.T) {
	sm := NewSessionManager()

	// Test empty sessions
	sessions := sm.GetActiveSessions()
	if len(sessions) != 0 {
		t.Error("GetActiveSessions should return empty slice for new manager")
	}

	// Create multiple sessions
	user1 := &User{ID: 1, Username: "user1"}
	user2 := &User{ID: 2, Username: "user2"}

	session1 := sm.CreateSession(user1, nil)
	session2 := sm.CreateSession(user2, nil)

	sessions = sm.GetActiveSessions()
	if len(sessions) != 2 {
		t.Errorf("Expected 2 active sessions, got %d", len(sessions))
	}

	// Verify sessions are in the list
	found1, found2 := false, false
	for _, session := range sessions {
		if session.ID == session1.ID {
			found1 = true
		}
		if session.ID == session2.ID {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Error("All created sessions should be in active sessions list")
	}

	// Remove one session and test again
	sm.RemoveSession(session1.ID)
	sessions = sm.GetActiveSessions()
	if len(sessions) != 1 {
		t.Errorf("Expected 1 active session after removal, got %d", len(sessions))
	}

	if sessions[0].ID != session2.ID {
		t.Error("Remaining session should be session2")
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1 := generateSessionID()
	id2 := generateSessionID()

	if id1 == "" {
		t.Error("Generated session ID should not be empty")
	}

	if id2 == "" {
		t.Error("Generated session ID should not be empty")
	}

	if id1 == id2 {
		t.Error("Generated session IDs should be unique")
	}

	// Test that ID has reasonable length (hex encoded, so should be even length)
	if len(id1)%2 != 0 {
		t.Error("Session ID should be even length (hex encoded)")
	}

	if len(id1) < 16 {
		t.Error("Session ID should be at least 16 characters for security")
	}
}

func TestConcurrentAccess(t *testing.T) {
	sm := NewSessionManager()
	user := &User{ID: 1, Username: "testuser"}

	// Test concurrent session creation
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			session := sm.CreateSession(user, nil)
			if session == nil {
				t.Error("Concurrent session creation failed")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	sessions := sm.GetActiveSessions()
	if len(sessions) != 10 {
		t.Errorf("Expected 10 concurrent sessions, got %d", len(sessions))
	}
}
