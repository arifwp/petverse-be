package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/arifwahyu/petverse-be/internal/modules/auth"
	"github.com/arifwahyu/petverse-be/internal/modules/user"
	"github.com/arifwahyu/petverse-be/internal/shared/apiresponse"
)

func registerSystemRoutes(mux *http.ServeMux, readiness ReadinessChecker) {
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /readyz", readyHandler(readiness))
}

func registerRoutes(
	mux *http.ServeMux,
	readiness ReadinessChecker,
	authHandler *auth.AuthHandler,
	userHandler *user.UserHandler,
	authMiddleware Middleware,
) {
	registerSystemRoutes(mux, readiness)
	registerAPIRoutes(mux, authHandler, userHandler, authMiddleware)
}

func registerAPIRoutes(
	mux *http.ServeMux,
	authHandler *auth.AuthHandler,
	userHandler *user.UserHandler,
	authMiddleware Middleware,
) {
	auth.RegisterRoutes(mux, authHandler)
	user.RegisterRoutes(mux, userHandler, authMiddleware)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	apiresponse.OK(
		w,
		map[string]string{
			"status": "ok",
		},
	)
}

func readyHandler(readiness ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC().Format(time.RFC3339)

		if readiness == nil {
			apiresponse.Error(
				w,
				http.StatusServiceUnavailable,
				"Readiness checker is not configured",
			)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := readiness.Ping(ctx); err != nil {
			apiresponse.Error(
				w,
				http.StatusServiceUnavailable,
				"Database is unavailable",
			)
			return
		}

		apiresponse.Success(
			w,
			http.StatusOK,
			"Service is ready",
			map[string]any{
				"status":   "ready",
				"database": "ok",
				"time":     now,
			},
		)
	}
}
