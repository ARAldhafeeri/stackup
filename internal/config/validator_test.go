package config

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Config",
			config: &Config{
				Tools: []Tool{
					{Name: "git", Version: "latest", Linux: &PlatformConfig{}},
				},
			},
			expectError: false,
		},
		{
			name: "No Tools",
			config: &Config{
				Tools: []Tool{},
			},
			expectError: true,
			errorMsg:    "no tools defined",
		},
		{
			name: "Missing Tool Name",
			config: &Config{
				Tools: []Tool{
					{Name: "", Version: "latest", Linux: &PlatformConfig{}},
				},
			},
			expectError: true,
			errorMsg:    "missing name",
		},
		{
			name: "Duplicate Tool Names",
			config: &Config{
				Tools: []Tool{
					{Name: "git", Version: "latest", Linux: &PlatformConfig{}},
					{Name: "git", Version: "2.0", Linux: &PlatformConfig{}},
				},
			},
			expectError: true,
			errorMsg:    "duplicate tool name",
		},
		{
			name: "Unknown Dependency",
			config: &Config{
				Tools: []Tool{
					{
						Name:         "docker",
						Version:      "latest",
						Dependencies: []string{"nonexistent"},
						Linux:        &PlatformConfig{},
					},
				},
			},
			expectError: true,
			errorMsg:    "unknown dependency",
		},
		{
			name: "Valid Dependencies",
			config: &Config{
				Tools: []Tool{
					{Name: "wsl", Version: "latest", Windows: &PlatformConfig{}},
					{
						Name:         "docker",
						Version:      "latest",
						Dependencies: []string{"wsl"},
						Windows:      &PlatformConfig{},
					},
				},
			},
			expectError: false,
		},
		{
			name: "No Platform Configuration",
			config: &Config{
				Tools: []Tool{
					{Name: "tool", Version: "latest"},
				},
			},
			expectError: true,
			errorMsg:    "no platform configuration",
		},
		{
			name: "Tool with Custom Install Only",
			config: &Config{
				Tools: []Tool{
					{
						Name:    "tool",
						Version: "latest",
						CustomInstall: []Command{
							{Command: "echo", Args: []string{"installing"}},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Error message %q does not contain %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateComplexDependencies(t *testing.T) {
	config := &Config{
		Tools: []Tool{
			{Name: "tool-a", Version: "1.0", Linux: &PlatformConfig{}},
			{Name: "tool-b", Version: "1.0", Dependencies: []string{"tool-a"}, Linux: &PlatformConfig{}},
			{Name: "tool-c", Version: "1.0", Dependencies: []string{"tool-a", "tool-b"}, Linux: &PlatformConfig{}},
		},
	}

	err := Validate(config)
	if err != nil {
		t.Errorf("Valid complex dependencies failed: %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		len(s) > len(substr)+1 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
