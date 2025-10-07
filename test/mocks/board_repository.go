package mocks

import (
	"errors"
	"sync"

	"github.com/leinonen/bbs/domain"
)

type BoardRepository struct {
	mu     sync.RWMutex
	boards map[int]*domain.Board
	nextID int
}

func NewBoardRepository() *BoardRepository {
	return &BoardRepository{
		boards: make(map[int]*domain.Board),
		nextID: 1,
	}
}

func (r *BoardRepository) Create(board *domain.Board) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate name
	for _, existingBoard := range r.boards {
		if existingBoard.Name == board.Name {
			return errors.New("board name already exists")
		}
	}

	board.ID = r.nextID
	r.nextID++
	r.boards[board.ID] = board
	return nil
}

func (r *BoardRepository) GetByID(id int) (*domain.Board, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	board, exists := r.boards[id]
	if !exists {
		return nil, errors.New("board not found")
	}
	return board, nil
}

func (r *BoardRepository) GetAll() ([]*domain.Board, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	boards := make([]*domain.Board, 0, len(r.boards))
	for _, board := range r.boards {
		boards = append(boards, board)
	}
	return boards, nil
}

func (r *BoardRepository) Update(board *domain.Board) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.boards[board.ID]; !exists {
		return errors.New("board not found")
	}
	r.boards[board.ID] = board
	return nil
}

func (r *BoardRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.boards[id]; !exists {
		return errors.New("board not found")
	}
	delete(r.boards, id)
	return nil
}
