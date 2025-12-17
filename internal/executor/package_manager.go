package executor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

// PackageManager handles installation via system package managers
type PackageManager struct {
	system *domain.System
}

// NewPackageManager creates a new package manager installer
func NewPackageManager(sys *domain.System) *PackageManager {
	return &PackageManager{system: sys}
}

// Install installs a tool using the appropriate package manager
func (pm *PackageManager) Install(tool *config.Tool, cfg *config.PlatformConfig) error {
	if pm.system.PackageManager == "" {
		return fmt.Errorf("no package manager available")
	}

	packageName := pm.getPackageName(tool, cfg)
	cmd := pm.buildInstallCommand(packageName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// getPackageName determines the correct package name for the current package manager
func (pm *PackageManager) getPackageName(tool *config.Tool, cfg *config.PlatformConfig) string {
	// Check for brew-specific name
	if cfg.Brew != "" && pm.system.PackageManager == domain.PackageManagerBrew {
		return cfg.Brew
	}

	// Check for package manager specific names
	if cfg.PackageNames != nil {
		if name, ok := cfg.PackageNames[pm.system.PackageManager]; ok {
			return name
		}
	}

	// Default to tool name
	return tool.Name
}

// buildInstallCommand creates the install command for the package manager
func (pm *PackageManager) buildInstallCommand(packageName string) *exec.Cmd {
	switch pm.system.PackageManager {
	case domain.PackageManagerAPT:
		return exec.Command("sudo", "apt-get", "install", "-y", packageName)
	case domain.PackageManagerDNF:
		return exec.Command("sudo", "dnf", "install", "-y", packageName)
	case domain.PackageManagerPacman:
		return exec.Command("sudo", "pacman", "-S", "--noconfirm", packageName)
	case domain.PackageManagerBrew:
		return exec.Command("brew", "install", packageName)
	case domain.PackageManagerWinget:
		return exec.Command("winget", "install", "-e", "--id", packageName)
	case domain.PackageManagerChoco:
		return exec.Command("choco", "install", packageName, "-y")
	default:
		return nil
	}
}
