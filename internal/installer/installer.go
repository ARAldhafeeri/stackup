// Package installer handles the orchestration of tool installation
package installer

import (
	"fmt"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
	"github.com/araldhafeeri/stackup/internal/executor"
	"github.com/araldhafeeri/stackup/internal/ui"
	"github.com/araldhafeeri/stackup/pkg/version"
)

// Installer orchestrates the installation process
type Installer struct {
	config         *config.Config
	system         *domain.System
	console        *ui.Console
	executor       *executor.Executor
	installedTools map[string]bool
}

// New creates a new Installer instance
func New(cfg *config.Config, sys *domain.System, console *ui.Console) *Installer {
	return &Installer{
		config:         cfg,
		system:         sys,
		console:        console,
		executor:       executor.New(sys),
		installedTools: make(map[string]bool),
	}
}

// Run executes the installation process
func (i *Installer) Run() error {
	i.console.PrintHeader(version.Version, i.system, i.config.Profile)

	// Pre-flight checks
	if err := i.runPreflightChecks(); err != nil {
		return fmt.Errorf("preflight checks failed: %w", err)
	}

	// Resolve dependencies
	toolsToInstall, err := i.resolveDependencies()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	needsReboot := false

	// Install each tool
	for idx, tool := range toolsToInstall {
		i.console.PrintToolHeader(idx+1, len(toolsToInstall), tool)

		if err := i.installTool(tool); err != nil {
			i.console.PrintError(tool.GetDisplayName(), err)
			continue
		}

		i.installedTools[tool.Name] = true

		// Verify if enabled
		if i.config.Settings.VerifyInstallations {
			if err := i.verifyTool(tool); err != nil {
				i.console.PrintWarning(tool.GetDisplayName(), "installed but verification failed")
			} else {
				i.console.PrintSuccess(tool.GetDisplayName(), "installed successfully")
			}
		} else {
			i.console.PrintSuccess(tool.GetDisplayName(), "installed")
		}

		if tool.RequiresReboot {
			needsReboot = true
		}
	}

	i.console.PrintComplete(needsReboot)

	return nil
}

// installTool installs a single tool
func (i *Installer) installTool(tool *config.Tool) error {
	if i.installedTools[tool.Name] {
		i.console.PrintInfo("Already installed, skipping...")
		return nil
	}

	// Pre-install commands
	if err := i.executor.RunCommands(tool.PreInstall, "Pre-install"); err != nil {
		return fmt.Errorf("pre-install failed: %w", err)
	}

	// Custom install commands
	if len(tool.CustomInstall) > 0 {
		if err := i.executor.RunCommands(tool.CustomInstall, "Custom install"); err != nil {
			return fmt.Errorf("custom install failed: %w", err)
		}
		return i.executor.RunCommands(tool.PostInstall, "Post-install")
	}

	// Platform-specific installation
	platformConfig := tool.GetPlatformConfig(i.system.OS)
	if platformConfig == nil {
		return fmt.Errorf("%w: %s on %s", domain.ErrNoPlatformConfig, tool.Name, i.system.OS)
	}

	// Platform custom commands
	if len(platformConfig.CustomCommands) > 0 {
		if err := i.executor.RunCommands(platformConfig.CustomCommands, "Platform install"); err != nil {
			return fmt.Errorf("platform install failed: %w", err)
		}
		return i.executor.RunCommands(tool.PostInstall, "Post-install")
	}

	// Try package manager
	if err := i.executor.InstallViaPackageManager(tool, platformConfig); err == nil {
		return i.executor.RunCommands(tool.PostInstall, "Post-install")
	}

	// Try direct download
	if platformConfig.Installer != "" {
		if err := i.executor.InstallViaDownload(tool, platformConfig); err != nil {
			return err
		}
		return i.executor.RunCommands(tool.PostInstall, "Post-install")
	}

	return domain.ErrNoInstallMethod
}
