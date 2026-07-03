package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, name, username, passwordHash string) (*User, error) {
	now := time.Now()
	user := &User{
		ID:           uuid.New(),
		Name:         stringPtr(name),
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateProfile(
	userID uuid.UUID,
	req UpdateProfileRequest,
) (*User, error) {
	updates := map[string]interface{}{}

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Username != nil {
		updates["username"] = *req.Username
	}

	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	if len(updates) == 0 {
		return r.GetUserById(userID)
	}

	var user User

	err := r.db.
		Model(&User{}).
		Where("id = ?", userID).
		Updates(updates).
		First(&user, "id = ?", userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
