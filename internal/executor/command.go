package executor

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

// CommandRunner handles execution of shell commands
type CommandRunner struct {
	system *domain.System
}

// NewCommandRunner creates a new command runner
func NewCommandRunner(sys *domain.System) *CommandRunner {
	return &CommandRunner{system: sys}
}

// Run executes a list of commands
func (r *CommandRunner) Run(commands []config.Command, stage string) error {
	if len(commands) == 0 {
		return nil
	}

	fmt.Printf("   Running %s commands...\n", stage)

	for _, cmdDef := range commands {
		if cmdDef.Description != "" {
			fmt.Printf("   → %s\n", cmdDef.Description)
		}

		cmd := r.buildCommand(cmdDef)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			if !cmdDef.IgnoreError {
				return fmt.Errorf("command '%s' failed: %w", cmdDef.Command, err)
			}
			fmt.Printf("   ⚠️  Command failed but continuing (ignore_error=true)\n")
		}

		if cmdDef.WaitFor > 0 {
			fmt.Printf("   Waiting %d seconds...\n", cmdDef.WaitFor)
			time.Sleep(time.Duration(cmdDef.WaitFor) * time.Second)
		}
	}

	return nil
}

// buildCommand creates an exec.Cmd from a Command definition
func (r *CommandRunner) buildCommand(cmdDef config.Command) *exec.Cmd {
	if cmdDef.Sudo && !r.system.IsWindows() {
		args := append([]string{cmdDef.Command}, cmdDef.Args...)
		return exec.Command("sudo", args...)
	}

	return exec.Command(cmdDef.Command, cmdDef.Args...)
}
