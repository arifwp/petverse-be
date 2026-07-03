package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
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

	if err := r.db.Create(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func (r *RefreshTokenRepository) GetRefreshToken(tokenString string) (*RefreshToken, error) {
	var token RefreshToken
	if err := r.db.Where("token = ?", tokenString).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(tokenString string) error {
	return r.db.Model(&RefreshToken{}).
		Where("token = ?", tokenString).
		Update("revoked", true).
		Error
}
