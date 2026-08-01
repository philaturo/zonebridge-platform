package server

import (
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	platformmiddleware "github.com/philaturo/zonebridge-platform/internal/platform/server/middleware"
)

// NewRouter initializes and configures the Chi router with platform middleware.
// The middleware order is strictly enforced to ensure correct context propagation
// and observability: RequestID -> RealIP -> Logger -> Recoverer -> Timeout.
func NewRouter(logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// 1. RequestID: Must be first to generate ID for all subsequent middleware
	r.Use(platformmiddleware.RequestID)

	// 2. RealIP: Resolve true client IP before logging
	r.Use(platformmiddleware.RealIP)

	// 3. Logger: Log the request with RequestID and RealIP already in context
	r.Use(platformmiddleware.Logger(logger))

	// 4. Recoverer: Catch panics and log them with full context (including RequestID)
	r.Use(platformmiddleware.Recoverer(logger))

	// 5. Timeout: Enforce request deadline (using Chi's robust, tested implementation - Chi's timeout requires a duration)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	return r
}
