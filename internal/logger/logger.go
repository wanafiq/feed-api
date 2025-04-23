package logger

import (
	"os"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	env := os.Getenv("ENV")
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
