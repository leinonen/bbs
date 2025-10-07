package domain

import (
	"testing"
	"time"
)

func TestNewBoard(t *testing.T) {
	name := "General Discussion"
	description := "A place for general chat"

	board := NewBoard(name, description)

	if board.Name != name {
		t.Errorf("Expected name %s, got %s", name, board.Name)
	}

	if board.Description != description {
		t.Errorf("Expected description %s, got %s", description, board.Description)
	}

	if board.PostCount != 0 {
		t.Errorf("Expected new board post count to be 0, got %d", board.PostCount)
	}

	if board.ID != 0 {
		t.Errorf("Expected new board ID to be 0, got %d", board.ID)
	}

	// Check that timestamp is set and recent
	now := time.Now()
	if board.CreatedAt.After(now) || board.CreatedAt.Before(now.Add(-time.Second)) {
		t.Error("CreatedAt timestamp should be recent")
	}
}

func TestBoardFields(t *testing.T) {
	board := &Board{
		ID:          1,
		Name:        "Tech Talk",
		Description: "Technology discussions",
		CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		PostCount:   15,
	}

	if board.ID != 1 {
		t.Errorf("Expected ID 1, got %d", board.ID)
	}

	if board.Name != "Tech Talk" {
		t.Errorf("Expected name 'Tech Talk', got %s", board.Name)
	}

	if board.PostCount != 15 {
		t.Errorf("Expected post count 15, got %d", board.PostCount)
	}

	expectedCreated := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !board.CreatedAt.Equal(expectedCreated) {
		t.Errorf("Expected CreatedAt %v, got %v", expectedCreated, board.CreatedAt)
	}
}
