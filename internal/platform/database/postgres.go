package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/arifwahyu/petverse-be/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrMissingURL = errors.New("database url is required")

type Store struct {
	db *gorm.DB
}

func Open(ctx context.Context, cfg config.DatabaseConfig) (*Store, error) {
	if cfg.URL == "" {
		return nil, ErrMissingURL
	}

	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get postgres connection pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) DB() *gorm.DB {
	return s.db
}

func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.db == nil {
		return ErrMissingURL
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
