package httpserver

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/arifwahyu/petverse-be/internal/authctx"
	"github.com/arifwahyu/petverse-be/internal/modules/auth"
	"github.com/arifwahyu/petverse-be/internal/shared/apiresponse"
	"github.com/google/uuid"
)

type Middleware func(http.Handler) http.Handler

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func requestLogger(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			next.ServeHTTP(w, r)
			log.Info(
				"http request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(startedAt).String(),
			)
		})
	}
}

func recoverer(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					log.Error(
						"panic recovered",
						"method", r.Method,
						"path", r.URL.Path,
						"panic", recovered,
						"stack", string(debug.Stack()),
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func authMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// http.Error(w, "Authorization header required", http.StatusUnauthorized)
				apiresponse.Error(
					w,
					http.StatusUnauthorized,
					"Authorization header required",
				)
				return
			}
			// Check Bearer token format
			parts := strings.Fields(authHeader)
			if len(parts) != 2 || parts[0] != "Bearer" {
				apiresponse.Error(
					w,
					http.StatusUnauthorized,
					"Invalid authorization format",
				)
				return
			}

			tokenString := parts[1]
			// Validate the token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				apiresponse.Error(
					w,
					http.StatusUnauthorized,
					"Invalid or expired token",
				)
				return
			}
			// Extract user ID from claims
			userIDStr, ok := claims["sub"].(string)
			if !ok {
				apiresponse.Error(
					w,
					http.StatusUnauthorized,
					"Invalid token claims",
				)
				return
			}
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				apiresponse.Error(
					w,
					http.StatusUnauthorized,
					"Invalid user ID in token",
				)
				return
			}
			// Add user ID to request context
			ctx := authctx.SetUserID(r.Context(), userID)
			// Call the next handler with the enhanced context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
