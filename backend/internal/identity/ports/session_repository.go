package ports

import (
	"context"

	"github.com/philaturo/zonebridge-platform/internal/identity/member"
	"github.com/philaturo/zonebridge-platform/internal/identity/session"
)

// SessionRepository defines the persistence contract for Session entities.
//
// The Identity domain depends upon this interface. Infrastructure adapters
// implement it. The domain never depends on adapter implementations.
type SessionRepository interface {
	// FindByID retrieves a Session by its unique identifier.
	//
	// Returns ErrNotFound if no Session with the given ID exists.
	FindByID(ctx context.Context, id session.ID) (*session.Session, error)

	// Save persists a Session entity.
	//
	// Implementations should use the Session's version field for optimistic
	// concurrency control. Returns ErrConflict if a concurrent modification
	// is detected.
	Save(ctx context.Context, s *session.Session) error

	// Delete removes a Session from the repository.
	//
	// Returns ErrNotFound if the Session does not exist.
	Delete(ctx context.Context, id session.ID) error

	// DeleteByMemberID removes all Sessions belonging to a specific Member.
	//
	// This is used during Member suspension or archival to invalidate all
	// active sessions. Returns nil if the Member has no sessions.
	DeleteByMemberID(ctx context.Context, memberID member.ID) error
}
