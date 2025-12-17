// Package platform handles system detection and platform-specific operations
package platform

import (
	"os/exec"
	"runtime"

	"github.com/araldhafeeri/stackup/internal/domain"
)

// Detect identifies the current system and available package manager
func Detect() *domain.System {
	sys := &domain.System{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	sys.PackageManager = detectPackageManager(sys.OS)

	return sys
}

// detectPackageManager finds the available package manager
func detectPackageManager(osName string) string {
	switch osName {
	case "linux":
		if commandExists("apt-get") {
			return domain.PackageManagerAPT
		} else if commandExists("dnf") {
			return domain.PackageManagerDNF
		} else if commandExists("pacman") {
			return domain.PackageManagerPacman
		}
	case "darwin":
		if commandExists("brew") {
			return domain.PackageManagerBrew
		}
	case "windows":
		if commandExists("winget") {
			return domain.PackageManagerWinget
		} else if commandExists("choco") {
			return domain.PackageManagerChoco
		}
	}

	return ""
}

// commandExists checks if a command is available in PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// IsElevated checks if the process is running with elevated privileges
func IsElevated(sys *domain.System) bool {
	// This is a simplified version - actual implementation would vary by OS
	return true
}
