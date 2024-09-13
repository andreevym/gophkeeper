package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var log *zap.Logger

const defaultLogLevel = "DEBUG"

func NewLogger(level string) (*zap.Logger, error) {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}
	lvl := zap.NewAtomicLevelAt(atomicLevel.Level())
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}
	log = zl
	return log, nil
}

func Logger() *zap.Logger {
	if log != nil {
		return log
	}

	log, _ = NewLogger(defaultLogLevel)
	log.Debug("debug log checked")
	log.Info("info log checked")
	log.Warn("warn log checked")
	log.Error("error log checked")

	return log
}
