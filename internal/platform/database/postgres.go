package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/arifwahyu/petverse-be/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrMissingURL = errors.New("database url is required")

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context, cfg config.DatabaseConfig) (*Store, error) {
	if cfg.URL == "" {
		return nil, ErrMissingURL
	}

	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.db == nil {
		return ErrMissingURL
	}

	return s.db.PingContext(ctx)
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	return s.db.Close()
}
