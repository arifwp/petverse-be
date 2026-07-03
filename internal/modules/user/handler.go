// internal/modules/user/handler.go
package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arifwahyu/petverse-be/internal/authctx"
	"github.com/arifwahyu/petverse-be/internal/shared/apiresponse"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepo *UserRepository
}

func NewUserHandler(userRepo *UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := authctx.GetUserID(r)
	if !ok {
		apiresponse.Error(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	user, err := h.userRepo.GetUserById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiresponse.Error(w, http.StatusNotFound, "User not found")
			return
		}

		apiresponse.Error(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type UpdateProfileRequest struct {
	Name      *string `json:"name"`
	Username  *string `json:"username"`
	AvatarURL *string `json:"avatar_url"`
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId, ok := authctx.GetUserID(r)
	if !ok {
		apiresponse.Unauthorized(w, "Unauthorized")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiresponse.BadRequest(w, "Invalid request payload")
		return
	}

	user, err := h.userRepo.UpdateProfile(userId, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			apiresponse.NotFound(w, "User not found")
			return
		}

		apiresponse.InternalServerError(w)
		return
	}

	apiresponse.Success(
		w,
		http.StatusOK,
		"Profile updated successfully",
		user,
	)
}
