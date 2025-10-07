package repository

import (
	"testing"
	"time"

	"github.com/leinonen/bbs/domain"
	"github.com/leinonen/bbs/test/mocks"
)

func TestPostRepository_Create(t *testing.T) {
	repo := mocks.NewPostRepository()
	post := domain.NewPost(1, 42, "testuser", "Test Post", "This is a test post")

	err := repo.Create(post)
	if err != nil {
		t.Errorf("Create should not return error: %v", err)
	}

	if post.ID == 0 {
		t.Error("Create should set post ID")
	}

	// Test creating a reply
	reply := domain.NewReply(1, 43, "otheruser", "This is a reply", post.ID)
	err = repo.Create(reply)
	if err != nil {
		t.Errorf("Create should not return error for reply: %v", err)
	}

	if reply.ReplyTo == nil || *reply.ReplyTo != post.ID {
		t.Error("Reply should have correct ReplyTo value")
	}
}

func TestPostRepository_GetByID(t *testing.T) {
	repo := mocks.NewPostRepository()
	post := domain.NewPost(1, 42, "testuser", "Test Post", "This is a test post")

	// Test getting non-existent post
	_, err := repo.GetByID(999)
	if err == nil {
		t.Error("GetByID should return error for non-existent post")
	}

	// Create post and test retrieval
	repo.Create(post)
	retrievedPost, err := repo.GetByID(post.ID)
	if err != nil {
		t.Errorf("GetByID should not return error: %v", err)
	}

	if retrievedPost.Title != post.Title {
		t.Errorf("Expected title %s, got %s", post.Title, retrievedPost.Title)
	}

	if retrievedPost.Content != post.Content {
		t.Errorf("Expected content %s, got %s", post.Content, retrievedPost.Content)
	}

	if retrievedPost.Username != post.Username {
		t.Errorf("Expected username %s, got %s", post.Username, retrievedPost.Username)
	}
}

