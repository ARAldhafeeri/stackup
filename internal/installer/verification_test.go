package installer

import (
	"os/exec"
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
	"github.com/araldhafeeri/stackup/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestVerifyTool(t *testing.T) {
	sys := &domain.System{OS: "linux"}
	console := ui.NewConsole()
	cfg := &config.Config{}
	installer := New(cfg, sys, console)

	// check if Go is installed for testing
	goInstalled := exec.Command("go", "--version").Run() == nil

	// Define test cases
	// Each test case includes the tool configuration and expected outcome

	tests := []struct {
		name        string
		tool        *config.Tool
		expectError bool
		description string
	}{
		{
			name: "Custom Verify Command - go",
			tool: &config.Tool{
				Name:          "go",
				VerifyCommand: "go version",
			},
			expectError: goInstalled,
			description: "Should pass if Go is installed",
		},
		{
			name: "Default Verify - go",
			tool: &config.Tool{
				Name:          "go",
				VerifyCommand: "",
			},
			expectError: !goInstalled,
			description: "Should try 'go --version' by default",
		},
		{
			name: "Non-existent Tool",
			tool: &config.Tool{
				Name:          "nonexistent-tool-12345",
				VerifyCommand: "",
			},
			expectError: true,
			description: "Should fail for non-existent tool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installer.verifyTool(tt.tool)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}
