package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App         AppConfig
	HTTP        HTTPConfig
	Security    SecurityConfig
	Database    DatabaseConfig
	Cache       CacheConfig
	ObjectStore ObjectStoreConfig
	Realtime    RealtimeConfig
}

type AppConfig struct {
	Name           string
	Env            string
	StartupTimeout time.Duration
}

type HTTPConfig struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func (c HTTPConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type SecurityConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AllowedOrigins     string
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	PingTimeout     time.Duration
}

type CacheConfig struct {
	URL string
}

type ObjectStoreConfig struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type RealtimeConfig struct {
	MaxConnections int
}

func Load() (Config, error) {
	cfg := Config{
		App: AppConfig{
			Name:           getEnv("APP_NAME", "petverse-be"),
			Env:            strings.ToLower(getEnv("APP_ENV", "local")),
			StartupTimeout: getDuration("APP_STARTUP_TIMEOUT", 250*time.Millisecond),
		},
		HTTP: HTTPConfig{
			Host:              getEnv("HTTP_HOST", "0.0.0.0"),
			Port:              getInt("HTTP_PORT", 8080),
			ReadHeaderTimeout: getDuration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second),
			ReadTimeout:       getDuration("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:      getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:       getDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout:   getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Security: SecurityConfig{
			AccessTokenSecret:  getEnv("JWT_ACCESS_SECRET", ""),
			RefreshTokenSecret: getEnv("JWT_REFRESH_SECRET", ""),
			AllowedOrigins:     getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", ""),
			MaxOpenConns:    getInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getInt("DATABASE_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
			PingTimeout:     getDuration("DATABASE_PING_TIMEOUT", 3*time.Second),
		},
		Cache: CacheConfig{
			URL: getEnv("REDIS_URL", ""),
		},
		ObjectStore: ObjectStoreConfig{
			Endpoint:  getEnv("OBJECT_STORE_ENDPOINT", ""),
			Bucket:    getEnv("OBJECT_STORE_BUCKET", "petverse-media"),
			AccessKey: getEnv("OBJECT_STORE_ACCESS_KEY", ""),
			SecretKey: getEnv("OBJECT_STORE_SECRET_KEY", ""),
			UseSSL:    getBool("OBJECT_STORE_USE_SSL", false),
		},
		Realtime: RealtimeConfig{
			MaxConnections: getInt("REALTIME_MAX_CONNECTIONS", 1000),
		},
	}

	return cfg, validate(cfg)
}

func validate(cfg Config) error {
	if !isAllowedEnv(cfg.App.Env) {
		return fmt.Errorf("%w: APP_ENV must be one of local, development, staging, production", ErrInvalidConfig)
	}

	if cfg.HTTP.Port <= 0 || cfg.HTTP.Port > 65535 {
		return fmt.Errorf("%w: HTTP_PORT must be between 1 and 65535: %d", ErrInvalidConfig, cfg.HTTP.Port)
	}

	if cfg.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("%w: DATABASE_MAX_OPEN_CONNS must be greater than 0", ErrInvalidConfig)
	}

	if cfg.Database.MaxIdleConns < 0 {
		return fmt.Errorf("%w: DATABASE_MAX_IDLE_CONNS cannot be negative", ErrInvalidConfig)
	}

	if cfg.App.Env != "local" {
		var missing []string
		if cfg.Security.AccessTokenSecret == "" {
			missing = append(missing, "JWT_ACCESS_SECRET")
		}
		if cfg.Security.RefreshTokenSecret == "" {
			missing = append(missing, "JWT_REFRESH_SECRET")
		}
		if cfg.Database.URL == "" {
			missing = append(missing, "DATABASE_URL")
		}
		if len(missing) > 0 {
			return fmt.Errorf("%w: missing %s configuration: %v", ErrInvalidConfig, cfg.App.Env, missing)
		}
	}

	return nil
}

func isAllowedEnv(env string) bool {
	switch strings.ToLower(env) {
	case "local", "development", "staging", "production":
		return true
	default:
		return false
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}

	return fallback
}

func getInt(key string, fallback int) int {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return value
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}

	return value
}

func getBool(key string, fallback bool) bool {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}

	return value
}

var ErrInvalidConfig = errors.New("invalid configuration")
