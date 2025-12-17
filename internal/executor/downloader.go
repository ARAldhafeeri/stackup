package executor

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

// DownloadInstaller handles installation by downloading installers
type DownloadInstaller struct {
	system *domain.System
}

// NewDownloadInstaller creates a new download installer
func NewDownloadInstaller(sys *domain.System) *DownloadInstaller {
	return &DownloadInstaller{system: sys}
}

// Install downloads and executes an installer
func (d *DownloadInstaller) Install(tool *config.Tool, cfg *config.PlatformConfig) error {
	fmt.Printf("   Downloading from %s...\n", cfg.Installer)

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "stackup-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download file
	filePath, err := d.downloadFile(cfg.Installer, tempDir, tool.Name, cfg.Type)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Printf("   Downloaded to %s\n", filePath)

	// Execute installer
	if err := d.executeInstaller(filePath, cfg); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	return nil
}

// downloadFile downloads a file from a URL to the specified directory
func (d *DownloadInstaller) downloadFile(url, destDir, toolName, fileType string) (string, error) {
	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Determine filename
	filename := d.determineFilename(url, toolName, fileType)
	filePath := filepath.Join(destDir, filename)

	// Create destination file
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy content
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("   Downloaded %d bytes\n", written)

	return filePath, nil
}

// determineFilename creates an appropriate filename for the downloaded installer
func (d *DownloadInstaller) determineFilename(url, toolName, fileType string) string {
	// If type is specified, use it
	if fileType != "" {
		return toolName + "." + fileType
	}

	// Try to extract filename from URL
	filename := filepath.Base(url)
	if filename != "" && filename != "." && filename != "/" {
		return filename
	}

	// Default fallback
	return toolName + ".installer"
}

// executeInstaller runs the downloaded installer with appropriate flags
func (d *DownloadInstaller) executeInstaller(path string, cfg *config.PlatformConfig) error {
	fmt.Printf("   Executing installer...\n")

	cmd := d.buildInstallerCommand(path, cfg)
	if cmd == nil {
		return fmt.Errorf("unsupported installer type: %s", cfg.Type)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// buildInstallerCommand creates the appropriate command to execute the installer
func (d *DownloadInstaller) buildInstallerCommand(path string, cfg *config.PlatformConfig) *exec.Cmd {
	switch cfg.Type {
	case "exe", "msi":
		return d.buildWindowsInstallerCommand(path, cfg)
	case "sh", "bash":
		return d.buildShellInstallerCommand(path)
	case "deb":
		return d.buildDebInstallerCommand(path)
	case "rpm":
		return d.buildRpmInstallerCommand(path)
	case "dmg":
		return d.buildDmgInstallerCommand(path)
	case "pkg":
		return d.buildPkgInstallerCommand(path)
	case "appimage":
		return d.buildAppImageCommand(path)
	default:
		// Generic executable - make it executable and run
		return d.buildGenericExecutableCommand(path)
	}
}

// buildWindowsInstallerCommand creates a command for Windows installers (.exe, .msi)
func (d *DownloadInstaller) buildWindowsInstallerCommand(path string, cfg *config.PlatformConfig) *exec.Cmd {
	args := cfg.SilentFlags

	// If no silent flags specified, use common defaults
	if len(args) == 0 {
		if cfg.Type == "msi" {
			// MSI silent install flags
			args = []string{"/i", path, "/quiet", "/norestart"}
			return exec.Command("msiexec", args...)
		}
		// EXE silent install flag
		args = []string{"/S"}
	}

	return exec.Command(path, args...)
}

// buildShellInstallerCommand creates a command for shell script installers
func (d *DownloadInstaller) buildShellInstallerCommand(path string) *exec.Cmd {
	// Make script executable
	os.Chmod(path, 0755)
	return exec.Command("bash", path)
}

// buildDebInstallerCommand creates a command for Debian package installers
func (d *DownloadInstaller) buildDebInstallerCommand(path string) *exec.Cmd {
	return exec.Command("sudo", "dpkg", "-i", path)
}

// buildRpmInstallerCommand creates a command for RPM package installers
func (d *DownloadInstaller) buildRpmInstallerCommand(path string) *exec.Cmd {
	return exec.Command("sudo", "rpm", "-i", path)
}

// buildDmgInstallerCommand creates a command for macOS DMG installers
func (d *DownloadInstaller) buildDmgInstallerCommand(path string) *exec.Cmd {
	// DMG installation is more complex - mount, copy app, unmount
	// This is a simplified version
	return exec.Command("hdiutil", "attach", path)
}

// buildPkgInstallerCommand creates a command for macOS PKG installers
func (d *DownloadInstaller) buildPkgInstallerCommand(path string) *exec.Cmd {
	return exec.Command("sudo", "installer", "-pkg", path, "-target", "/")
}

// buildAppImageCommand creates a command for Linux AppImage installers
func (d *DownloadInstaller) buildAppImageCommand(path string) *exec.Cmd {
	// Make AppImage executable
	os.Chmod(path, 0755)
	// AppImages typically need to be moved to a location in PATH
	// This is a simplified version that just makes it executable
	return exec.Command("chmod", "+x", path)
}

// buildGenericExecutableCommand creates a command for generic executable files
func (d *DownloadInstaller) buildGenericExecutableCommand(path string) *exec.Cmd {
	// Make file executable
	os.Chmod(path, 0755)
	return exec.Command(path)
}

// DownloadToFile is a utility function to download a file without installing
// Useful for scripts or files that need custom handling
func DownloadToFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
