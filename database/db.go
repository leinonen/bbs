package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at DATETIME NOT NULL,
		last_login DATETIME NOT NULL,
		is_admin BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS boards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		board_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		title TEXT,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		reply_to INTEGER,
		FOREIGN KEY (board_id) REFERENCES boards(id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (reply_to) REFERENCES posts(id)
	);

	CREATE INDEX IF NOT EXISTS idx_posts_board ON posts(board_id);
	CREATE INDEX IF NOT EXISTS idx_posts_user ON posts(user_id);
	CREATE INDEX IF NOT EXISTS idx_posts_reply ON posts(reply_to);
	CREATE INDEX IF NOT EXISTS idx_posts_updated ON posts(updated_at);

	INSERT OR IGNORE INTO boards (id, name, description, created_at)
	VALUES
		(1, 'general', 'General discussion', datetime('now')),
		(2, 'tech', 'Technology and programming', datetime('now')),
		(3, 'random', 'Random topics', datetime('now'));
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}

	return nil
}