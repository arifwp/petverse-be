package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/arifwahyu/petverse-be/internal/modules/auth"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func AuthMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }
            // Check Bearer token format
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
                return
            }
            tokenString := parts[1]
            // Validate the token
            claims, err := authService.ValidateToken(tokenString)
            if err != nil {
                http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
                return
            }
            // Extract user ID from claims
            userIDStr, ok := claims["sub"].(string)
            if !ok {
                http.Error(w, "Invalid token claims", http.StatusUnauthorized)
                return
            }
            userID, err := uuid.Parse(userIDStr)
            if err != nil {
                http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
                return
            }
            // Add user ID to request context
            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            // Call the next handler with the enhanced context
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	return userID, ok
}