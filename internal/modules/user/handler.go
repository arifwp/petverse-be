package user

import (
	"encoding/json"
	"net/http"

	"github.com/arifwahyu/petverse-be/internal/authctx"
	"github.com/arifwahyu/petverse-be/internal/shared/apiresponse"
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
		apiresponse.Error(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	user, err := h.userRepo.GetUserById(userId)
	if err != nil {
		apiresponse.Error(
			w,
			http.StatusNotFound,
			"User not found",
		)
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
