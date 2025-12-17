package installer

import (
	"os/exec"
	"strings"

	"github.com/araldhafeeri/stackup/internal/config"
)

// verifyTool checks if a tool was installed correctly
func (i *Installer) verifyTool(tool *config.Tool) error {
	// Use custom verify command if specified
	if tool.VerifyCommand != "" {
		parts := strings.Fields(tool.VerifyCommand)
		cmd := exec.Command(parts[0], parts[1:]...)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}

	// Default: try --version
	cmd := exec.Command(tool.Name, "--version")
	return cmd.Run()
}
