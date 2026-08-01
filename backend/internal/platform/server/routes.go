package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/philaturo/zonebridge-platform/internal/platform/clock"
	"github.com/philaturo/zonebridge-platform/internal/platform/server/handlers"
)

// RegisterRoutes attaches all platform and domain routes to the router.
func RegisterRoutes(r *chi.Mux, c clock.Clock) {
	// Platform routes
	r.Get("/health", handlers.Health(c))

	// Future domain routes will be registered here, e.g.:
	// r.Route("/api/v1", func(r chi.Router) {
	//     identity.RegisterRoutes(r)
	// })
}
