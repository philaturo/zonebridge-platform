package logger

import (
	"context"
	"log/slog"
)

// ctxKey is an unexported type to prevent context key collisions.
type ctxKey struct{}

// loggerKey is the specific key used to store the logger in a context.
var loggerKey = ctxKey{}

// WithContext returns a new context with the provided logger attached.
// This enables request-scoped logging where handlers can attach metadata
// (like request IDs) before passing the context down to services.
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves the logger from the context.
// If no logger is found, it returns nil. Callers should check for nil
// and provide a fallback (e.g., a default logger) if necessary.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return nil
}