// internal/modules/user/routes.go
package user

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	h *UserHandler,
	authMiddleware func(http.Handler) http.Handler,
) {
	mux.Handle(
		"GET /api/profile",
		authMiddleware(http.HandlerFunc(h.Profile)),
	)

	mux.Handle(
		"PUT /api/profile",
		authMiddleware(http.HandlerFunc(h.UpdateProfile)),
	)
}
