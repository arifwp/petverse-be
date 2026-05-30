package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func registerSystemRoutes(mux *http.ServeMux, readiness ReadinessChecker) {
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /readyz", readyHandler(readiness))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func readyHandler(readiness ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if readiness == nil {
			writeJSON(w, http.StatusOK, map[string]any{
				"status":   "ready",
				"database": "not_configured",
				"time":     time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := readiness.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"status":   "not_ready",
				"database": "unavailable",
				"time":     time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"status":   "ready",
			"database": "ok",
			"time":     time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
