package member

import (
	"fmt"
	"time"
)

// ID is the unique identifier for a Member within ZoneBridge.
// It is the canonical reference used by every future business domain.
type ID string

// Empty reports whether the ID is unset.
func (id ID) Empty() bool {
	return id == ""
}

// Member is the aggregate root of the Identity domain.
//
// A Member represents an authenticated participant within the ZoneBridge community.
// Every authenticated person is represented by exactly one Member.
// The Member owns its lifecycle state and enforces all transition invariants.
//
// Fields are unexported to preserve aggregate invariants. All mutations occur
// through methods that validate state transitions before applying changes.
type Member struct {
	id          ID
	status      Status
	providerID  string
	email       string
	displayName string
	createdAt   time.Time
	updatedAt   time.Time
	version     int64
}

// NewMember constructs a new Member in the Active state.
//
// A Member is created only after successful authentication. The initial lifecycle
// state is always Active; the Guest state represents the absence of a Member.
//
// The now parameter is the authoritative timestamp for creation. Callers (typically
// services) are responsible for sourcing this from the platform clock abstraction.
func NewMember(id ID, providerID, email, displayName string, now time.Time) (*Member, error) {
	if id.Empty() {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidMember)
	}
	if providerID == "" {
		return nil, fmt.Errorf("%w: providerID is required", ErrInvalidMember)
	}
	if now.IsZero() {
		return nil, fmt.Errorf("%w: creation timestamp is required", ErrInvalidMember)
	}

	return &Member{
		id:          id,
		status:      StatusActive,
		providerID:  providerID,
		email:       email,
		displayName: displayName,
		createdAt:   now,
		updatedAt:   now,
		version:     1,
	}, nil
}

// Rehydrate reconstructs a Member from persisted state.
//
// This constructor is intended for repository implementations that need to
// reconstitute a Member from storage. It performs validation but does not
// enforce lifecycle semantics beyond status validity.
func Rehydrate(id ID, status Status, providerID, email, displayName string, createdAt, updatedAt time.Time, version int64) (*Member, error) {
	if id.Empty() {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidMember)
	}
	if !status.Valid() {
		return nil, fmt.Errorf("%w: status %q is not valid", ErrInvalidMember, status)
	}
	if version < 1 {
		return nil, fmt.Errorf("%w: version must be at least 1", ErrInvalidMember)
	}

	return &Member{
		id:          id,
		status:      status,
		providerID:  providerID,
		email:       email,
		displayName: displayName,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		version:     version,
	}, nil
}

// ID returns the Member's unique identifier.
func (m *Member) ID() ID { return m.id }

// Status returns the Member's current lifecycle state.
func (m *Member) Status() Status { return m.status }

// ProviderID returns the external identity provider's identifier for this Member.
func (m *Member) ProviderID() string { return m.providerID }

// Email returns the Member's email address.
func (m *Member) Email() string { return m.email }

// DisplayName returns the Member's display name.
func (m *Member) DisplayName() string { return m.displayName }

// CreatedAt returns the timestamp when the Member was created.
func (m *Member) CreatedAt() time.Time { return m.createdAt }

// UpdatedAt returns the timestamp of the most recent mutation.
func (m *Member) UpdatedAt() time.Time { return m.updatedAt }

// Version returns the optimistic concurrency version.
// Repositories use this to detect conflicting updates.
func (m *Member) Version() int64 { return m.version }

// IsActive reports whether the Member is in the Active state.
func (m *Member) IsActive() bool { return m.status == StatusActive }

// IsSuspended reports whether the Member is in the Suspended state.
func (m *Member) IsSuspended() bool { return m.status == StatusSuspended }

// IsArchived reports whether the Member is in the Archived state.
func (m *Member) IsArchived() bool { return m.status == StatusArchived }

// Suspend transitions the Member from Active to Suspended.
//
// Suspension preserves identity while temporarily restricting participation.
// The now parameter is the authoritative timestamp for the transition.
//
// Returns ErrInvalidTransition if the Member is not currently Active.
func (m *Member) Suspend(now time.Time) error {
	if m.status != StatusActive {
		return fmt.Errorf(
			"%w: cannot suspend member in %q state (expected %q)",
			ErrInvalidTransition, m.status, StatusActive,
		)
	}
	m.status = StatusSuspended
	m.updatedAt = now
	m.version++
	return nil
}

// Reactivate transitions the Member from Suspended to Active.
//
// Reactivation restores full participation privileges.
// The now parameter is the authoritative timestamp for the transition.
//
// Returns ErrInvalidTransition if the Member is not currently Suspended.
func (m *Member) Reactivate(now time.Time) error {
	if m.status != StatusSuspended {
		return fmt.Errorf(
			"%w: cannot reactivate member in %q state (expected %q)",
			ErrInvalidTransition, m.status, StatusSuspended,
		)
	}
	m.status = StatusActive
	m.updatedAt = now
	m.version++
	return nil
}

// Archive transitions the Member from Active to Archived.
//
// Archival is permanent. An Archived Member remains in the system for historical
// integrity but can never participate or be reactivated.
// The now parameter is the authoritative timestamp for the transition.
//
// Returns ErrInvalidTransition if the Member is not currently Active.
func (m *Member) Archive(now time.Time) error {
	if m.status != StatusActive {
		return fmt.Errorf(
			"%w: cannot archive member in %q state (expected %q)",
			ErrInvalidTransition, m.status, StatusActive,
		)
	}
	m.status = StatusArchived
	m.updatedAt = now
	m.version++
	return nil
}
