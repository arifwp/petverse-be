package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateRefreshToken(userID uuid.UUID, ttl time.Duration) (*RefreshToken, error) {
	tokenID := uuid.New()
	expiredAt := time.Now().Add(ttl)

	token := &RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     tokenID.String(),
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
		Revoked:   false,
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
