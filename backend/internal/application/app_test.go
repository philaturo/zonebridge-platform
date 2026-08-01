package application

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Use a standard valid port. New() only validates config, it doesn't bind.
	t.Setenv("ZONEBRIDGE_SERVER_PORT", "8080")
	t.Setenv("ZONEBRIDGE_LOG_LEVEL", "error")
	t.Setenv("ZONEBRIDGE_LOG_FORMAT", "text")
	t.Setenv("ZONEBRIDGE_ENVIRONMENT", "test")

	app, err := New()
	if err != nil {
		t.Fatalf("expected New() to succeed, got error: %v", err)
	}
	if app == nil {
		t.Fatal("expected app to be non-nil")
	}
	if app.cfg == nil || app.logger == nil || app.clock == nil || app.srv == nil {
		t.Fatal("expected all platform dependencies to be initialized")
	}
}

func TestNew_ConfigurationFailure(t *testing.T) {
	// Use an invalid port to force config.Load() to fail
	t.Setenv("ZONEBRIDGE_SERVER_PORT", "99999")

	app, err := New()
	if err == nil {
		t.Fatal("expected New() to fail with invalid configuration")
	}
	if app != nil {
		t.Fatal("expected app to be nil on failure")
	}
}

func TestRun_GracefulShutdown(t *testing.T) {
	// Use a high ephemeral port to avoid "address already in use" conflicts during testing
	t.Setenv("ZONEBRIDGE_SERVER_PORT", "59999")
	t.Setenv("ZONEBRIDGE_LOG_LEVEL", "error")
	t.Setenv("ZONEBRIDGE_LOG_FORMAT", "text")
	t.Setenv("ZONEBRIDGE_ENVIRONMENT", "test")
	t.Setenv("ZONEBRIDGE_SERVER_SHUTDOWN_TIMEOUT", "2s")

	app, err := New()
	if err != nil {
		t.Fatalf("expected New() to succeed, got error: %v", err)
	}

	// Use a cancellable context to simulate a shutdown trigger deterministically
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx)
	}()

	// Give the server a moment to start listening
	time.Sleep(100 * time.Millisecond)

	// Trigger shutdown by canceling the context
	cancel()

	// Wait for Run() to return cleanly
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("expected Run() to return nil on graceful shutdown, got: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("expected Run() to return within shutdown timeout")
	}
}
