package member

import "errors"

// Sentinel errors for the Member aggregate.
var (
	// ErrInvalidMember indicates that a Member failed validation during construction.
	ErrInvalidMember = errors.New("invalid member")

	// ErrInvalidTransition indicates that a lifecycle state transition is not permitted
	// from the Member's current state.
	ErrInvalidTransition = errors.New("invalid lifecycle transition")

	// ErrNotFound indicates that a Member with the given identifier does not exist.
	ErrNotFound = errors.New("member not found")
)
