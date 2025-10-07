package sqlite

import (
	"database/sql"
	"errors"

	"github.com/leinonen/bbs/domain"
)

type BoardRepository struct {
	db *sql.DB
}

func NewBoardRepository(db *sql.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(board *domain.Board) error {
	query := `
		INSERT INTO boards (name, description, created_at)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, board.Name, board.Description, board.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	board.ID = int(id)
	return nil
}

func (r *BoardRepository) GetByID(id int) (*domain.Board, error) {
	board := &domain.Board{}
	query := `
		SELECT b.id, b.name, b.description, b.created_at, COUNT(p.id) as post_count
		FROM boards b
		LEFT JOIN posts p ON b.id = p.board_id
		WHERE b.id = ?
		GROUP BY b.id
	`

	err := r.db.QueryRow(query, id).Scan(
		&board.ID,
		&board.Name,
		&board.Description,
		&board.CreatedAt,
		&board.PostCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	return board, nil
}

func (r *BoardRepository) GetAll() ([]*domain.Board, error) {
	query := `
		SELECT b.id, b.name, b.description, b.created_at, COUNT(p.id) as post_count
		FROM boards b
		LEFT JOIN posts p ON b.id = p.board_id
		GROUP BY b.id
		ORDER BY b.name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []*domain.Board
	for rows.Next() {
		board := &domain.Board{}
		err := rows.Scan(
			&board.ID,
			&board.Name,
			&board.Description,
			&board.CreatedAt,
			&board.PostCount,
		)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, rows.Err()
}

func (r *BoardRepository) Update(board *domain.Board) error {
	query := `
		UPDATE boards
		SET name = ?, description = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, board.Name, board.Description, board.ID)
	return err
}

func (r *BoardRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM boards WHERE id = ?", id)
	return err
}
