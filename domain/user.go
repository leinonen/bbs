package domain

import "time"

type User struct {
	ID        int
	Username  string
	Password  string // This will be hashed
	Email     string
	CreatedAt time.Time
	LastLogin time.Time
	IsAdmin   bool
}

func NewUser(username, email string) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Email:     email,
		CreatedAt: now,
		LastLogin: now,
		IsAdmin:   false,
	}
}
