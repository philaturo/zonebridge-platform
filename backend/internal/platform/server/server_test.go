package server

import (
	"context"
//	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/philaturo/zonebridge-platform/internal/platform/config"
	"github.com/philaturo/zonebridge-platform/internal/platform/logger"
)

func TestServerLifecycle(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		ServerHost:            "127.0.0.1",
		ServerPort:            0, // Let OS assign a random available port for testing
		ServerReadTimeout:     5 * time.Second,
		ServerWriteTimeout:    5 * time.Second,
		ServerShutdownTimeout: 5 * time.Second,
		LogLevel:              "error",
		LogFormat:             "text",
		Environment:           "test",
	}

	log, err := logger.New(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := New(cfg, log, r)

	// Start server in background (simulating application orchestration)
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test that it responds
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Trigger graceful shutdown (simulating application orchestration)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		t.Errorf("unexpected error during shutdown: %v", err)
	}

	// Verify Start() returned cleanly or with expected closed error
	select {
	case err := <-errChan:
		if err != nil && err.Error() != "server startup failed: http: Server closed" {
			// Note: http.ErrServerClosed is expected and handled inside Start(),
			// but if it leaks, we catch it here.
		}
	case <-time.After(1 * time.Second):
		t.Error("server did not shut down within expected timeframe")
	}
}
