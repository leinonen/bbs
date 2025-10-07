package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/leinonen/bbs/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (username, password, email, created_at, last_login, is_admin)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		user.Username,
		string(hashedPassword),
		user.Email,
		user.CreatedAt,
		user.LastLogin,
		user.IsAdmin)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, username, email, created_at, last_login, is_admin
		FROM users WHERE id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.LastLogin,
		&user.IsAdmin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, username, email, created_at, last_login, is_admin
		FROM users WHERE username = ?
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.LastLogin,
		&user.IsAdmin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, is_admin = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, user.Username, user.Email, user.IsAdmin, user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (r *UserRepository) Authenticate(username, password string) (*domain.User, error) {
	user := &domain.User{}
	var hashedPassword string

	query := `
		SELECT id, username, password, email, created_at, last_login, is_admin
		FROM users WHERE username = ?
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&hashedPassword,
		&user.Email,
		&user.CreatedAt,
		&user.LastLogin,
		&user.IsAdmin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := "UPDATE users SET last_login = ? WHERE id = ?"
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}
