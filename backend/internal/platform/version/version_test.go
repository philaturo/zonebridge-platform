package version

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()
	v := Get()
	if v == "" {
		t.Error("Get() returned empty string")
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		content     string
		wantVersion string
		wantErr     bool
	}{
		{"valid version", "2.0.0\n", "2.0.0", false},
		{"version with spaces", "  2.1.0  ", "2.1.0", false},
		{"empty file", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			// Create a temporary file for isolated testing
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "VERSION")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			got, err := LoadFromFile(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantVersion {
				t.Errorf("LoadFromFile() = %q, want %q", got, tt.wantVersion)
			}
		})
	}
}