// Package clock provides a time abstraction for deterministic testing.
package clock

import "time"

// Clock provides an abstraction over time operations.
// This allows tests to inject a fake clock for deterministic behavior.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
}

// RealClock implements Clock using the actual system time.
type RealClock struct{}

// Now returns the current system time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// New returns a new RealClock instance.
func New() Clock {
	return RealClock{}
}