// Package role provides the Role entity for the Identity domain.
//
// A Role groups permissions into reusable capability sets. Roles simplify
// authorization by assigning collections of permissions rather than individual
// capabilities. A Role contains no business logic; it is a pure data structure
// that groups immutable Permission IDs.
package role

import (
	"fmt"

	"github.com/philaturo/zonebridge-platform/internal/identity/permission"
)

// ID is the unique identifier for a Role.
type ID string

// Empty reports whether the ID is unset.
func (id ID) Empty() bool {
	return id == ""
}

// Well-known platform roles.
//
// These constants define the canonical roles recognized by ZoneBridge.
// Platform governance is implemented through roles and permissions rather
// than privileged account types.
const (
	// RolePlatformMaintainer represents a Platform Maintainer with elevated governance permissions.
	RolePlatformMaintainer ID = "platform.maintainer"

	// RoleCommunityOwner represents the owner of a specific Community.
	RoleCommunityOwner ID = "community.owner"

	// RoleCommunityModerator represents a moderator of a specific Community.
	RoleCommunityModerator ID = "community.moderator"

	// RoleMember represents the default role assigned to every authenticated Member.
	RoleMember ID = "member"
)

// Role groups permissions into reusable capability sets.
//
// A Role contains one or more permissions and may be assigned to multiple Members.
// Roles contain no business logic and do not authenticate Members or evaluate permissions.
type Role struct {
	id          ID
	name        string
	description string
	permissions []permission.ID
}

// New constructs a new Role with the given permissions.
//
// Returns an error if the ID is empty, the name is empty, or no permissions are provided.
// Duplicate permission IDs are deduplicated while preserving order.
func New(id ID, name, description string, permissions []permission.ID) (*Role, error) {
	if id.Empty() {
		return nil, fmt.Errorf("role id is required")
	}
	if name == "" {
		return nil, fmt.Errorf("role name is required")
	}
	if len(permissions) == 0 {
		return nil, fmt.Errorf("role must contain at least one permission")
	}

	// Validate and deduplicate permissions.
	seen := make(map[permission.ID]struct{}, len(permissions))
	deduped := make([]permission.ID, 0, len(permissions))
	for _, pid := range permissions {
		if !pid.Valid() {
			return nil, fmt.Errorf("role %q contains invalid permission id: %q", id, pid)
		}
		if _, exists := seen[pid]; exists {
			continue
		}
		seen[pid] = struct{}{}
		deduped = append(deduped, pid)
	}

	// Copy to prevent external mutation.
	perms := make([]permission.ID, len(deduped))
	copy(perms, deduped)

	return &Role{
		id:          id,
		name:        name,
		description: description,
		permissions: perms,
	}, nil
}

// ID returns the role's unique identifier.
func (r *Role) ID() ID { return r.id }

// Name returns the role's display name.
func (r *Role) Name() string { return r.name }

// Description returns a human-readable description of the role.
func (r *Role) Description() string { return r.description }

// Permissions returns a copy of the permission IDs granted by this role.
// The returned slice is a defensive copy; callers cannot mutate the role's state.
func (r *Role) Permissions() []permission.ID {
	out := make([]permission.ID, len(r.permissions))
	copy(out, r.permissions)
	return out
}

// HasPermission reports whether the role grants the specified permission.
func (r *Role) HasPermission(pid permission.ID) bool {
	for _, p := range r.permissions {
		if p == pid {
			return true
		}
	}
	return false
}
