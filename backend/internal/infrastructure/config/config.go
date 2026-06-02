package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration read from the environment at startup.
// All fields are validated before the server starts — fail fast, never silently.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port int
	Env  string // "development" | "production"
}

type DatabaseConfig struct {
	DSN string // postgres://user:pass@host:port/db?sslmode=...
}

type JWTConfig struct {
	Secret         string
	AccessTTLMin   int
	RefreshTTLDays int
}

// Load reads and validates all required environment variables.
// Returns an error listing every missing or invalid variable — never partial config.
func Load() (*Config, error) {
	var errs []string

	// Server
	portStr := requireEnv("SERVER_PORT", &errs)
	env := requireEnv("SERVER_ENV", &errs)

	// Database
	dsn := requireEnv("GRAMA_DB_DSN", &errs)

	// JWT
	jwtSecret := requireEnv("JWT_SECRET", &errs)
	accessTTLStr := requireEnv("JWT_ACCESS_TTL_MINUTES", &errs)
	refreshTTLStr := requireEnv("JWT_REFRESH_TTL_DAYS", &errs)

	if len(errs) > 0 {
		return nil, fmt.Errorf("missing required environment variables:\n  %s", strings.Join(errs, "\n  "))
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("SERVER_PORT must be a valid port number (1-65535), got %q", portStr)
	}

	if env != "development" && env != "production" {
		return nil, fmt.Errorf("SERVER_ENV must be 'development' or 'production', got %q", env)
	}

	if len(jwtSecret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 characters")
	}

	accessTTL, err := strconv.Atoi(accessTTLStr)
	if err != nil || accessTTL < 1 {
		return nil, fmt.Errorf("JWT_ACCESS_TTL_MINUTES must be a positive integer, got %q", accessTTLStr)
	}

	refreshTTL, err := strconv.Atoi(refreshTTLStr)
	if err != nil || refreshTTL < 1 {
		return nil, fmt.Errorf("JWT_REFRESH_TTL_DAYS must be a positive integer, got %q", refreshTTLStr)
	}

	return &Config{
		Server:   ServerConfig{Port: port, Env: env},
		Database: DatabaseConfig{DSN: dsn},
		JWT:      JWTConfig{Secret: jwtSecret, AccessTTLMin: accessTTL, RefreshTTLDays: refreshTTL},
	}, nil
}

func requireEnv(key string, errs *[]string) string {
	v := os.Getenv(key)
	if v == "" {
		*errs = append(*errs, key)
	}
	return v
}
