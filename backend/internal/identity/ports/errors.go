// Package ports defines the contracts between the Identity domain and external systems.
//
// Ports are interfaces that the domain depends upon. Infrastructure adapters
// (databases, OAuth providers, caches) implement these interfaces. The domain
// never depends on adapter implementations.
package ports

import "errors"

// Common sentinel errors used across repository ports.
var (
	// ErrNotFound indicates that the requested entity does not exist.
	ErrNotFound = errors.New("entity not found")

	// ErrAlreadyExists indicates that an entity with the same unique identifier already exists.
	ErrAlreadyExists = errors.New("entity already exists")

	// ErrConflict indicates an optimistic concurrency conflict.
	ErrConflict = errors.New("concurrent modification detected")
)
