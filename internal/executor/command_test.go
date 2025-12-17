package executor

import (
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

func TestCommandRunner(t *testing.T) {
	sys := &domain.System{OS: "linux", Arch: "amd64"}
	runner := NewCommandRunner(sys)

	if runner == nil {
		t.Fatal("NewCommandRunner returned nil")
	}
}

func TestBuildCommand(t *testing.T) {
	tests := []struct {
		name       string
		system     *domain.System
		cmdDef     config.Command
		expectSudo bool
	}{
		{
			name:   "Simple Command Linux",
			system: &domain.System{OS: "linux"},
			cmdDef: config.Command{
				Command: "echo",
				Args:    []string{"hello"},
				Sudo:    false,
			},
			expectSudo: false,
		},
		{
			name:   "Sudo Command Linux",
			system: &domain.System{OS: "linux"},
			cmdDef: config.Command{
				Command: "apt-get",
				Args:    []string{"update"},
				Sudo:    true,
			},
			expectSudo: true,
		},
		{
			name:   "Sudo Command Windows",
			system: &domain.System{OS: "windows"},
			cmdDef: config.Command{
				Command: "echo",
				Args:    []string{"hello"},
				Sudo:    true, // Should be ignored on Windows
			},
			expectSudo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewCommandRunner(tt.system)
			cmd := runner.buildCommand(tt.cmdDef)

			if cmd == nil {
				t.Fatal("buildCommand returned nil")
			}

			// Check if sudo is in the command path
			hasSudo := cmd.Path == "/usr/bin/sudo" || cmd.Path == "/bin/sudo" ||
				(len(cmd.Args) > 0 && cmd.Args[0] == "sudo")

			if hasSudo != tt.expectSudo {
				t.Errorf("Command has sudo = %v, want %v", hasSudo, tt.expectSudo)
			}
		})
	}
}

func TestRunCommandsEmpty(t *testing.T) {
	sys := &domain.System{OS: "linux"}
	runner := NewCommandRunner(sys)

	err := runner.Run([]config.Command{}, "test")
	if err != nil {
		t.Errorf("Run with empty commands should not error: %v", err)
	}
}

func TestRunCommandsIgnoreError(t *testing.T) {
	sys := &domain.System{OS: "linux"}
	runner := NewCommandRunner(sys)

	commands := []config.Command{
		{
			Command:     "false", // Command that always fails
			IgnoreError: true,
		},
	}

	err := runner.Run(commands, "test")
	if err != nil {
		t.Errorf("Run should not error with ignore_error=true: %v", err)
	}
}
