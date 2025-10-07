package domain

import "time"

type Board struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	PostCount   int
}

func NewBoard(name, description string) *Board {
	return &Board{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		PostCount:   0,
	}
}