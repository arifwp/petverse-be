package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Name          *string    `gorm:"column:name" json:"name,omitempty"`
	Username      string     `gorm:"column:username;not null;uniqueIndex" json:"username"`
	Email         string     `gorm:"column:email;not null;uniqueIndex" json:"email"`
	AvatarURL     *string    `gorm:"column:avatar_url" json:"avatar_url,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null" json:"updated_at"`
	LastLoginAt   *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	EmailVerified bool       `gorm:"column:email_verified;not null;default:false" json:"email_verified"`
	PasswordHash  string     `gorm:"column:password;not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}
