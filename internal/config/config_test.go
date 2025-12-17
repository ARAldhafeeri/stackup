package config

import (
	"testing"
)

func TestToolGetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		tool     Tool
		expected string
	}{
		{
			name: "With DisplayName",
			tool: Tool{
				Name:        "wsl",
				DisplayName: "Windows Subsystem for Linux",
			},
			expected: "Windows Subsystem for Linux",
		},
		{
			name: "Without DisplayName",
			tool: Tool{
				Name:        "git",
				DisplayName: "",
			},
			expected: "git",
		},
		{
			name: "Empty Name and DisplayName",
			tool: Tool{
				Name:        "",
				DisplayName: "",
			},
			expected: "ToolNameNotSet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tool.GetDisplayName()
			if result != tt.expected {
				t.Errorf("GetDisplayName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestToolGetPlatformConfig(t *testing.T) {
	tool := Tool{
		Name: "test-tool",
		Windows: &PlatformConfig{
			Installer: "<https://example.com/windows.exe>",
			Type:      "exe",
		},
		Linux: &PlatformConfig{
			PackageNames: map[string]string{"apt": "test-package"},
		},
		MacOS: &PlatformConfig{
			Brew: "test-tool",
		},
	}

	tests := []struct {
		name     string
		osName   string
		expected *PlatformConfig
	}{
		{
			name:     "Windows Platform",
			osName:   "windows",
			expected: tool.Windows,
		},
		{
			name:     "Linux Platform",
			osName:   "linux",
			expected: tool.Linux,
		},
		{
			name:     "macOS Platform",
			osName:   "darwin",
			expected: tool.MacOS,
		},
		{
			name:     "Unknown Platform",
			osName:   "unknown",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.GetPlatformConfig(tt.osName)
			if result != tt.expected {
				t.Errorf("GetPlatformConfig(%s) = %v, want %v", tt.osName, result, tt.expected)
			}
		})
	}
}

func TestCommandStructure(t *testing.T) {
	cmd := Command{
		Command:     "wsl",
		Args:        []string{"--install"},
		Description: "Install WSL",
		Sudo:        false,
		WaitFor:     10,
		IgnoreError: false,
	}

	if cmd.Command != "wsl" {
		t.Errorf("Command = %q, want %q", cmd.Command, "wsl")
	}

	if len(cmd.Args) != 1 || cmd.Args[0] != "--install" {
		t.Errorf("Args = %v, want [\"--install\"]", cmd.Args)
	}

	if cmd.WaitFor != 10 {
		t.Errorf("WaitFor = %d, want 10", cmd.WaitFor)
	}

	if cmd.Sudo {
		t.Error("Sudo should be false")
	}

	if cmd.IgnoreError {
		t.Error("IgnoreError should be false")
	}
}
