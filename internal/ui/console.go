// Package ui handles user interface output and formatting
package ui

import (
	"fmt"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

// Console handles formatted console output
type Console struct {
	// Could add options here like verbose mode, color support, etc.
	verbose bool
	noColor bool
}

// NewConsole creates a new console UI
func NewConsole() *Console {
	return &Console{
		verbose: false,
		noColor: false,
	}
}

// NewConsoleWithOptions creates a console with specific options
func NewConsoleWithOptions(verbose, noColor bool) *Console {
	return &Console{
		verbose: verbose,
		noColor: noColor,
	}
}

// PrintHeader prints the application header with system info
func (c *Console) PrintHeader(version string, sys *domain.System, profile string) {
	fmt.Println("🚀 StackUp v" + version)
	fmt.Println("============================")
	fmt.Printf("OS: %s | Arch: %s | Package Manager: %s\n",
		sys.OS, sys.Arch, c.getPackageManagerDisplay(sys.PackageManager))

	if profile != "" {
		fmt.Printf("Profile: %s\n", profile)
	}
	fmt.Println()
}

// PrintToolHeader prints the header for a tool installation
func (c *Console) PrintToolHeader(current, total int, tool *config.Tool) {
	fmt.Printf("[%d/%d] Installing %s...\n", current, total, tool.GetDisplayName())
	if tool.Description != "" {
		fmt.Printf("    %s\n", tool.Description)
	}
}

// PrintSuccess prints a success message
func (c *Console) PrintSuccess(name, message string) {
	if name != "" {
		fmt.Printf("✅ %s %s\n\n", name, message)
	} else {
		fmt.Printf("✅ %s\n\n", message)
	}
}

// PrintError prints an error message
func (c *Console) PrintError(name string, err error) {
	fmt.Printf("❌ Failed to install %s: %v\n\n", name, err)
}

// PrintWarning prints a warning message
func (c *Console) PrintWarning(name, message string) {
	if name != "" {
		fmt.Printf("⚠️  %s: %s\n\n", name, message)
	} else {
		fmt.Printf("⚠️  %s\n", message)
	}
}

// PrintInfo prints an informational message
func (c *Console) PrintInfo(message string) {
	fmt.Printf("   %s\n", message)
}

// PrintComplete prints the completion message
func (c *Console) PrintComplete(needsReboot bool) {
	fmt.Println("🎉 Installation complete!")

	if needsReboot {
		fmt.Println()
		fmt.Println("⚠️  Some tools require a system reboot to complete installation.")
		fmt.Println("   Please restart your computer when convenient.")
	}
}

// PrintSeparator prints a visual separator
func (c *Console) PrintSeparator() {
	fmt.Println("----------------------------")
}

// getPackageManagerDisplay returns a user-friendly name for the package manager
func (c *Console) getPackageManagerDisplay(pm string) string {
	if pm == "" {
		return "None"
	}
	return pm
}

// Verbose prints a message only if verbose mode is enabled
func (c *Console) Verbose(format string, args ...interface{}) {
	if c.verbose {
		fmt.Printf("   [VERBOSE] "+format+"\n", args...)
	}
}

// Debug prints debug information
func (c *Console) Debug(format string, args ...interface{}) {
	if c.verbose {
		fmt.Printf("   [DEBUG] "+format+"\n", args...)
	}
}
