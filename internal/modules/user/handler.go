package user

import (
	"encoding/json"
	"net/http"

	"github.com/arifwahyu/petverse-be/internal/modules/auth"
	"github.com/arifwahyu/petverse-be/internal/modules/middleware"
)

type UserHandler struct {
	userRepo *auth.UserRepository
}

func NewUserHandler(userRepo *auth.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

type UserResponse struct {
	ID string `json:"id"`	
	Email string `json:"email"`
	Username string `json:"username"`
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetUserById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := UserResponse {
		ID: user.ID.String(),
		Email: user.Email,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}