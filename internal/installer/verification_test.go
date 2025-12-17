package installer

import (
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

	tests := []struct {
		name        string
		tool        *config.Tool
		shouldPass  bool
		description string
	}{
		{
			name: "Custom Verify Command - go",
			tool: &config.Tool{
				Name:          "go",
				VerifyCommand: "go version",
			},
			description: "Should pass if Go is installed",
		},
		{
			name: "Default Verify - go",
			tool: &config.Tool{
				Name:          "go",
				VerifyCommand: "",
			},
			description: "Should try 'go --version' by default",
		},
		{
			name: "Non-existent Tool",
			tool: &config.Tool{
				Name:          "nonexistent-tool-12345",
				VerifyCommand: "",
			},
			shouldPass:  false,
			description: "Should fail for non-existent tool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := installer.verifyTool(tt.tool)
			passed := err == nil

			assert.Equal(t, tt.shouldPass, passed, "Verification result mismatch")

		})
	}
}
