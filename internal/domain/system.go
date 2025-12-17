// Package domain contains core domain models
package domain

// System represents the detected system information
type System struct {
	OS             string
	Arch           string
	PackageManager string
}

// PackageManager types
const (
	PackageManagerAPT    = "apt"
	PackageManagerDNF    = "dnf"
	PackageManagerPacman = "pacman"
	PackageManagerBrew   = "brew"
	PackageManagerWinget = "winget"
	PackageManagerChoco  = "choco"
)

// IsLinux returns true if the system is Linux
func (s *System) IsLinux() bool {
	return s.OS == "linux"
}

// IsWindows returns true if the system is Windows
func (s *System) IsWindows() bool {
	return s.OS == "windows"
}

// IsMacOS returns true if the system is macOS
func (s *System) IsMacOS() bool {
	return s.OS == "darwin"
}

// HasPackageManager returns true if a package manager is available
func (s *System) HasPackageManager() bool {
	return s.PackageManager != ""
}