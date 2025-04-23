package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger() (*zap.SugaredLogger, error) {
	env := os.Getenv("ENV")
	if env == "production" {
		return prodLogger().Sugar(), nil
	}
	return devLogger().Sugar(), nil
}

func prodLogger() *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()

	cfg.CallerKey = "caller"
	cfg.MessageKey = "message"
	cfg.NameKey = "name"
	cfg.LevelKey = "level"
	cfg.StacktraceKey = "stacktrace"
	cfg.TimeKey = "timestamp"

	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.LineEnding = zapcore.DefaultLineEnding

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     cfg,
		InitialFields:     map[string]interface{}{},
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	return zap.Must(config.Build())
}

func devLogger() *zap.Logger {
	cfg := zap.NewDevelopmentEncoderConfig()

	cfg.CallerKey = "caller"
	cfg.MessageKey = "message"
	cfg.NameKey = "name"
	cfg.LevelKey = "level"
	cfg.StacktraceKey = "stacktrace"
	cfg.TimeKey = "timestamp"

	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.LineEnding = zapcore.DefaultLineEnding

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     cfg,
		InitialFields:     map[string]interface{}{},
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	return zap.Must(config.Build())
}
