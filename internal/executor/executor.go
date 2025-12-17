// Package executor handles command execution and installation operations
package executor

import (
	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

// Executor handles execution of commands and installations
type Executor struct {
	system            *domain.System
	commandRunner     *CommandRunner
	packageManager    *PackageManager
	downloadInstaller *DownloadInstaller
}

// New creates a new Executor
func New(sys *domain.System) *Executor {
	return &Executor{
		system:            sys,
		commandRunner:     NewCommandRunner(sys),
		packageManager:    NewPackageManager(sys),
		downloadInstaller: NewDownloadInstaller(sys),
	}
}

// RunCommands executes a list of commands
func (e *Executor) RunCommands(commands []config.Command, stage string) error {
	return e.commandRunner.Run(commands, stage)
}

// InstallViaPackageManager installs a tool using the system package manager
func (e *Executor) InstallViaPackageManager(tool *config.Tool, cfg *config.PlatformConfig) error {
	return e.packageManager.Install(tool, cfg)
}

// InstallViaDownload installs a tool by downloading an installer
func (e *Executor) InstallViaDownload(tool *config.Tool, cfg *config.PlatformConfig) error {
	return e.downloadInstaller.Install(tool, cfg)
}
