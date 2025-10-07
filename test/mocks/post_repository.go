package mocks

import (
	"errors"
	"sort"
	"sync"

	"github.com/leinonen/bbs/domain"
)

type PostRepository struct {
	mu     sync.RWMutex
	posts  map[int]*domain.Post
	nextID int
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		posts:  make(map[int]*domain.Post),
		nextID: 1,
	}
}

func (r *PostRepository) Create(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post.ID = r.nextID
	r.nextID++
	r.posts[post.ID] = post
	return nil
}

func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (r *PostRepository) GetByBoard(boardID int, limit, offset int) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var boardPosts []*domain.Post
	for _, post := range r.posts {
		if post.BoardID == boardID && post.ReplyTo == nil {
			boardPosts = append(boardPosts, post)
		}
	}

	// Sort by created time (newest first)
	sort.Slice(boardPosts, func(i, j int) bool {
		return boardPosts[i].CreatedAt.After(boardPosts[j].CreatedAt)
	})

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(boardPosts) {
		return []*domain.Post{}, nil
	}
	if end > len(boardPosts) {
		end = len(boardPosts)
	}

	return boardPosts[start:end], nil
}

func (r *PostRepository) GetReplies(postID int) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var replies []*domain.Post
	for _, post := range r.posts {
		if post.ReplyTo != nil && *post.ReplyTo == postID {
			replies = append(replies, post)
		}
	}

	// Sort by created time (oldest first for replies)
	sort.Slice(replies, func(i, j int) bool {
		return replies[i].CreatedAt.Before(replies[j].CreatedAt)
	})

	return replies, nil
}

func (r *PostRepository) GetRecent(limit int) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allPosts []*domain.Post
	for _, post := range r.posts {
		allPosts = append(allPosts, post)
	}

	// Sort by created time (newest first)
	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)
	})

	if limit > len(allPosts) {
		limit = len(allPosts)
	}

	return allPosts[:limit], nil
}

func (r *PostRepository) Update(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.ID]; !exists {
		return errors.New("post not found")
	}
	r.posts[post.ID] = post
	return nil
}

func (r *PostRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[id]; !exists {
		return errors.New("post not found")
	}
	delete(r.posts, id)
	return nil
}

func (r *PostRepository) CountByBoard(boardID int) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, post := range r.posts {
		if post.BoardID == boardID {
			count++
		}
	}
	return count, nil
}
