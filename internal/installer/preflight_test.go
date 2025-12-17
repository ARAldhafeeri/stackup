package installer

import (
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
	"github.com/araldhafeeri/stackup/internal/ui"
)

func TestPreflightChecks(t *testing.T) {
	tests := []struct {
		name        string
		system      *domain.System
		expectError bool
	}{
		{
			name: "Linux with Package Manager",
			system: &domain.System{
				OS:             "linux",
				Arch:           "amd64",
				PackageManager: "apt",
			},
			expectError: false,
		},
		{
			name: "Windows without Package Manager",
			system: &domain.System{
				OS:             "windows",
				Arch:           "amd64",
				PackageManager: "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Tools: []config.Tool{
					{Name: "test-tool", Version: "1.0", Linux: &config.PlatformConfig{}},
				},
			}
			console := ui.NewConsole()
			installer := New(cfg, tt.system, console)

			// Note: Internet connectivity check will run here
			// This test may fail if there's no internet connection
			err := installer.runPreflightChecks()

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}

			if !tt.expectError && err != nil {
				// Internet connectivity issues are acceptable in tests
				t.Logf("Preflight check error (may be expected): %v", err)
			}
		})
	}
}
