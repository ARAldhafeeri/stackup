package installer

import (
	"net/http"
	"os"

	"github.com/araldhafeeri/stackup/internal/domain"
)

// runPreflightChecks performs pre-installation system checks
func (i *Installer) runPreflightChecks() error {
	i.console.PrintInfo("Running preflight checks...")

	// Check permissions
	if !i.system.IsWindows() && os.Geteuid() != 0 {
		i.console.PrintWarning("", "Not running as root - some installations may require sudo")
	}

	// Check package manager
	if !i.system.HasPackageManager() {
		i.console.PrintWarning("", "No package manager detected - will use direct downloads")
	}

	// Check internet connectivity
	if err := checkInternet(); err != nil {
		return domain.ErrNoInternet
	}

	i.console.PrintSuccess("", "Preflight checks passed")
	return nil
}

// checkInternet verifies internet connectivity
func checkInternet() error {
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
