package executor

import (
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

func TestGetPackageName(t *testing.T) {
	tests := []struct {
		name        string
		tool        *config.Tool
		platformCfg *config.PlatformConfig
		pkgManager  string
		expected    string
	}{
		{
			name: "Brew Specific Name",
			tool: &config.Tool{Name: "vscode"},
			platformCfg: &config.PlatformConfig{
				Brew: "visual-studio-code",
			},
			pkgManager: domain.PackageManagerBrew,
			expected:   "visual-studio-code",
		},
		{
			name: "Package Manager Specific Name",
			tool: &config.Tool{Name: "git"},
			platformCfg: &config.PlatformConfig{
				PackageNames: map[string]string{
					"apt": "git-all",
					"dnf": "git-core",
				},
			},
			pkgManager: "apt",
			expected:   "git-all",
		},
		{
			name:        "Default to Tool Name",
			tool:        &config.Tool{Name: "docker"},
			platformCfg: &config.PlatformConfig{},
			pkgManager:  "apt",
			expected:    "docker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := &domain.System{PackageManager: tt.pkgManager}
			pm := NewPackageManager(sys)

			result := pm.getPackageName(tt.tool, tt.platformCfg)
			if result != tt.expected {
				t.Errorf("getPackageName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestBuildInstallCommand(t *testing.T) {
	tests := []struct {
		name        string
		pkgManager  string
		packageName string
		expectNil   bool
	}{
		{
			name:        "APT",
			pkgManager:  domain.PackageManagerAPT,
			packageName: "git",
			expectNil:   false,
		},
		{
			name:        "DNF",
			pkgManager:  domain.PackageManagerDNF,
			packageName: "docker",
			expectNil:   false,
		},
		{
			name:        "Brew",
			pkgManager:  domain.PackageManagerBrew,
			packageName: "node",
			expectNil:   false,
		},
		{
			name:        "Winget",
			pkgManager:  domain.PackageManagerWinget,
			packageName: "Git.Git",
			expectNil:   false,
		},
		{
			name:        "Unsupported",
			pkgManager:  "unsupported",
			packageName: "tool",
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := &domain.System{PackageManager: tt.pkgManager}
			pm := NewPackageManager(sys)

			cmd := pm.buildInstallCommand(tt.packageName)

			if tt.expectNil {
				if cmd != nil {
					t.Error("Expected nil command for unsupported package manager")
				}
			} else {
				if cmd == nil {
					t.Error("Expected non-nil command")
				}
			}
		})
	}
}
