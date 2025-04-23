package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Env         string
	DatabaseURL string
	Port        string
	JWTSecret   string
	Smtp        *smtp
	Url         *url
}

type smtp struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type url struct {
	Web string
}

func LoadConfig() (*Config, error) {
	if err := validateRequiredConfig(); err != nil {
		return nil, err
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	smptpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, err
	}

	smtp := &smtp{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     smptpPort,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}

	url := &url{
		Web: os.Getenv("WEB_URL"),
	}

	return &Config{
		Env:         env,
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Smtp:        smtp,
		Url:         url,
	}, nil
}

func validateRequiredConfig() error {
	if isEmpty := os.Getenv("DATABASE_URL") == ""; isEmpty {
		return errors.New("DATABASE_URL is not set")
	}

	if isEmpty := os.Getenv("JWT_SECRET") == ""; isEmpty {
		return errors.New("JWT_SECRET is not set")
	}

	if isEmpty := os.Getenv("SMTP_HOST") == ""; isEmpty {
		return errors.New("SMTP_HOST is not set")
	}

	if isEmpty := os.Getenv("SMTP_PORT") == ""; isEmpty {
		return errors.New("SMTP_PORT is not set")
	}

	if isEmpty := os.Getenv("SMTP_USERNAME") == ""; isEmpty {
		return errors.New("SMTP_USERNAME is not set")
	}

	if isEmpty := os.Getenv("SMTP_PASSWORD") == ""; isEmpty {
		return errors.New("SMTP_PASSWORD is not set")
	}

	if isEmpty := os.Getenv("SMTP_FROM") == ""; isEmpty {
		return errors.New("SMTP_FROM is not set")
	}

	return nil
}