func TestPostRepository_GetByBoard(t *testing.T) {
	repo := mocks.NewPostRepository()

	// Test empty board
	posts, err := repo.GetByBoard(1, 10, 0)
	if err != nil {
		t.Errorf("GetByBoard should not return error: %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("Expected 0 posts, got %d", len(posts))
	}

	// Create posts for different boards
	post1 := domain.NewPost(1, 42, "user1", "Post 1", "Content 1")
	post2 := domain.NewPost(1, 43, "user2", "Post 2", "Content 2")
	post3 := domain.NewPost(2, 44, "user3", "Post 3", "Content 3")
	reply := domain.NewReply(1, 45, "user4", "Reply content", 0) // Will be set after post1 creation

	// Adjust times to test ordering
	post1.CreatedAt = time.Now().Add(-2 * time.Hour)
	post2.CreatedAt = time.Now().Add(-1 * time.Hour)
	post3.CreatedAt = time.Now()

	repo.Create(post1)
	repo.Create(post2)
	repo.Create(post3)

	// Create reply to post1
	*reply.ReplyTo = post1.ID
	repo.Create(reply)

	// Test getting posts for board 1
	posts, err = repo.GetByBoard(1, 10, 0)
	if err != nil {
		t.Errorf("GetByBoard should not return error: %v", err)
	}

	// Should get 2 posts (not the reply) in newest-first order
	if len(posts) != 2 {
		t.Errorf("Expected 2 posts for board 1, got %d", len(posts))
	}

	// Check ordering (newest first)
	if posts[0].Title != "Post 2" {
		t.Errorf("Expected first post to be 'Post 2', got %s", posts[0].Title)
	}

	if posts[1].Title != "Post 1" {
		t.Errorf("Expected second post to be 'Post 1', got %s", posts[1].Title)
	}

	// Test pagination
	posts, err = repo.GetByBoard(1, 1, 0)
	if err != nil {
		t.Errorf("GetByBoard should not return error: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post with limit 1, got %d", len(posts))
	}

	posts, err = repo.GetByBoard(1, 1, 1)
	if err != nil {
		t.Errorf("GetByBoard should not return error: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post with offset 1, got %d", len(posts))
	}

	if posts[0].Title != "Post 1" {
		t.Errorf("Expected offset post to be 'Post 1', got %s", posts[0].Title)
	}
}

func TestPostRepository_GetReplies(t *testing.T) {
	repo := mocks.NewPostRepository()

	// Create a post
	post := domain.NewPost(1, 42, "testuser", "Test Post", "This is a test post")
	repo.Create(post)

	// Test getting replies for post with no replies
	replies, err := repo.GetReplies(post.ID)
	if err != nil {
		t.Errorf("GetReplies should not return error: %v", err)
	}

	if len(replies) != 0 {
		t.Errorf("Expected 0 replies, got %d", len(replies))
	}

	// Create replies with different times to test ordering
	reply1 := domain.NewReply(1, 43, "user2", "First reply", post.ID)
	reply2 := domain.NewReply(1, 44, "user3", "Second reply", post.ID)
	reply3 := domain.NewReply(1, 45, "user4", "Third reply", post.ID)

	reply1.CreatedAt = time.Now().Add(-2 * time.Hour)
	reply2.CreatedAt = time.Now().Add(-1 * time.Hour)
	reply3.CreatedAt = time.Now()

	repo.Create(reply1)
	repo.Create(reply2)
	repo.Create(reply3)

	// Test getting replies
	replies, err = repo.GetReplies(post.ID)
	if err != nil {
		t.Errorf("GetReplies should not return error: %v", err)
	}

	if len(replies) != 3 {
		t.Errorf("Expected 3 replies, got %d", len(replies))
	}

	// Check ordering (oldest first for replies)
	if replies[0].Content != "First reply" {
		t.Errorf("Expected first reply to be 'First reply', got %s", replies[0].Content)
	}

	if replies[2].Content != "Third reply" {
		t.Errorf("Expected last reply to be 'Third reply', got %s", replies[2].Content)
	}

	// Test getting replies for non-existent post
	replies, err = repo.GetReplies(999)
	if err != nil {
		t.Errorf("GetReplies should not return error for non-existent post: %v", err)
	}

	if len(replies) != 0 {
		t.Errorf("Expected 0 replies for non-existent post, got %d", len(replies))
	}
}

func TestPostRepository_GetRecent(t *testing.T) {
	repo := mocks.NewPostRepository()

	// Test empty repository
	posts, err := repo.GetRecent(10)
	if err != nil {
		t.Errorf("GetRecent should not return error: %v", err)
	}

	if len(posts) != 0 {
		t.Errorf("Expected 0 recent posts, got %d", len(posts))
	}

	// Create posts and replies with different times
	post1 := domain.NewPost(1, 42, "user1", "Old Post", "Old content")
	post2 := domain.NewPost(1, 43, "user2", "Recent Post", "Recent content")
	reply := domain.NewReply(1, 44, "user3", "Recent reply", 0)

	post1.CreatedAt = time.Now().Add(-2 * time.Hour)
	post2.CreatedAt = time.Now().Add(-1 * time.Hour)
	reply.CreatedAt = time.Now()

	repo.Create(post1)
	repo.Create(post2)
	*reply.ReplyTo = post1.ID
	repo.Create(reply)

	// Test getting recent posts
	posts, err = repo.GetRecent(10)
	if err != nil {
		t.Errorf("GetRecent should not return error: %v", err)
	}

	if len(posts) != 3 {
		t.Errorf("Expected 3 recent posts (including reply), got %d", len(posts))
	}

	// Check ordering (newest first)
	if posts[0].Content != "Recent reply" {
		t.Errorf("Expected first post to be recent reply, got %s", posts[0].Content)
	}

	// Test limit
	posts, err = repo.GetRecent(2)
	if err != nil {
		t.Errorf("GetRecent should not return error: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 recent posts with limit 2, got %d", len(posts))
	}
}

func TestPostRepository_Update(t *testing.T) {
	repo := mocks.NewPostRepository()
	post := domain.NewPost(1, 42, "testuser", "Test Post", "Original content")

	// Test updating non-existent post
	err := repo.Update(post)
	if err == nil {
		t.Error("Update should return error for non-existent post")
	}

	// Create post and test update
	repo.Create(post)
	post.Title = "Updated Title"
	post.Content = "Updated content"
	post.UpdatedAt = time.Now()

	err = repo.Update(post)
	if err != nil {
		t.Errorf("Update should not return error: %v", err)
	}

	retrievedPost, _ := repo.GetByID(post.ID)
	if retrievedPost.Title != "Updated Title" {
		t.Errorf("Expected updated title, got %s", retrievedPost.Title)
	}

	if retrievedPost.Content != "Updated content" {
		t.Errorf("Expected updated content, got %s", retrievedPost.Content)
	}
}

func TestPostRepository_Delete(t *testing.T) {
	repo := mocks.NewPostRepository()
	post := domain.NewPost(1, 42, "testuser", "Test Post", "Test content")

	// Test deleting non-existent post
	err := repo.Delete(999)
	if err == nil {
		t.Error("Delete should return error for non-existent post")
	}

	// Create post and test deletion
	repo.Create(post)
	err = repo.Delete(post.ID)
	if err != nil {
		t.Errorf("Delete should not return error: %v", err)
	}

	// Verify post is deleted
	_, err = repo.GetByID(post.ID)
	if err == nil {
		t.Error("GetByID should return error for deleted post")
	}
}

func TestPostRepository_CountByBoard(t *testing.T) {
	repo := mocks.NewPostRepository()

	// Test empty board
	count, err := repo.CountByBoard(1)
	if err != nil {
		t.Errorf("CountByBoard should not return error: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 posts, got %d", count)
	}

	// Create posts for different boards
	post1 := domain.NewPost(1, 42, "user1", "Post 1", "Content 1")
	post2 := domain.NewPost(1, 43, "user2", "Post 2", "Content 2")
	post3 := domain.NewPost(2, 44, "user3", "Post 3", "Content 3")
	reply := domain.NewReply(1, 45, "user4", "Reply content", 0)

	repo.Create(post1)
	repo.Create(post2)
	repo.Create(post3)
	*reply.ReplyTo = post1.ID
	repo.Create(reply)

	// Test count for board 1 (should include replies)
	count, err = repo.CountByBoard(1)
	if err != nil {
		t.Errorf("CountByBoard should not return error: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 posts for board 1 (including reply), got %d", count)
	}

	// Test count for board 2
	count, err = repo.CountByBoard(2)
	if err != nil {
		t.Errorf("CountByBoard should not return error: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 post for board 2, got %d", count)
	}
}
