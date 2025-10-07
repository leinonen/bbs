package domain

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	email := "test@example.com"

	user := NewUser(username, email)

	if user.Username != username {
		t.Errorf("Expected username %s, got %s", username, user.Username)
	}

	if user.Email != email {
		t.Errorf("Expected email %s, got %s", email, user.Email)
	}

	if user.IsAdmin {
		t.Error("Expected new user to not be admin")
	}

	if user.ID != 0 {
		t.Errorf("Expected new user ID to be 0, got %d", user.ID)
	}

	// Check that timestamps are set and recent
	now := time.Now()
	if user.CreatedAt.After(now) || user.CreatedAt.Before(now.Add(-time.Second)) {
		t.Error("CreatedAt timestamp should be recent")
	}

	if user.LastLogin.After(now) || user.LastLogin.Before(now.Add(-time.Second)) {
		t.Error("LastLogin timestamp should be recent")
	}
}

func TestUserFields(t *testing.T) {
	user := &User{
		ID:        42,
		Username:  "admin",
		Password:  "hashed_password",
		Email:     "admin@example.com",
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		LastLogin: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		IsAdmin:   true,
	}

	if user.ID != 42 {
		t.Errorf("Expected ID 42, got %d", user.ID)
	}

	if user.Username != "admin" {
		t.Errorf("Expected username admin, got %s", user.Username)
	}

	if !user.IsAdmin {
		t.Error("Expected user to be admin")
	}

	expectedCreated := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !user.CreatedAt.Equal(expectedCreated) {
		t.Errorf("Expected CreatedAt %v, got %v", expectedCreated, user.CreatedAt)
	}
}
