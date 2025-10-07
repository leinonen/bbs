package domain

import "time"

type Post struct {
	ID        int
	BoardID   int
	UserID    int
	Username  string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	ReplyTo   *int // nil if not a reply
	Replies   int
}

func NewPost(boardID, userID int, username, title, content string) *Post {
	now := time.Now()
	return &Post{
		BoardID:   boardID,
		UserID:    userID,
		Username:  username,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		Replies:   0,
	}
}

func NewReply(boardID, userID int, username, content string, replyTo int) *Post {
	now := time.Now()
	return &Post{
		BoardID:   boardID,
		UserID:    userID,
		Username:  username,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
		ReplyTo:   &replyTo,
		Replies:   0,
	}
}