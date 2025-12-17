package version

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Error("Version should not be empty")
	}

	if v != Version {
		t.Errorf("GetVersion() = %q, want %q", v, Version)
	}
}

func TestGetFullVersion(t *testing.T) {
	// Save original values
	originalCommit := GitCommit
	originalVersion := Version

	defer func() {
		GitCommit = originalCommit
	}()

	tests := []struct {
		name      string
		gitCommit string
		contains  []string
	}{
		{
			name:      "With Git Commit",
			gitCommit: "abc123",
			contains:  []string{originalVersion, "abc123"},
		},
		{
			name:      "Without Git Commit",
			gitCommit: "unknown",
			contains:  []string{originalVersion},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GitCommit = tt.gitCommit
			result := GetFullVersion()

			for _, needle := range tt.contains {
				if !strings.Contains(result, needle) {
					t.Errorf("GetFullVersion() = %q should contain %q", result, needle)
				}
			}
		})
	}
}

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()

	if info == "" {
		t.Error("BuildInfo should not be empty")
	}

	requiredFields := []string{"Version:", "Git Commit:", "Build Date:"}
	for _, field := range requiredFields {
		if !strings.Contains(info, field) {
			t.Errorf("BuildInfo should contain %q", field)
		}
	}
}
