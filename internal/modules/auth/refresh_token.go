package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id;not null;index" json:"user_id"`
	Token     string    `gorm:"column:token;not null;uniqueIndex" json:"token"`
	ExpiredAt time.Time `gorm:"column:expired_at;not null" json:"expired_at"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	Revoked   bool      `gorm:"column:revoked;not null;default:false" json:"revoked"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
