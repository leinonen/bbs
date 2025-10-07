package test

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/leinonen/bbs/domain"
	"github.com/leinonen/bbs/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Create temporary database file
	tmpFile, err := ioutil.TempFile("", "test_bbs_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	// Open database connection
	db, err := sql.Open("sqlite3", tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Load test schema
	schemaPath := filepath.Join("testdata", "test_schema.sql")
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}

	// Execute schema
	if _, err := db.Exec(string(schemaBytes)); err != nil {
		t.Fatalf("Failed to execute schema: %v", err)
	}

	// Clean up function
	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpFile.Name())
	})

	return db
}

func TestSQLiteUserRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewUserRepository(db)

	// Test Create
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"

	err := repo.Create(user)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("Create should set user ID")
	}

	// Test GetByID
	retrievedUser, err := repo.GetByID(user.ID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if retrievedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedUser.Username)
	}

	// Test GetByUsername
	userByName, err := repo.GetByUsername("testuser")
	if err != nil {
		t.Errorf("GetByUsername failed: %v", err)
	}

	if userByName.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, userByName.ID)
	}

	// Test Authenticate
	authUser, err := repo.Authenticate("testuser", "password123")
	if err != nil {
		t.Errorf("Authenticate failed: %v", err)
	}

	if authUser.ID != user.ID {
		t.Errorf("Expected authenticated user ID %d, got %d", user.ID, authUser.ID)
	}

	// Test wrong password
	_, err = repo.Authenticate("testuser", "wrongpassword")
	if err == nil {
		t.Error("Authenticate should fail with wrong password")
	}

	// Test Update
	user.Email = "newemail@example.com"
	user.IsAdmin = true

	err = repo.Update(user)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	updatedUser, _ := repo.GetByID(user.ID)
	if updatedUser.Email != "newemail@example.com" {
		t.Errorf("Expected updated email, got %s", updatedUser.Email)
	}

	if !updatedUser.IsAdmin {
		t.Error("Expected user to be admin after update")
	}

	// Test UpdateLastLogin
	err = repo.UpdateLastLogin(user.ID)
	if err != nil {
		t.Errorf("UpdateLastLogin failed: %v", err)
	}

	// Test Delete
	err = repo.Delete(user.ID)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(user.ID)
	if err == nil {
		t.Error("GetByID should fail for deleted user")
	}
}

func TestSQLiteBoardRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	repo := sqlite.NewBoardRepository(db)

	// Test Create
	board := domain.NewBoard("General", "General discussion")

	err := repo.Create(board)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if board.ID == 0 {
		t.Error("Create should set board ID")
	}

	// Test GetByID
	retrievedBoard, err := repo.GetByID(board.ID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if retrievedBoard.Name != board.Name {
		t.Errorf("Expected name %s, got %s", board.Name, retrievedBoard.Name)
	}

	// Test GetAll
	board2 := domain.NewBoard("Tech", "Technology discussions")
	repo.Create(board2)

	boards, err := repo.GetAll()
	if err != nil {
		t.Errorf("GetAll failed: %v", err)
	}

	if len(boards) != 2 {
		t.Errorf("Expected 2 boards, got %d", len(boards))
	}

	// Test Update
	board.Name = "Updated General"
	board.Description = "Updated description"

	err = repo.Update(board)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	updatedBoard, _ := repo.GetByID(board.ID)
	if updatedBoard.Name != "Updated General" {
		t.Errorf("Expected updated name, got %s", updatedBoard.Name)
	}

	// Test Delete
	err = repo.Delete(board.ID)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(board.ID)
	if err == nil {
		t.Error("GetByID should fail for deleted board")
	}

	// Verify it's not in GetAll results
	boards, _ = repo.GetAll()
	if len(boards) != 1 {
		t.Errorf("Expected 1 board after deletion, got %d", len(boards))
	}
}

func TestSQLitePostRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	userRepo := sqlite.NewUserRepository(db)
	boardRepo := sqlite.NewBoardRepository(db)
	postRepo := sqlite.NewPostRepository(db)

	// Set up test data
	user := domain.NewUser("testuser", "test@example.com")
	user.Password = "password123"
	userRepo.Create(user)

	board := domain.NewBoard("General", "General discussion")
	boardRepo.Create(board)

	// Test Create Post
	post := domain.NewPost(board.ID, user.ID, user.Username, "Test Post", "This is a test post")

	err := postRepo.Create(post)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if post.ID == 0 {
		t.Error("Create should set post ID")
	}

	// Test GetByID
	retrievedPost, err := postRepo.GetByID(post.ID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if retrievedPost.Title != post.Title {
		t.Errorf("Expected title %s, got %s", post.Title, retrievedPost.Title)
	}

	if retrievedPost.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrievedPost.Username)
	}

	// Test Create Reply
	reply := domain.NewReply(board.ID, user.ID, user.Username, "This is a reply", post.ID)

	err = postRepo.Create(reply)
	if err != nil {
		t.Errorf("Create reply failed: %v", err)
	}

	// Test GetByBoard (should not include replies)
	posts, err := postRepo.GetByBoard(board.ID, 10, 0)
	if err != nil {
		t.Errorf("GetByBoard failed: %v", err)
	}

	if len(posts) != 1 {
		t.Errorf("Expected 1 post (excluding reply), got %d", len(posts))
	}

	if posts[0].ID != post.ID {
		t.Error("GetByBoard should return original post, not reply")
	}

	// Test GetReplies
	replies, err := postRepo.GetReplies(post.ID)
	if err != nil {
		t.Errorf("GetReplies failed: %v", err)
	}

	if len(replies) != 1 {
		t.Errorf("Expected 1 reply, got %d", len(replies))
	}

	if replies[0].Content != "This is a reply" {
		t.Errorf("Expected reply content, got %s", replies[0].Content)
	}

	// Test GetRecent
	recentPosts, err := postRepo.GetRecent(10)
	if err != nil {
		t.Errorf("GetRecent failed: %v", err)
	}

	if len(recentPosts) != 2 {
		t.Errorf("Expected 2 recent posts (post + reply), got %d", len(recentPosts))
	}

	// Test CountByBoard
	count, err := postRepo.CountByBoard(board.ID)
	if err != nil {
		t.Errorf("CountByBoard failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 posts in board (post + reply), got %d", count)
	}

	// Test Update
	post.Title = "Updated Title"
	post.Content = "Updated content"

	err = postRepo.Update(post)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	updatedPost, _ := postRepo.GetByID(post.ID)
	if updatedPost.Title != "Updated Title" {
		t.Errorf("Expected updated title, got %s", updatedPost.Title)
	}

	// Test Delete
	err = postRepo.Delete(reply.ID)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	replies, _ = postRepo.GetReplies(post.ID)
	if len(replies) != 0 {
		t.Errorf("Expected 0 replies after deletion, got %d", len(replies))
	}
}
