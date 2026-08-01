// Package errors provides sentinel errors for platform-level concerns.
package errors

import "errors"

// Platform-level sentinel errors.
var (
	// ErrInvalidConfiguration indicates that the provided configuration is invalid.
	ErrInvalidConfiguration = errors.New("invalid configuration")

	// ErrVersionNotFound indicates that the application version could not be determined.
	ErrVersionNotFound = errors.New("version information not found")

	// ErrServerStartup indicates that the server failed to start.
	ErrServerStartup = errors.New("server startup failed")

	// ErrGracefulShutdown indicates that graceful shutdown failed.
	ErrGracefulShutdown = errors.New("graceful shutdown failed")
)
