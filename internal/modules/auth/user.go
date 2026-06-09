package auth

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	Name          *string    `db:"name" json:"name,omitempty"`
	Username      string     `db:"username" json:"username"`
	Email         string     `db:"email" json:"email"`
	AvatarURL     *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	LastLoginAt   *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	EmailVerified bool       `db:"email_verified" json:"email_verified"`
	PasswordHash  string     `db:"password" json:"-"`
}

