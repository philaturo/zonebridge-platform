// Package member provides the Member aggregate for the Identity domain.
// The Member is the canonical identity entity within ZoneBridge.
package member

// Status represents the operational lifecycle state of a Member.
// A Member exists in exactly one status at any time.
//
// Note: "Guest" is not a Member status. A Guest has not yet been
// authenticated and therefore possesses no Member identity.
// The Member aggregate begins its lifecycle in the Active state
// upon successful authentication.
type Status string

// Valid Member lifecycle statuses.
const (
	// StatusActive represents a verified Member permitted to participate.
	StatusActive Status = "active"

	// StatusSuspended represents a Member whose participation is temporarily restricted.
	StatusSuspended Status = "suspended"

	// StatusArchived represents a permanently inactive Member retained for historical integrity.
	StatusArchived Status = "archived"
)

// Valid reports whether the status is a recognized Member lifecycle state.
func (s Status) Valid() bool {
	switch s {
	case StatusActive, StatusSuspended, StatusArchived:
		return true
	}
	return false
}

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}
