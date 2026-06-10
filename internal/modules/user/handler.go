package user

import (
	"encoding/json"
	"net/http"

	"github.com/arifwahyu/petverse-be/internal/authctx"
)

type UserHandler struct {
	userRepo *UserRepository
}

func NewUserHandler(userRepo *UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId, ok := authctx.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetUserById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
