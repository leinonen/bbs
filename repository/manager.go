package repository

import (
	"database/sql"

	"github.com/leinonen/bbs/repository/sqlite"
)

type Manager struct {
	User  UserRepository
	Board BoardRepository
	Post  PostRepository
	db    *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		User:  sqlite.NewUserRepository(db),
		Board: sqlite.NewBoardRepository(db),
		Post:  sqlite.NewPostRepository(db),
		db:    db,
	}
}

func (m *Manager) DB() *sql.DB {
	return m.db
}

func (m *Manager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}
