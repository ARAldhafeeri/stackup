package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return LoadFromBytes(data)
}

// LoadFromBytes parses configuration from byte slice
func LoadFromBytes(data []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Set default display names
	for i := range config.Tools {
		if config.Tools[i].DisplayName == "" {
			config.Tools[i].DisplayName = config.Tools[i].Name
		}
	}

	return &config, nil
}