package repository

import (
	"github.com/leinonen/bbs/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id int) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int) error
	Authenticate(username, password string) (*domain.User, error)
	UpdateLastLogin(userID int) error
}

type BoardRepository interface {
	Create(board *domain.Board) error
	GetByID(id int) (*domain.Board, error)
	GetAll() ([]*domain.Board, error)
	Update(board *domain.Board) error
	Delete(id int) error
}

type PostRepository interface {
	Create(post *domain.Post) error
	GetByID(id int) (*domain.Post, error)
	GetByBoard(boardID int, limit, offset int) ([]*domain.Post, error)
	GetReplies(postID int) ([]*domain.Post, error)
	GetRecent(limit int) ([]*domain.Post, error)
	Update(post *domain.Post) error
	Delete(id int) error
	CountByBoard(boardID int) (int, error)
}