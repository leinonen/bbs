package sqlite

import (
	"database/sql"
	"errors"

	"github.com/leinonen/bbs/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *domain.Post) error {
	query := `
		INSERT INTO posts (board_id, user_id, title, content, created_at, updated_at, reply_to)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var replyTo sql.NullInt64
	if post.ReplyTo != nil {
		replyTo = sql.NullInt64{Int64: int64(*post.ReplyTo), Valid: true}
	}

	result, err := r.db.Exec(query,
		post.BoardID,
		post.UserID,
		post.Title,
		post.Content,
		post.CreatedAt,
		post.UpdatedAt,
		replyTo)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = int(id)
	return nil
}

func (r *PostRepository) GetByID(id int) (*domain.Post, error) {
	post := &domain.Post{}
	var replyTo sql.NullInt64

	query := `
		SELECT p.id, p.board_id, p.user_id, u.username, p.title, p.content,
		       p.created_at, p.updated_at, p.reply_to,
		       (SELECT COUNT(*) FROM posts WHERE reply_to = p.id) as reply_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.BoardID,
		&post.UserID,
		&post.Username,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&replyTo,
		&post.Replies,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	if replyTo.Valid {
		replyToInt := int(replyTo.Int64)
		post.ReplyTo = &replyToInt
	}

	return post, nil
}

func (r *PostRepository) GetByBoard(boardID int, limit, offset int) ([]*domain.Post, error) {
	query := `
		SELECT p.id, p.board_id, p.user_id, u.username, p.title, p.content,
		       p.created_at, p.updated_at, p.reply_to,
		       (SELECT COUNT(*) FROM posts WHERE reply_to = p.id) as reply_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.board_id = ? AND p.reply_to IS NULL
		ORDER BY p.updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, boardID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		var replyTo sql.NullInt64

		err := rows.Scan(
			&post.ID,
			&post.BoardID,
			&post.UserID,
			&post.Username,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&replyTo,
			&post.Replies,
		)
		if err != nil {
			return nil, err
		}

		if replyTo.Valid {
			replyToInt := int(replyTo.Int64)
			post.ReplyTo = &replyToInt
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (r *PostRepository) GetReplies(postID int) ([]*domain.Post, error) {
	query := `
		SELECT p.id, p.board_id, p.user_id, u.username, p.title, p.content,
		       p.created_at, p.updated_at, p.reply_to, 0 as reply_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.reply_to = ?
		ORDER BY p.created_at ASC
	`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		var replyTo sql.NullInt64

		err := rows.Scan(
			&post.ID,
			&post.BoardID,
			&post.UserID,
			&post.Username,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&replyTo,
			&post.Replies,
		)
		if err != nil {
			return nil, err
		}

		if replyTo.Valid {
			replyToInt := int(replyTo.Int64)
			post.ReplyTo = &replyToInt
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (r *PostRepository) GetRecent(limit int) ([]*domain.Post, error) {
	query := `
		SELECT p.id, p.board_id, p.user_id, u.username, p.title, p.content,
		       p.created_at, p.updated_at, p.reply_to,
		       (SELECT COUNT(*) FROM posts WHERE reply_to = p.id) as reply_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		var replyTo sql.NullInt64

		err := rows.Scan(
			&post.ID,
			&post.BoardID,
			&post.UserID,
			&post.Username,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&replyTo,
			&post.Replies,
		)
		if err != nil {
			return nil, err
		}

		if replyTo.Valid {
			replyToInt := int(replyTo.Int64)
			post.ReplyTo = &replyToInt
		}

		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (r *PostRepository) Update(post *domain.Post) error {
	query := `
		UPDATE posts
		SET title = ?, content = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, post.Title, post.Content, post.UpdatedAt, post.ID)
	return err
}

func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}

func (r *PostRepository) CountByBoard(boardID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM posts WHERE board_id = ?"
	err := r.db.QueryRow(query, boardID).Scan(&count)
	return count, err
}