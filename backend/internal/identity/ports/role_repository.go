package ports

import (
	"context"

	"github.com/philaturo/zonebridge-platform/internal/identity/role"
)

// RoleRepository defines the persistence contract for Role entities.
//
// The Identity domain depends upon this interface. Infrastructure adapters
// implement it. The domain never depends on adapter implementations.
type RoleRepository interface {
	// FindByID retrieves a Role by its unique identifier.
	//
	// Returns ErrNotFound if no Role with the given ID exists.
	FindByID(ctx context.Context, id role.ID) (*role.Role, error)

	// FindAll retrieves all Roles defined in the system.
	FindAll(ctx context.Context) ([]*role.Role, error)

	// FindByMemberID retrieves all Roles assigned to a specific Member.
	FindByMemberID(ctx context.Context, memberID string) ([]*role.Role, error)
}
