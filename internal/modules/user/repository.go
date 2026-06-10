// internal/modules/user/repository.go
package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(email, name, username, passwordHash string) (*User, error) {
	user := &User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `
		INSERT INTO users (id, email, username, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user User

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id uuid.UUID) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var user User

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.EmailVerified,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
