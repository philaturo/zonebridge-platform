// Package permission provides the Permission entity for the Identity domain.
//
// A Permission represents a single atomic capability within ZoneBridge.
// Permissions are immutable and are granted to Members through Roles.
// Authorization evaluates capabilities, not identities.
package permission

import (
	"fmt"
	"strings"
)

// ID is the unique identifier for a Permission.
//
// Permission IDs follow a domain.action convention (e.g., "community.create").
// They are immutable once defined.
type ID string

// Empty reports whether the ID is unset.
func (id ID) Empty() bool {
	return id == ""
}

// Valid reports whether the ID conforms to the domain.action convention.
func (id ID) Valid() bool {
	if id.Empty() {
		return false
	}
	s := string(id)
	parts := strings.SplitN(s, ".", 2)
	if len(parts) != 2 {
		return false
	}
	return parts[0] != "" && parts[1] != ""
}

// Domain returns the domain portion of the permission ID.
// Returns empty string if the ID is invalid.
func (id ID) Domain() string {
	parts := strings.SplitN(string(id), ".", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// Action returns the action portion of the permission ID.
// Returns empty string if the ID is invalid.
func (id ID) Action() string {
	parts := strings.SplitN(string(id), ".", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// String returns the string representation of the permission ID.
func (id ID) String() string {
	return string(id)
}

// Well-known platform permissions.
//
// These constants define the canonical capabilities recognized by ZoneBridge.
// Future milestones may introduce additional permissions within their respective
// domains without modifying this authorization architecture.
const (
	// Community permissions.
	CommunityCreate           ID = "community.create"
	CommunityDelete           ID = "community.delete"
	CommunityModerate         ID = "community.moderate"

	// Event permissions.
	EventCreate               ID = "event.create"
	EventPublish              ID = "event.publish"
	EventDelete               ID = "event.delete"

	// Knowledge permissions.
	KnowledgePublish          ID = "knowledge.publish"
	KnowledgeModerate         ID = "knowledge.moderate"

	// Help request permissions.
	HelpRequestCreate         ID = "help.request.create"
	HelpRequestRespond        ID = "help.request.respond"

	// Audit coordination permissions.
	AuditRequest              ID = "audit.request"
	AuditRespond              ID = "audit.respond"

	// Platform governance permissions.
	PlatformMaintain          ID = "platform.maintain"
	PlatformModerate          ID = "platform.moderate"
	MemberSuspend             ID = "member.suspend"
	MemberArchive             ID = "member.archive"
)

// Permission is an immutable capability within the platform.
//
// Permissions are grouped into Roles and evaluated during authorization.
// A Permission never authenticates Members, owns resources, or evaluates itself.
type Permission struct {
	id          ID
	description string
}

// New constructs a new Permission.
//
// Returns an error if the ID is not a valid domain.action identifier.
func New(id ID, description string) (*Permission, error) {
	if !id.Valid() {
		return nil, fmt.Errorf("invalid permission id: %q (must follow domain.action convention)", id)
	}
	return &Permission{
		id:          id,
		description: description,
	}, nil
}

// ID returns the permission's unique identifier.
func (p *Permission) ID() ID { return p.id }

// Description returns a human-readable description of the capability.
func (p *Permission) Description() string { return p.description }
