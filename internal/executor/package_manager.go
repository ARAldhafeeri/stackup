// PackageManager handles installation via system package managers
package executor

import (
	"fmt"
	"strings"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

type PackageManager struct {
	system *domain.System
}

func NewPackageManager(sys *domain.System) *PackageManager {
	return &PackageManager{
		system: sys,
	}
}

func (pm *PackageManager) Install(tool *config.Tool, cfg *config.PlatformConfig) error {
	// Use custom package manager commands if specified
	if len(cfg.PackageManagerCommands) > 0 {
		commands := pm.makeNonInteractive(cfg.PackageManagerCommands)
		runner := NewCommandRunner(pm.system)
		return runner.Run(commands, "Package manager install")
	}

	// Try to determine and use system package manager
	switch pm.system.OS {
	case "windows":
		return pm.installWindows(tool, cfg)
	case "darwin":
		return pm.installDarwin(tool, cfg)
	case "linux":
		return pm.installLinux(tool, cfg)
	default:
		return fmt.Errorf("unsupported OS: %s", pm.system.OS)
	}
}

func (pm *PackageManager) makeNonInteractive(commands []config.Command) []config.Command {
	nonInteractiveCmds := make([]config.Command, len(commands))

	for i, cmd := range commands {
		commandStr := cmd.Command

		// Add non-interactive flags based on package manager
		if strings.Contains(commandStr, "winget install") {
			commandStr = commandStr + " --accept-package-agreements --accept-source-agreements"
		} else if strings.Contains(commandStr, "apt install") ||
			strings.Contains(commandStr, "apt-get install") {
			if !strings.Contains(commandStr, " -y") && !strings.Contains(commandStr, " --yes") {
				commandStr = commandStr + " -y"
			}
		} else if strings.Contains(commandStr, "yum install") ||
			strings.Contains(commandStr, "dnf install") {
			if !strings.Contains(commandStr, " -y") {
				commandStr = commandStr + " -y"
			}
		} else if strings.Contains(commandStr, "pacman -S") {
			if !strings.Contains(commandStr, " --noconfirm") {
				commandStr = commandStr + " --noconfirm"
			}
		} else if strings.Contains(commandStr, "choco install") {
			if !strings.Contains(commandStr, " -y") {
				commandStr = commandStr + " -y"
			}
		} else if strings.Contains(commandStr, "brew install") {
			// Homebrew is usually non-interactive
		}

		nonInteractiveCmds[i] = config.Command{
			Command: commandStr,
			WorkDir: cmd.WorkDir,
			Env:     cmd.Env,
		}
	}

	return nonInteractiveCmds
}

func (pm *PackageManager) installWindows(tool *config.Tool, cfg *config.PlatformConfig) error {
	// Try winget first (Windows 10+)
	wingetCmd := fmt.Sprintf("winget install --id %s --exact --accept-package-agreements --accept-source-agreements",
		cfg.PackageName)

	commands := []config.Command{
		{Command: wingetCmd},
	}

	runner := NewCommandRunner(pm.system)
	return runner.Run(commands, "Winget install")
}

// ... rest of your existing package manager methods ...
