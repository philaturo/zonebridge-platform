// Package scope provides the Scope entity for the Identity domain.
//
// A Scope defines where a permission applies. Authorization decisions are
// always evaluated within a specific scope, enabling fine-grained capability
// control (e.g., a Moderator may moderate only their assigned Community).
package scope

import "fmt"

// Type represents the category of resource a Scope applies to.
//
// Future domains may introduce additional scope types without modifying
// the authorization architecture.
type Type string

// Valid Scope types.
const (
	// TypePlatform represents platform-wide scope.
	TypePlatform Type = "platform"

	// TypeCommunity represents a specific Community.
	TypeCommunity Type = "community"

	// TypeEvent represents a specific Event.
	TypeEvent Type = "event"

	// TypeKnowledgeArticle represents a specific Knowledge Article.
	TypeKnowledgeArticle Type = "knowledge.article"
)

// Valid reports whether the type is a recognized scope type.
func (t Type) Valid() bool {
	switch t {
	case TypePlatform, TypeCommunity, TypeEvent, TypeKnowledgeArticle:
		return true
	}
	return false
}

// String returns the string representation of the scope type.
func (t Type) String() string {
	return string(t)
}

// Scope defines where a permission applies.
//
// A Scope combines a Type with an optional ResourceID. Platform-wide scopes
// have an empty ResourceID; resource-specific scopes require one.
type Scope struct {
	scopeType  Type
	resourceID string
}

// New constructs a new Scope.
//
// Platform scopes (TypePlatform) must have an empty resourceID.
// Resource-specific scopes must have a non-empty resourceID.
func New(scopeType Type, resourceID string) (*Scope, error) {
	if !scopeType.Valid() {
		return nil, fmt.Errorf("invalid scope type: %q", scopeType)
	}

	if scopeType == TypePlatform && resourceID != "" {
		return nil, fmt.Errorf("platform scope must have empty resourceID, got %q", resourceID)
	}

	if scopeType != TypePlatform && resourceID == "" {
		return nil, fmt.Errorf("scope type %q requires a non-empty resourceID", scopeType)
	}

	return &Scope{
		scopeType:  scopeType,
		resourceID: resourceID,
	}, nil
}

// Platform constructs a platform-wide Scope.
func Platform() *Scope {
	return &Scope{scopeType: TypePlatform, resourceID: ""}
}

// ForResource constructs a resource-specific Scope.
func ForResource(scopeType Type, resourceID string) (*Scope, error) {
	return New(scopeType, resourceID)
}

// Type returns the scope's type.
func (s *Scope) Type() Type { return s.scopeType }

// ResourceID returns the scope's resource identifier.
// Returns empty string for platform-wide scopes.
func (s *Scope) ResourceID() string { return s.resourceID }

// IsPlatform reports whether this is a platform-wide scope.
func (s *Scope) IsPlatform() bool { return s.scopeType == TypePlatform }

// Equal reports whether two scopes are identical.
func (s *Scope) Equal(other *Scope) bool {
	if s == nil || other == nil {
		return s == other
	}
	return s.scopeType == other.scopeType && s.resourceID == other.resourceID
}
