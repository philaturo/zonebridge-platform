// Package application provides the composition root and lifecycle management
// for the ZoneBridge platform. It assembles platform dependencies and
// coordinates startup and graceful shutdown.
package application

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/philaturo/zonebridge-platform/internal/platform/clock"
	"github.com/philaturo/zonebridge-platform/internal/platform/config"
	"github.com/philaturo/zonebridge-platform/internal/platform/logger"
	"github.com/philaturo/zonebridge-platform/internal/platform/server"
)

// App is the composition root of the ZoneBridge platform.
// It owns the lifecycle of all long-lived platform resources.
type App struct {
	cfg    *config.Config
	logger *slog.Logger
	clock  clock.Clock
	srv    *server.Server
}

// New assembles the platform dependencies and returns a fully initialized App.
// It is strictly deterministic: it performs no network activity, starts no
// goroutines, and handles no OS signals.
func New() (*App, error) {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	// 2. Initialize logger
	lgr, err := logger.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	lgr.Info("configuration loaded")
	lgr.Info("logger initialized")

	// 3. Initialize clock
	clk := clock.New()

	// 4. Construct router and register routes
	router := server.NewRouter(lgr)
	server.RegisterRoutes(router, clk) // <-- FIXED: Removed the logger argument
	lgr.Info("router initialized")

	// 5. Construct HTTP server
	srv := server.New(cfg, lgr, router)
	lgr.Info("server initialized")

	lgr.Info("application initialized")

	return &App{
		cfg:    cfg,
		logger: lgr,
		clock:  clk,
		srv:    srv,
	}, nil
}

// Run owns the lifetime of the application. It starts the HTTP server,
// waits for shutdown triggers (OS signals or context cancellation),
// coordinates graceful shutdown, and blocks until the shutdown is complete.
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("starting http server")

	// Channel to capture server startup errors
	errCh := make(chan error, 1)
	go func() {
		if err := a.srv.Start(); err != nil {
			errCh <- fmt.Errorf("start server: %w", err)
		}
		close(errCh)
	}()

	// Channel to capture OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh) // Ensure cleanup of signal registration

	var shutdownReason string

	select {
	case err := <-errCh:
		// Server failed to start on its own
		return fmt.Errorf("server failed to start: %w", err)
	case sig := <-sigCh:
		shutdownReason = fmt.Sprintf("os signal: %s", sig.String())
	case <-ctx.Done():
		shutdownReason = "context cancelled"
	}

	a.logger.Info("shutdown signal received", "reason", shutdownReason)
	a.logger.Info("graceful shutdown started")

	// Use context.Background() for shutdown to ensure the timeout
	// is respected even if the parent context is already cancelled.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ServerShutdownTimeout)
	defer cancel()

	if err := a.srv.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	a.logger.Info("graceful shutdown completed")
	return nil
}
