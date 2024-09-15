package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

var (
	log             *zap.Logger
	once            sync.Once
	defaultLogLevel = "DEBUG"
)

// NewLogger creates a new zap.Logger based on the provided log level.
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
	return zl, nil
}

// Logger returns the global zap.Logger instance. It initializes the logger if not already set.
func Logger() *zap.Logger {
	once.Do(func() {
		var err error
		log, err = NewLogger(defaultLogLevel)
		if err != nil {
			// Fallback to a no-op logger or handle the error appropriately
			fmt.Printf("Error initializing logger: %v\n", err)
			log = zap.NewNop()
		} else {
			// Test logging to ensure logger is working correctly
			log.Debug("debug log checked")
			log.Info("info log checked")
			log.Warn("warn log checked")
			log.Error("error log checked")
		}
	})
	return log
}
