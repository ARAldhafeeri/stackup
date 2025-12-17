package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
		validate    func(*testing.T, *Config)
	}{
		{
			name: "Valid Simple Config",
			content: `profile: test
settings:
  verify_installations: true
tools:
  - name: git
    version: latest
`,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Profile != "test" {
					t.Errorf("Profile = %q, want %q", cfg.Profile, "test")
				}
				if len(cfg.Tools) != 1 {
					t.Errorf("len(Tools) = %d, want 1", len(cfg.Tools))
				}
				if cfg.Tools[0].Name != "git" {
					t.Errorf("Tools[0].Name = %q, want %q", cfg.Tools[0].Name, "git")
				}
			},
		},
		{
			name: "Config with Dependencies",
			content: `tools:
  - name: docker
    version: latest
    dependencies: ["wsl"]
  - name: wsl
    version: latest
`,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if len(cfg.Tools) != 2 {
					t.Errorf("len(Tools) = %d, want 2", len(cfg.Tools))
				}
				docker := cfg.Tools[0]
				if len(docker.Dependencies) != 1 || docker.Dependencies[0] != "wsl" {
					t.Errorf("Dependencies = %v, want [\"wsl\"]", docker.Dependencies)
				}
			},
		},
		{
			name: "Config with Custom Commands",
			content: `tools:
  - name: wsl
    version: latest
    windows:
      custom_commands:
        - command: wsl
          args: ["--install"]
          description: "Enable WSL"
          wait_for: 5
          ignore_error: false
`,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				tool := cfg.Tools[0]
				if tool.Windows == nil {
					t.Fatal("Windows config is nil")
				}
				if len(tool.Windows.CustomCommands) != 1 {
					t.Fatalf("len(CustomCommands) = %d, want 1", len(tool.Windows.CustomCommands))
				}
				cmd := tool.Windows.CustomCommands[0]
				if cmd.Command != "wsl" {
					t.Errorf("Command = %q, want %q", cmd.Command, "wsl")
				}
				if cmd.WaitFor != 5 {
					t.Errorf("WaitFor = %d, want 5", cmd.WaitFor)
				}
			},
		},
		{
			name: "Config with Presets",
			content: `tools:
  - name: git
    version: latest
  - name: docker
    version: latest
presets:
  web-dev:
    description: "Web development"
    tools: ["git", "docker"]
`,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if len(cfg.Presets) != 1 {
					t.Errorf("len(Presets) = %d, want 1", len(cfg.Presets))
				}
				preset, exists := cfg.Presets["web-dev"]
				if !exists {
					t.Fatal("web-dev preset not found")
				}
				if preset.Description != "Web development" {
					t.Errorf("Description = %q, want %q", preset.Description, "Web development")
				}
				if len(preset.Tools) != 2 {
					t.Errorf("len(Tools) = %d, want 2", len(preset.Tools))
				}
			},
		},
		{
			name:        "Invalid YAML",
			content:     `invalid: [unclosed array`,
			expectError: true,
		},
		{
			name:        "Empty Config",
			content:     ``,
			expectError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg == nil {
					t.Error("Config should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.yaml")

			if err := os.WriteFile(configPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			cfg, err := LoadFromFile(configPath)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestLoadFromFileNonExistent(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadFromBytes(t *testing.T) {
	content := []byte(`profile: test
tools:
  - name: git
    version: latest
`)

	cfg, err := LoadFromBytes(content)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cfg.Profile != "test" {
		t.Errorf("Profile = %q, want %q", cfg.Profile, "test")
	}

	if len(cfg.Tools) != 1 {
		t.Errorf("len(Tools) = %d, want 1", len(cfg.Tools))
	}

	// Test that DisplayName defaults to Name
	if cfg.Tools[0].GetDisplayName() != "git" {
		t.Errorf("DisplayName = %q, want %q", cfg.Tools[0].GetDisplayName(), "git")
	}
}

func TestComplexConfigParsing(t *testing.T) {
	content := `profile: complex-test

settings:
  auto_update_path: true
  verify_installations: true

tools:
  - name: complex-tool
    display_name: "Complex Tool"
    version: "2.0"
    description: "A complex tool"
    dependencies: ["dependency-tool"]
    requires_reboot: true
    verify_command: "complex-tool --version"

    pre_install:
      - command: echo
        args: ["Pre-install"]
        description: "Preparation"

    custom_install:
      - command: echo
        args: ["Custom install"]
        wait_for: 2

    post_install:
      - command: echo
        args: ["Post-install"]

    windows:
      installer: "<https://example.com/installer.exe>"
      type: exe
      silent_flags: ["/S", "/quiet"]
      custom_commands:
        - command: cmd
          args: ["/c", "echo", "Windows"]

    linux:
      package_names:
        apt: complex-tool
        dnf: complex-tool
      custom_commands:
        - command: bash
          args: ["-c", "echo Linux"]
          sudo: true

    macos:
      brew: complex-tool
      custom_commands:
        - command: bash
          args: ["-c", "echo macOS"]

  - name: dependency-tool
    version: "1.0"
`

	cfg, err := LoadFromBytes([]byte(content))
	if err != nil {
		t.Fatalf("Failed to parse complex config: %v", err)
	}

	tool := cfg.Tools[0]

	// Validate basic fields
	if tool.Name != "complex-tool" {
		t.Errorf("Name = %q, want %q", tool.Name, "complex-tool")
	}

	if tool.DisplayName != "Complex Tool" {
		t.Errorf("DisplayName = %q, want %q", tool.DisplayName, "Complex Tool")
	}

	if !tool.RequiresReboot {
		t.Error("RequiresReboot should be true")
	}

	// Validate dependencies
	if len(tool.Dependencies) != 1 || tool.Dependencies[0] != "dependency-tool" {
		t.Errorf("Dependencies = %v, want [\"dependency-tool\"]", tool.Dependencies)
	}

	// Validate lifecycle commands
	if len(tool.PreInstall) != 1 {
		t.Errorf("len(PreInstall) = %d, want 1", len(tool.PreInstall))
	}

	if len(tool.CustomInstall) != 1 {
		t.Errorf("len(CustomInstall) = %d, want 1", len(tool.CustomInstall))
	}

	if len(tool.PostInstall) != 1 {
		t.Errorf("len(PostInstall) = %d, want 1", len(tool.PostInstall))
	}

	// Validate platform configs
	if tool.Windows == nil {
		t.Fatal("Windows config is nil")
	}

	if tool.Windows.Type != "exe" {
		t.Errorf("Windows.Type = %q, want %q", tool.Windows.Type, "exe")
	}

	if len(tool.Windows.SilentFlags) != 2 {
		t.Errorf("len(Windows.SilentFlags) = %d, want 2", len(tool.Windows.SilentFlags))
	}

	if tool.Linux == nil {
		t.Fatal("Linux config is nil")
	}

	if tool.Linux.PackageNames["apt"] != "complex-tool" {
		t.Errorf("Linux.PackageNames[apt] = %q, want %q", tool.Linux.PackageNames["apt"], "complex-tool")
	}

	if tool.MacOS == nil {
		t.Fatal("MacOS config is nil")
	}

	if tool.MacOS.Brew != "complex-tool" {
		t.Errorf("MacOS.Brew = %q, want %q", tool.MacOS.Brew, "complex-tool")
	}
}

// Benchmark config parsing
func BenchmarkLoadFromBytes(b *testing.B) {
	content := []byte(ExampleConfig)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadFromBytes(content)
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
	}
}
