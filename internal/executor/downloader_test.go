package executor

import (
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

func TestDetermineFilename(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		toolName string
		fileType string
		expected string
	}{
		{
			name:     "With File Type",
			url:      "https://example.com/installer",
			toolName: "docker",
			fileType: "exe",
			expected: "docker.exe",
		},
		{
			name:     "From URL",
			url:      "https://example.com/downloads/installer.msi",
			toolName: "tool",
			fileType: "",
			expected: "installer.msi",
		},
		{
			name:     "Fallback to Default",
			url:      "https://example.com/mytool.installer",
			toolName: "mytool",
			fileType: "",
			expected: "mytool.installer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := &domain.System{OS: "linux"}
			d := NewDownloadInstaller(sys)

			result := d.determineFilename(tt.url, tt.toolName, tt.fileType)
			if result != tt.expected {
				t.Errorf("determineFilename() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestBuildInstallerCommand(t *testing.T) {
	tests := []struct {
		name      string
		fileType  string
		expectNil bool
	}{
		{name: "EXE", fileType: "exe", expectNil: false},
		{name: "MSI", fileType: "msi", expectNil: false},
		{name: "DEB", fileType: "deb", expectNil: false},
		{name: "RPM", fileType: "rpm", expectNil: false},
		{name: "SH", fileType: "sh", expectNil: false},
		{name: "DMG", fileType: "dmg", expectNil: false},
		{name: "PKG", fileType: "pkg", expectNil: false},
		{name: "AppImage", fileType: "appimage", expectNil: false},
		{name: "Generic", fileType: "", expectNil: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := &domain.System{OS: "linux"}
			d := NewDownloadInstaller(sys)

			cfg := &config.PlatformConfig{Type: tt.fileType}
			cmd := d.buildInstallerCommand("/tmp/test", cfg)

			if tt.expectNil {
				if cmd != nil {
					t.Error("Expected nil command")
				}
			} else {
				if cmd == nil {
					t.Error("Expected non-nil command")
				}
			}
		})
	}
}

func TestBuildWindowsInstallerCommand(t *testing.T) {
	sys := &domain.System{OS: "windows"}
	d := NewDownloadInstaller(sys)

	tests := []struct {
		name        string
		fileType    string
		silentFlags []string
		checkArgs   func(*testing.T, []string)
	}{
		{
			name:        "EXE with custom flags",
			fileType:    "exe",
			silentFlags: []string{"/S", "/quiet"},
			checkArgs: func(t *testing.T, args []string) {
				if len(args) < 2 {
					t.Error("Expected at least 2 arguments")
				}
			},
		},
		{
			name:        "EXE with default flags",
			fileType:    "exe",
			silentFlags: nil,
			checkArgs: func(t *testing.T, args []string) {
				if len(args) == 0 {
					t.Error("Expected default silent flag")
				}
			},
		},
		{
			name:        "MSI",
			fileType:    "msi",
			silentFlags: nil,
			checkArgs: func(t *testing.T, args []string) {
				// MSI uses msiexec, different structure
				if len(args) < 3 {
					t.Error("Expected msiexec arguments")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.PlatformConfig{
				Type:        tt.fileType,
				SilentFlags: tt.silentFlags,
			}

			cmd := d.buildWindowsInstallerCommand("/tmp/test."+tt.fileType, cfg)
			if cmd == nil {
				t.Fatal("Expected non-nil command")
			}

			if tt.checkArgs != nil {
				tt.checkArgs(t, cmd.Args)
			}
		})
	}
}
