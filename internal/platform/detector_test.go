package platform

import (
	"runtime"
	"testing"

	"github.com/araldhafeeri/stackup/internal/domain"
)

func TestDetect(t *testing.T) {
	sys := Detect()

	if sys == nil {
		t.Fatal("Detect() returned nil")
	}

	if sys.OS == "" {
		t.Error("OS should not be empty")
	}

	if sys.Arch == "" {
		t.Error("Arch should not be empty")
	}

	expectedOS := runtime.GOOS
	if sys.OS != expectedOS {
		t.Errorf("OS = %q, want %q", sys.OS, expectedOS)
	}

	expectedArch := runtime.GOARCH
	if sys.Arch != expectedArch {
		t.Errorf("Arch = %q, want %q", sys.Arch, expectedArch)
	}

	t.Logf("Detected: OS=%s, Arch=%s, PackageManager=%s",
		sys.OS, sys.Arch, sys.PackageManager)
}

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name    string
		command string
		// We can't predict exact results across systems, so we just test it doesn't panic
	}{
		{"Go Command", "go"},
		{"NonExistent Command", "this-command-does-not-exist-12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := commandExists(tt.command)
			t.Logf("%s exists: %v", tt.command, result)
		})
	}
}

func TestSystemMethods(t *testing.T) {
	tests := []struct {
		name      string
		system    *domain.System
		isLinux   bool
		isWindows bool
		isMacOS   bool
	}{
		{
			name:      "Linux System",
			system:    &domain.System{OS: "linux"},
			isLinux:   true,
			isWindows: false,
			isMacOS:   false,
		},
		{
			name:      "Windows System",
			system:    &domain.System{OS: "windows"},
			isLinux:   false,
			isWindows: true,
			isMacOS:   false,
		},
		{
			name:      "macOS System",
			system:    &domain.System{OS: "darwin"},
			isLinux:   false,
			isWindows: false,
			isMacOS:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.system.IsLinux() != tt.isLinux {
				t.Errorf("IsLinux() = %v, want %v", tt.system.IsLinux(), tt.isLinux)
			}
			if tt.system.IsWindows() != tt.isWindows {
				t.Errorf("IsWindows() = %v, want %v", tt.system.IsWindows(), tt.isWindows)
			}
			if tt.system.IsMacOS() != tt.isMacOS {
				t.Errorf("IsMacOS() = %v, want %v", tt.system.IsMacOS(), tt.isMacOS)
			}
		})
	}
}

func TestHasPackageManager(t *testing.T) {
	tests := []struct {
		name     string
		system   *domain.System
		expected bool
	}{
		{
			name:     "With Package Manager",
			system:   &domain.System{PackageManager: "apt"},
			expected: true,
		},
		{
			name:     "Without Package Manager",
			system:   &domain.System{PackageManager: ""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.system.HasPackageManager()
			if result != tt.expected {
				t.Errorf("HasPackageManager() = %v, want %v", result, tt.expected)
			}
		})
	}
}
