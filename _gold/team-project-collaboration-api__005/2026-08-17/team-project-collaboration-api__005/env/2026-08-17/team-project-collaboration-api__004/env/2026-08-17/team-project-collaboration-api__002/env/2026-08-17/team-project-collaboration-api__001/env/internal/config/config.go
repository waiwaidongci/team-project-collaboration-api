package config

import (
	"os"
	"time"
)

const (
	defaultDatabaseURL = "postgres://postgres:postgres@127.0.0.1:55432/team_project_task?sslmode=disable"
	defaultHTTPAddr    = ":18108"
	defaultJWTSecret   = "dev-only-change-me"
	defaultTokenTTL    = 72 * time.Hour
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
	JWTSecret   string
	TokenTTL    time.Duration
	Environment string
}

func Load() Config {
	return Config{
		DatabaseURL: getenv("DATABASE_URL", defaultDatabaseURL),
		HTTPAddr:    getenv("HTTP_ADDR", defaultHTTPAddr),
		JWTSecret:   getenv("JWT_SECRET", defaultJWTSecret),
		TokenTTL:    getDuration("TOKEN_TTL", defaultTokenTTL),
		Environment: getenv("APP_ENV", "development"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}
