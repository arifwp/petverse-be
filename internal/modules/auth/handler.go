// internal/modules/auth/handler.go
package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/arifwahyu/petverse-be/internal/shared/apiresponse"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
		)
		return
	}

	if req.Email == "" {
		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Email is required",
		)
		return
	}
	if req.Name == "" {
		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Name is required",
		)
		return
	}
	if req.Username == "" {
		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Username is required",
		)
		return
	}
	if req.Password == "" {

		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Password is required",
		)
		return
	}

	user, err := h.authService.Register(req.Email, req.Name, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrEmailInUse) {
			apiresponse.Error(
				w,
				http.StatusConflict,
				"Email already in use",
			)
			return
		}

		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Error creating user",
		)
		return
	}

	response := RegisterResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
		)
		return
	}

	// refresh token 7 days
	refreshTokenTTL := 7 * 24 * time.Hour

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password, refreshTokenTTL)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {

			apiresponse.Error(
				w,
				http.StatusUnauthorized,
				"Invalid credentials",
			)
		} else {
			apiresponse.Error(
				w,
				http.StatusInternalServerError,
				"Internal server error",
			)
		}
		return
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Token string `json:"access_token"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiresponse.Error(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
		)
		return
	}

	token, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrExpiredToken) {

			apiresponse.Error(
				w,
				http.StatusUnauthorized,
				"Invalid or expired refresh token",
			)
		} else {
			apiresponse.Error(
				w,
				http.StatusInternalServerError,
				"Internal server error",
			)
		}
		return
	}

	response := RefreshResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
