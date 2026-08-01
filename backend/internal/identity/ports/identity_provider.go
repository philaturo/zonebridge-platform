package ports

import (
	"context"
	"time"
)

// VerifiedIdentity represents the trusted identity claims returned by an
// Identity Provider after successful authentication.
//
// The Identity domain converts these claims into the canonical Member representation.
// Business domains never consume provider-specific claims directly.
type VerifiedIdentity struct {
	// ProviderID is the unique identifier assigned by the external identity provider.
	ProviderID string

	// Email is the verified email address of the authenticated user.
	Email string

	// DisplayName is the human-readable name of the authenticated user.
	DisplayName string

	// VerifiedAt is the timestamp when the identity was verified by the provider.
	VerifiedAt time.Time
}

// AuthenticationRequest contains the provider-specific material required to
// initiate authentication.
//
// For OAuth flows, Code typically contains the authorization code.
// For development and testing adapters, Code may contain a member identifier.
type AuthenticationRequest struct {
	// Code is the provider-specific authentication material.
	Code string
}

// IdentityProvider verifies identity on behalf of ZoneBridge.
//
// ZoneBridge delegates authentication to approved providers. Authentication
// credentials are never stored by ZoneBridge. The provider returns trusted
// identity claims that the domain uses to resolve or create Members.
//
// Implementations include OAuth adapters (e.g., Gitea) and development mocks.
type IdentityProvider interface {
	// Authenticate verifies the provided credentials and returns a VerifiedIdentity.
	//
	// Returns an error if the credentials are invalid, the provider is unavailable,
	// or the identity claims are malformed. Authentication failures never expose
	// sensitive implementation details.
	Authenticate(ctx context.Context, request AuthenticationRequest) (*VerifiedIdentity, error)

	// Name returns the human-readable name of the identity provider.
	// Used for logging and observability.
	Name() string
}
