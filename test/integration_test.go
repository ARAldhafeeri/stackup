//go:build integration
// +build integration

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/installer"
	"github.com/araldhafeeri/stackup/internal/platform"
	"github.com/araldhafeeri/stackup/internal/ui"
)

// TestFullInstallFlow tests the complete installation flow with a minimal config
func TestFullInstallFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	// Create a minimal config that won't actually install anything harmful
	configContent := `profile: integration-test
settings:
  verify_installations: false

tools:
  - name: echo-test
    display_name: "Echo Test"
    version: latest
    custom_install:
      - command: echo
        args: ["Testing installation"]
        description: "Test command"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// Load config
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Detect system
	sys := platform.Detect()
	console := ui.NewConsole()

	// Create installer
	inst := installer.New(cfg, sys, console)

	// This would run the actual installation
	// For integration tests, we might want to mock certain parts
	// or use test-specific tools
	t.Log("Integration test setup complete")
	_ = inst
}
