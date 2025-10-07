package repository

import (
	"testing"

	"github.com/leinonen/bbs/domain"
	"github.com/leinonen/bbs/test/mocks"
)

func TestUserRepository_Create(t *testing.T) {
	repo := mocks.NewUserRepository()
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	err := repo.Create(user)
	if err != nil {
		t.Errorf("Create should not return error: %v", err)
	}

	if user.ID == 0 {
		t.Error("Create should set user ID")
	}

	// Test duplicate username
	user2 := domain.NewUser("testuser", "test2@example.com")
	err = repo.Create(user2)
	if err == nil {
		t.Error("Create should return error for duplicate username")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := mocks.NewUserRepository()
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	// Test getting non-existent user
	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("GetByID should return error for non-existent user")
	}

	// Create user and test retrieval
	repo.Create(user)
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("GetByID should not return error: %v", err)
	}

	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	repo := mocks.NewUserRepository()
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	// Test getting non-existent user
	_, err := repo.GetByUsername("nonexistent")
	if err == nil {
		t.Error("GetByUsername should return error for non-existent user")
	}

	// Create user and test retrieval
	repo.Create(user)
	retrievedUser, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Errorf("GetByUsername should not return error: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, retrievedUser.ID)
	}

	if retrievedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := mocks.NewUserRepository()
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	// Test updating non-existent user
	err := repo.Update(user)
	if err == nil {
		t.Error("Update should return error for non-existent user")
	}

	// Create user and test update
	repo.Create(user)
	user.Email = "newemail@example.com"
	user.IsAdmin = true

	err = repo.Update(user)
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}

	retrievedUser, _ := repo.GetByID(user.ID)
	if retrievedUser.Email != "newemail@example.com" {
		t.Errorf("Expected updated email, got %s", retrievedUser.Email)
	}

	if !retrievedUser.IsAdmin {
		t.Error("Expected user to be admin after update")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := mocks.NewUserRepository()
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	// Test deleting non-existent user
	err := repo.Delete(999)
	if err == nil {
		t.Error("Delete should return error for non-existent user")
	}

	// Create user and test deletion
	repo.Create(user)
	err = repo.Delete(user.ID)
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}

	// Verify user is deleted
	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Error("GetByID should return error for deleted user")
	}
}

func TestUserRepository_Authenticate(t *testing.T) {
	repo := mocks.NewUserRepository()
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	// Test authenticating non-existent user
	_, err := repo.Authenticate("nonexistent", "password")
	if err == nil {
		t.Error("Authenticate should return error for non-existent user")
	}

	// Create user and test authentication
	repo.Create(user)

	// Test correct credentials
	authUser, err := repo.Authenticate("testuser", "password123")
	if err != nil {
		t.Errorf("Authenticate should not return error for correct credentials: %v", err)
	}

	if authUser.ID != user.ID {
		t.Errorf("Expected authenticated user ID %d, got %d", user.ID, authUser.ID)
	}

	// Test incorrect password
	_, err = repo.Authenticate("testuser", "wrongpassword")
	if err == nil {
		t.Error("Authenticate should return error for incorrect password")
	}
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	repo := mocks.NewUserRepository()
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	// Test updating last login for non-existent user
	err := repo.UpdateLastLogin(999)
	if err == nil {
		t.Error("UpdateLastLogin should return error for non-existent user")
	}

	// Create user and test last login update
	repo.Create(user)
	err = repo.UpdateLastLogin(user.ID)
	if err != nil {
		t.Errorf("UpdateLastLogin should not return error: %v", err)
	}
}
