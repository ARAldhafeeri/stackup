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
	system      *domain.System
	interactive bool
}

// NewPackageManager creates a new package manager installer
func NewPackageManager(sys *domain.System) *PackageManager {
	return &PackageManager{system: sys, interactive: true}
}

// Install installs a tool using the appropriate package manager
func (pm *PackageManager) Install(tool *config.Tool, cfg *config.PlatformConfig) error {
	if pm.system.PackageManager == "" {
		return fmt.Errorf("no package manager available")
	}

	manager, packageName := pm.getPackageManagerAndName(tool, cfg)

	cmd := pm.buildInstallCommand(packageName, manager)

	if cmd == nil {
		return fmt.Errorf("unsupported package manager: %s", pm.system.PackageManager)
	}

	// Setup stdio - this is crucial for interactive commands
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Important: Inherit environment to ensure proper execution
	cmd.Env = os.Environ()

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

func (pm *PackageManager) getPackageManagerAndName(tool *config.Tool, cfg *config.PlatformConfig) (string, string) {
	// First, check if a specific manager is configured for this tool
	if tool.Manager != "" {
		// If manager is specified, use it with the package names mapping
		if cfg.PackageNames != nil {
			if name, ok := cfg.PackageNames[tool.Manager]; ok {
				return tool.Manager, name
			}
		}
		// If no package name mapping, use tool name
		return tool.Manager, tool.Name
	}

	// If no tool-specific manager, check package names for the system's package manager
	if cfg.PackageNames != nil {
		// First, try the system's detected package manager
		if name, ok := cfg.PackageNames[pm.system.PackageManager]; ok {
			return pm.system.PackageManager, name
		}

		// If system package manager not in package names, use the first available
		for manager, name := range cfg.PackageNames {
			return manager, name
		}
	}

	// Check for brew-specific name (macOS only)
	if cfg.Brew != "" && pm.system.PackageManager == domain.PackageManagerBrew {
		return domain.PackageManagerBrew, cfg.Brew
	}

	// Default to system package manager with tool name
	return pm.system.PackageManager, tool.Name
}

// buildInstallCommand creates the install command for the package manager
func (pm *PackageManager) buildInstallCommand(packageName string, manager string) *exec.Cmd {

	switch manager {
	case domain.PackageManagerAPT:
		if pm.interactive {
			// Interactive mode - let user respond to prompts
			return exec.Command("sudo", "apt-get", "install", packageName)
		}
		return exec.Command("sudo", "apt-get", "install", "-y", packageName)

	case domain.PackageManagerDNF:
		if pm.interactive {
			return exec.Command("sudo", "dnf", "install", packageName)
		}
		return exec.Command("sudo", "dnf", "install", "-y", packageName)

	case domain.PackageManagerPacman:
		if pm.interactive {
			return exec.Command("sudo", "pacman", "-S", packageName)
		}
		return exec.Command("sudo", "pacman", "-S", "--noconfirm", packageName)

	case domain.PackageManagerBrew:
		// Homebrew doesn't typically require interactive confirmation
		return exec.Command("brew", "install", packageName)

	case domain.PackageManagerWinget:
		if pm.interactive {
			return exec.Command("sudo", "winget", "install", "-e", "--id", packageName)
		}
		return exec.Command("sudo", "winget", "install", "-e", "--id", packageName, "--accept-package-agreements", "--accept-source-agreements")

	case domain.PackageManagerChoco:
		if pm.interactive {
			return exec.Command("sudo", "choco", "install", packageName)
		}
		return exec.Command("sudo", "choco", "install", packageName, "-y")

	default:
		return nil
	}
}
