package auth

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *AuthHandler) {
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("POST /api/auth/refresh", h.RefreshToken)
}
