package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/arifwahyu/petverse-be/internal/config"
	"github.com/arifwahyu/petverse-be/internal/modules/auth"
	"github.com/arifwahyu/petverse-be/internal/modules/user"
	"github.com/arifwahyu/petverse-be/internal/platform/database"
	"github.com/arifwahyu/petverse-be/internal/platform/httpserver"
	"github.com/arifwahyu/petverse-be/internal/platform/logger"
	"github.com/arifwahyu/petverse-be/internal/platform/shutdown"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func Run(parent context.Context) error {
	loadEnv()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg.App.Env)
	ctx, stop := shutdown.Context(parent)
	defer stop()

	var store *database.Store
	if cfg.Database.URL != "" {
		store, err = database.Open(ctx, cfg.Database)
		if err != nil {
			return err
		}
		defer func() {
			if err := store.Close(); err != nil {
				log.Error("close database", "error", err)
			}
		}()
		log.Info("postgres connected")
	}

	var readiness httpserver.ReadinessChecker
	if store != nil {
		readiness = store
	}

	userRepo := user.NewUserRepository(store.DB())
	refreshTokenRepo := auth.NewRefreshTokenRepository(store.DB())

	authService := auth.NewAuthService(
		userRepo,
		refreshTokenRepo,
		cfg.Security.AccessTokenSecret,
		15*time.Minute,
	)

	authHandler := auth.NewAuthHandler(authService)
	userHandler := user.NewUserHandler(userRepo)

	server := httpserver.New(
		cfg.HTTP,
		log,
		readiness,
		authService,
		authHandler,
		userHandler,
	)
	errCh := make(chan error, 1)

	go func() {
		log.Info(
			"starting api",
			"service", cfg.App.Name,
			"environment", cfg.App.Env,
			"address", cfg.HTTP.Address(),
		)
		errCh <- server.Start()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		log.Info("shutting down api")
		return server.Stop(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		slog.Error("http server failed", "error", err)
		return err
	case <-time.After(cfg.App.StartupTimeout):
		log.Info("api startup window completed")
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
			defer cancel()
			return server.Stop(shutdownCtx)
		case err := <-errCh:
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		}
	}
}
