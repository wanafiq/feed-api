package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env         string // "development" or "production"
	DatabaseURL string
	Port        string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // In production, use a secure key.
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		DatabaseURL: databaseURL,
		Port:        port,
		JWTSecret:   jwtSecret,
		Env:         env,
	}, nil
}
