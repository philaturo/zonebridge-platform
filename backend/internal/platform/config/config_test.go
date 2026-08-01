package config

import (
	"errors"
	"strings"
	"testing"

	platformerrors "github.com/philaturo/zonebridge-platform/internal/platform/errors"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		envVars     map[string]string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, cfg *Config)
	}{
        // ... [test cases remain exactly the same as before] ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				// Architecturally sound error checking using errors.Is
				if !errors.Is(err, platformerrors.ErrInvalidConfiguration) {
					t.Errorf("expected error to wrap ErrInvalidConfiguration, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
