package domain

import (
	"time"

	"golang.org/x/term"
)

type Session struct {
	ID           string
	User         *User
	Terminal     *term.Terminal
	CreatedAt    time.Time
	LastActivity time.Time
}
