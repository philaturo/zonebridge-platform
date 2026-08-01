// Package version provides application version management.
package version

import (
	"fmt"
	"os"
	"strings"
)

// Version is the application version.
// It should be set at build time using ldflags:
//
//	go build -ldflags "-X github.com/philaturo/zonebridge-platform/internal/platform/version.Version=2.0.0"
var Version = "dev"

func init() {
	// If Version was already set via ldflags, do nothing.
	if Version != "dev" {
		return
	}

	// Fallback for local development: attempt to read from the VERSION file.
	// This avoids fragile absolute path assumptions by relying on the
	// standard practice of running `go run` or the binary from the repository root.
	v, err := LoadFromFile("VERSION")
	if err == nil && v != "" {
		Version = v
	}
}

// Get returns the application version.
func Get() string {
	return Version
}

// LoadFromFile reads and trims the version from the given file path
// This is exported to provide production value and to allow deterministic
// isolated testing without init() hacks.

func LoadFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read version file: %w", err)

	}
	v := strings.TrimSpace(string(data))
	if v == "" {
		return "", fmt.Errorf("version file is empty")
	}
	return v, nil
}
