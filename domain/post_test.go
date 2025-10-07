package domain

import (
	"testing"
	"time"
)

func TestNewPost(t *testing.T) {
	boardID := 1
	userID := 42
	username := "testuser"
	title := "Test Post"
	content := "This is a test post content"

	post := NewPost(boardID, userID, username, title, content)

	if post.BoardID != boardID {
		t.Errorf("Expected BoardID %d, got %d", boardID, post.BoardID)
	}

	if post.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, post.UserID)
	}

	if post.Username != username {
		t.Errorf("Expected Username %s, got %s", username, post.Username)
	}

	if post.Title != title {
		t.Errorf("Expected Title %s, got %s", title, post.Title)
	}

	if post.Content != content {
		t.Errorf("Expected Content %s, got %s", content, post.Content)
	}

	if post.ReplyTo != nil {
		t.Error("Expected new post ReplyTo to be nil")
	}

	if post.Replies != 0 {
		t.Errorf("Expected new post Replies to be 0, got %d", post.Replies)
	}

	if post.ID != 0 {
		t.Errorf("Expected new post ID to be 0, got %d", post.ID)
	}

	// Check that timestamps are set and recent
	now := time.Now()
	if post.CreatedAt.After(now) || post.CreatedAt.Before(now.Add(-time.Second)) {
		t.Error("CreatedAt timestamp should be recent")
	}

	if post.UpdatedAt.After(now) || post.UpdatedAt.Before(now.Add(-time.Second)) {
		t.Error("UpdatedAt timestamp should be recent")
	}

	// CreatedAt and UpdatedAt should be equal for new posts
	if !post.CreatedAt.Equal(post.UpdatedAt) {
		t.Error("CreatedAt and UpdatedAt should be equal for new posts")
	}
}

func TestNewReply(t *testing.T) {
	boardID := 1
	userID := 42
	username := "testuser"
	content := "This is a reply"
	replyTo := 10

	reply := NewReply(boardID, userID, username, content, replyTo)

	if reply.BoardID != boardID {
		t.Errorf("Expected BoardID %d, got %d", boardID, reply.BoardID)
	}

	if reply.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, reply.UserID)
	}

	if reply.Username != username {
		t.Errorf("Expected Username %s, got %s", username, reply.Username)
	}

	if reply.Content != content {
		t.Errorf("Expected Content %s, got %s", content, reply.Content)
	}

	if reply.Title != "" {
		t.Errorf("Expected reply Title to be empty, got %s", reply.Title)
	}

	if reply.ReplyTo == nil {
		t.Error("Expected reply ReplyTo to be set")
	} else if *reply.ReplyTo != replyTo {
		t.Errorf("Expected ReplyTo %d, got %d", replyTo, *reply.ReplyTo)
	}

	if reply.Replies != 0 {
		t.Errorf("Expected new reply Replies to be 0, got %d", reply.Replies)
	}

	// Check that timestamps are set and recent
	now := time.Now()
	if reply.CreatedAt.After(now) || reply.CreatedAt.Before(now.Add(-time.Second)) {
		t.Error("CreatedAt timestamp should be recent")
	}

	if reply.UpdatedAt.After(now) || reply.UpdatedAt.Before(now.Add(-time.Second)) {
		t.Error("UpdatedAt timestamp should be recent")
	}
}

func TestPostReplyToPointer(t *testing.T) {
	// Test that ReplyTo properly handles pointer semantics
	post := &Post{
		ID:      1,
		ReplyTo: nil,
	}

	if post.ReplyTo != nil {
		t.Error("ReplyTo should be nil")
	}

	replyToValue := 42
	post.ReplyTo = &replyToValue

	if post.ReplyTo == nil {
		t.Error("ReplyTo should not be nil after assignment")
	}

	if *post.ReplyTo != 42 {
		t.Errorf("Expected ReplyTo value 42, got %d", *post.ReplyTo)
	}
}
