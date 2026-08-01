// Package session provides the Session entity for the Identity domain.
//
// A Session represents an authenticated interaction between a Member and ZoneBridge.
// Sessions are transient; Members are permanent. Sessions preserve authenticated
// continuity across requests but never become the canonical identity.
//
// This file defines only the Session entity. Session lifecycle management
// (creation, validation, expiration, revocation) is implemented in the
// Identity Services phase.
package session

import (
	"fmt"
	"time"

	"github.com/philaturo/zonebridge-platform/internal/identity/member"
)

// ID is the unique identifier for a Session.
type ID string

// Empty reports whether the ID is unset.
func (id ID) Empty() bool {
	return id == ""
}

// Status represents the lifecycle state of a Session.
type Status string

// Valid Session lifecycle statuses.
const (
	// StatusCreated represents a session that has been established but is not yet active.
	StatusCreated Status = "created"

	// StatusActive represents a valid session that may be used for authentication.
	StatusActive Status = "active"

	// StatusExpired represents a session whose lifetime has elapsed.
	StatusExpired Status = "expired"

	// StatusRevoked represents a session that has been explicitly invalidated.
	StatusRevoked Status = "revoked"
)

// Valid reports whether the status is a recognized session state.
func (s Status) Valid() bool {
	switch s {
	case StatusCreated, StatusActive, StatusExpired, StatusRevoked:
		return true
	}
	return false
}

// Session represents an authenticated interaction between a Member and ZoneBridge.
//
// Every Session belongs to exactly one Member. Sessions may expire or be revoked.
// A Session never grants permissions; it only preserves authenticated continuity.
type Session struct {
	id         ID
	memberID   member.ID
	status     Status
	ipAddress  string
	userAgent  string
	createdAt  time.Time
	expiresAt  time.Time
	revokedAt  *time.Time
	lastUsedAt time.Time
	version    int64
}

// New constructs a new Session in the Active state.
//
// The session is immediately usable for authenticated requests. The expiresAt
// timestamp determines when the session will naturally expire.
func New(id ID, memberID member.ID, ipAddress, userAgent string, createdAt, expiresAt time.Time) (*Session, error) {
	if id.Empty() {
		return nil, fmt.Errorf("session id is required")
	}
	if memberID.Empty() {
		return nil, fmt.Errorf("member id is required")
	}
	if expiresAt.Before(createdAt) || expiresAt.Equal(createdAt) {
		return nil, fmt.Errorf("session expiration must be after creation")
	}

	return &Session{
		id:         id,
		memberID:   memberID,
		status:     StatusActive,
		ipAddress:  ipAddress,
		userAgent:  userAgent,
		createdAt:  createdAt,
		expiresAt:  expiresAt,
		lastUsedAt: createdAt,
		version:    1,
	}, nil
}

// Rehydrate reconstructs a Session from persisted state.
func Rehydrate(id ID, memberID member.ID, status Status, ipAddress, userAgent string, createdAt, expiresAt time.Time, revokedAt *time.Time, lastUsedAt time.Time, version int64) (*Session, error) {
	if id.Empty() {
		return nil, fmt.Errorf("session id is required")
	}
	if memberID.Empty() {
		return nil, fmt.Errorf("member id is required")
	}
	if !status.Valid() {
		return nil, fmt.Errorf("invalid session status: %q", status)
	}

	return &Session{
		id:         id,
		memberID:   memberID,
		status:     status,
		ipAddress:  ipAddress,
		userAgent:  userAgent,
		createdAt:  createdAt,
		expiresAt:  expiresAt,
		revokedAt:  revokedAt,
		lastUsedAt: lastUsedAt,
		version:    version,
	}, nil
}

// ID returns the session's unique identifier.
func (s *Session) ID() ID { return s.id }

// MemberID returns the identifier of the Member who owns this session.
func (s *Session) MemberID() member.ID { return s.memberID }

// Status returns the session's current lifecycle state.
func (s *Session) Status() Status { return s.status }

// IPAddress returns the IP address from which the session was created.
func (s *Session) IPAddress() string { return s.ipAddress }

// UserAgent returns the user agent string from which the session was created.
func (s *Session) UserAgent() string { return s.userAgent }

// CreatedAt returns the timestamp when the session was created.
func (s *Session) CreatedAt() time.Time { return s.createdAt }

// ExpiresAt returns the timestamp when the session will naturally expire.
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }

// RevokedAt returns the timestamp when the session was explicitly revoked.
// Returns the zero time if the session has not been revoked.
func (s *Session) RevokedAt() time.Time {
	if s.revokedAt == nil {
		return time.Time{}
	}
	return *s.revokedAt
}

// LastUsedAt returns the timestamp of the most recent session usage.
func (s *Session) LastUsedAt() time.Time { return s.lastUsedAt }

// Version returns the optimistic concurrency version.
func (s *Session) Version() int64 { return s.version }

// IsExpired reports whether the session has expired as of the given time.
func (s *Session) IsExpired(now time.Time) bool {
	return !now.Before(s.expiresAt)
}

// IsRevoked reports whether the session has been explicitly revoked.
func (s *Session) IsRevoked() bool {
	return s.status == StatusRevoked
}

// IsActive reports whether the session is currently valid for use.
// A session is active if it has not been revoked and has not expired.
func (s *Session) IsActive(now time.Time) bool {
	return s.status == StatusActive && !s.IsExpired(now)
}
