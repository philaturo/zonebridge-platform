// Package logger provides structured logging infrastructure for the ZoneBridge platform.
// It uses Go's standard log/slog package to ensure zero external dependencies
// while providing production-ready, machine-readable diagnostic output.
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/philaturo/zonebridge-platform/internal/platform/config"
)

// New creates a new structured logger based on the provided configuration.
// It configures the log level and output format (JSON or text).
// It returns an error if the configuration contains an invalid log level or format.
func New(cfg *config.Config) (*slog.Logger, error) {
	level, err := parseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}

	switch strings.ToLower(cfg.LogFormat) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		return nil, fmt.Errorf("unsupported log format: %q", cfg.LogFormat)
	}

	return slog.New(handler), nil
}

// parseLevel converts a string log level to slog.Level.
// It returns an error if the level is not one of: debug, info, warn, error.
func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid log level: %q (must be debug, info, warn, or error)", level)
	}
}