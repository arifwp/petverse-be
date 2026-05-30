package httpserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/arifwahyu/petverse-be/internal/config"
)

type ReadinessChecker interface {
	Ping(context.Context) error
}

type Server struct {
	httpServer *http.Server
}

func New(cfg config.HTTPConfig, log *slog.Logger, readiness ReadinessChecker) *Server {
	mux := http.NewServeMux()
	registerSystemRoutes(mux, readiness)

	handler := recoverer(log)(requestLogger(log)(securityHeaders(mux)))

	return &Server{
		httpServer: &http.Server{
			Addr:              cfg.Address(),
			Handler:           handler,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			BaseContext: func(_ net.Listener) context.Context {
				return context.Background()
			},
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
