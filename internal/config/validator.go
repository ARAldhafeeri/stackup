package config

import (
	"fmt"
)

// Validate checks if the configuration is valid
func Validate(cfg *Config) error {
	if len(cfg.Tools) == 0 {
		return fmt.Errorf("no tools defined in configuration")
	}

	toolNames := make(map[string]bool)

	for i, tool := range cfg.Tools {
		// Check for required fields
		if tool.Name == "" {
			return fmt.Errorf("tool at index %d missing name", i)
		}

		// Check for duplicate tool names
		if toolNames[tool.Name] {
			return fmt.Errorf("duplicate tool name: %s", tool.Name)
		}
		toolNames[tool.Name] = true

		// Validate dependencies exist
		for _, dep := range tool.Dependencies {
			if !toolNames[dep] && !hasTool(cfg.Tools, dep) {
				return fmt.Errorf("tool %s has unknown dependency: %s", tool.Name, dep)
			}
		}

		// Check that at least one platform is configured
		if tool.Windows == nil && tool.Linux == nil && tool.MacOS == nil && len(tool.CustomInstall) == 0 {
			return fmt.Errorf("tool %s has no platform configuration", tool.Name)
		}
	}

	return nil
}

func hasTool(tools []Tool, name string) bool {
	for _, tool := range tools {
		if tool.Name == name {
			return true
		}
	}
	return false
}