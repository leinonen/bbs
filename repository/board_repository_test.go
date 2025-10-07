package repository

import (
	"testing"

	"github.com/leinonen/bbs/domain"
	"github.com/leinonen/bbs/test/mocks"
)

func TestBoardRepository_Create(t *testing.T) {
	repo := mocks.NewBoardRepository()
	board := domain.NewBoard("General", "General discussion")

	err := repo.Create(board)
	if err != nil {
		t.Errorf("Create should not return error: %v", err)
	}

	if board.ID == 0 {
		t.Error("Create should set board ID")
	}

	// Test duplicate name
	board2 := domain.NewBoard("General", "Another general board")
	err = repo.Create(board2)
	if err == nil {
		t.Error("Create should return error for duplicate board name")
	}
}

func TestBoardRepository_GetByID(t *testing.T) {
	repo := mocks.NewBoardRepository()
	board := domain.NewBoard("Tech", "Technology discussions")

	// Test getting non-existent board
	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("GetByID should return error for non-existent board")
	}

	// Create board and test retrieval
	repo.Create(board)
	retrievedBoard, err := repo.GetByID(board.ID)
	if err != nil {
		t.Errorf("GetByID should not return error: %v", err)
	}

	if retrievedBoard.Name != board.Name {
		t.Errorf("Expected name %s, got %s", board.Name, retrievedBoard.Name)
	}

	if retrievedBoard.Description != board.Description {
		t.Errorf("Expected description %s, got %s", board.Description, retrievedBoard.Description)
	}
}

func TestBoardRepository_GetAll(t *testing.T) {
	repo := mocks.NewBoardRepository()

	// Test empty repository
	boards, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll should not return error: %v", err)
	}

	if len(boards) != 0 {
		t.Errorf("Expected 0 boards, got %d", len(boards))
	}

	// Create multiple boards
	board1 := domain.NewBoard("General", "General discussion")
	board2 := domain.NewBoard("Tech", "Technology discussions")
	board3 := domain.NewBoard("Random", "Random topics")

	repo.Create(board1)
	repo.Create(board2)
	repo.Create(board3)

	boards, err = repo.GetAll()
	if err != nil {
		t.Errorf("GetAll should not return error: %v", err)
	}

	if len(boards) != 3 {
		t.Errorf("Expected 3 boards, got %d", len(boards))
	}

	// Verify all boards are present
	boardNames := make(map[string]bool)
	for _, board := range boards {
		boardNames[board.Name] = true
	}

	expectedNames := []string{"General", "Tech", "Random"}
	for _, name := range expectedNames {
		if !boardNames[name] {
			t.Errorf("Expected board %s not found in results", name)
		}
	}
}

func TestBoardRepository_Update(t *testing.T) {
	repo := mocks.NewBoardRepository()
	board := domain.NewBoard("General", "General discussion")

	// Test updating non-existent board
	err := repo.Update(board)
	if err == nil {
		t.Error("Update should return error for non-existent board")
	}

	// Create board and test update
	repo.Create(board)
	board.Name = "Updated General"
	board.Description = "Updated description"

	err = repo.Update(board)
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}

	retrievedBoard, _ := repo.GetByID(board.ID)
	if retrievedBoard.Name != "Updated General" {
		t.Errorf("Expected updated name, got %s", retrievedBoard.Name)
	}

	if retrievedBoard.Description != "Updated description" {
		t.Errorf("Expected updated description, got %s", retrievedBoard.Description)
	}
}

func TestBoardRepository_Delete(t *testing.T) {
	repo := mocks.NewBoardRepository()
	board := domain.NewBoard("General", "General discussion")

	// Test deleting non-existent board
	err := repo.Delete(999)
	if err == nil {
		t.Error("Delete should return error for non-existent board")
	}

	// Create board and test deletion
	repo.Create(board)
	err = repo.Delete(board.ID)
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}

	// Verify board is deleted
	_, err = repo.GetByID(board.ID)
	if err == nil {
		t.Error("GetByID should return error for deleted board")
	}

	// Verify it's not in GetAll results
	boards, _ := repo.GetAll()
	for _, b := range boards {
		if b.ID == board.ID {
			t.Error("Deleted board should not appear in GetAll results")
		}
	}
}
