package ports

import (
	"context"

	"github.com/philaturo/zonebridge-platform/internal/identity/member"
)

// MemberRepository defines the persistence contract for Member aggregates.
//
// The Identity domain depends upon this interface. Infrastructure adapters
// (e.g., PostgreSQL, in-memory) implement it. The domain never depends on
// adapter implementations.
type MemberRepository interface {
	// FindByID retrieves a Member by its unique identifier.
	//
	// Returns ErrNotFound if no Member with the given ID exists.
	FindByID(ctx context.Context, id member.ID) (*member.Member, error)

	// FindByProviderID retrieves a Member by their external identity provider identifier.
	//
	// Returns ErrNotFound if no Member with the given provider ID exists.
	// This is used during authentication to locate existing Members.
	FindByProviderID(ctx context.Context, providerID string) (*member.Member, error)

	// Save persists a Member aggregate.
	//
	// If the Member is new (not previously persisted), it is created.
	// If the Member exists, it is updated. Implementations should use the
	// Member's version field for optimistic concurrency control.
	//
	// Returns ErrConflict if a concurrent modification is detected.
	Save(ctx context.Context, m *member.Member) error
}
