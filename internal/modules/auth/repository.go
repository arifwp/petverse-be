package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

type RefreshTokenRepository struct {
    db *sql.DB
}


func (r *UserRepository) CreateUser(email, name, username, passwordHash string) (*User, error) {
	user := &User {
		ID: uuid.New(),
		Email: email,
		Username: username,
		PasswordHash: passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (r *UserRepository) GetUserByEmail(email string) (*User,error) {
	query := `SELECT email FROM users WHERE email = $1`

	var user User
	
	err := r.db.QueryRow(query,email).Scan(&user.Email)

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

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}


func (r *RefreshTokenRepository) CreateRefreshToken(userID uuid.UUID, ttl time.Duration) (*RefreshToken, error) {
	tokenID := uuid.New()
	expiredAt := time.Now().Add(ttl)

	token := &RefreshToken {
		ID: tokenID,
		UserID:    userID,
		Token: tokenID.String(),
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
		Revoked: false,
	}

	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expired_at, created_at, revoked)
        VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiredAt, token.CreatedAt, token.Revoked)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *RefreshTokenRepository) GetRefreshToken(tokenString string) (*RefreshToken, error) {
	query := `SELECT id, user_id, token, expired_at, revoked FROM refresh_tokens WHERE token = $1`

	var token RefreshToken
	err := r.db.QueryRow(query, tokenString).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiredAt,
		&token.Revoked,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(tokenString string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1`

	_, err := r.db.Exec(query, tokenString)
	return err
}