// Package config handles configuration structures and management
package config

// Config represents the main application configuration
type Config struct {
	Profile  string            `yaml:"profile"`
	Settings Settings          `yaml:"settings"`
	Tools    []Tool            `yaml:"tools"`
	Presets  map[string]Preset `yaml:"presets,omitempty"`
}

// Settings contains global installation settings
type Settings struct {
	AutoUpdatePath      bool `yaml:"auto_update_path"`
	VerifyInstallations bool `yaml:"verify_installations"`
}

// Preset defines a named collection of tools
type Preset struct {
	Description string   `yaml:"description"`
	Tools       []string `yaml:"tools"`
}

// Tool represents a single installable tool
type Tool struct {
	Name            string          `yaml:"name"`
	DisplayName     string          `yaml:"display_name,omitempty"`
	Version         string          `yaml:"version"`
	Description     string          `yaml:"description,omitempty"`
	Manager         string          `yaml:"manager,omitempty"`
	Windows         *PlatformConfig `yaml:"windows,omitempty"`
	Linux           *PlatformConfig `yaml:"linux,omitempty"`
	MacOS           *PlatformConfig `yaml:"macos,omitempty"`
	PreInstall      []Command       `yaml:"pre_install,omitempty"`
	CustomInstall   []Command       `yaml:"custom_install,omitempty"`
	PostInstall     []Command       `yaml:"post_install,omitempty"`
	VerifyCommand   string          `yaml:"verify_command,omitempty"`
	RequiresReboot  bool            `yaml:"requires_reboot,omitempty"`
	Dependencies    []string        `yaml:"dependencies,omitempty"`
}

// PlatformConfig contains platform-specific installation details
type PlatformConfig struct {
	Installer      string            `yaml:"installer,omitempty"`
	Type           string            `yaml:"type,omitempty"`
	SilentFlags    []string          `yaml:"silent_flags,omitempty"`
	PackageNames   map[string]string `yaml:"package_names,omitempty"`
	Brew           string            `yaml:"brew,omitempty"`
	CustomCommands []Command         `yaml:"custom_commands,omitempty"`
}

// Command represents a command to execute
type Command struct {
	Command     string   `yaml:"command"`
	Args        []string `yaml:"args,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Sudo        bool     `yaml:"sudo,omitempty"`
	WaitFor     int      `yaml:"wait_for,omitempty"` // seconds to wait after command
	IgnoreError bool     `yaml:"ignore_error,omitempty"`
}

// GetDisplayName returns the display name or falls back to name
func (t *Tool) GetDisplayName() string {
	// use display name if set
	if t.DisplayName != "" {
		return t.DisplayName
	}

	// fallback to tool name
	if t.Name != "" {
		return t.Name 
	}
	
	return "ToolNameNotSet"
}

// GetPlatformConfig returns the appropriate platform configuration
func (t *Tool) GetPlatformConfig(osName string) *PlatformConfig {
	switch osName {
	case "windows":
		return t.Windows
	case "linux":
		return t.Linux
	case "darwin":
		return t.MacOS
	default:
		return nil
	}
}