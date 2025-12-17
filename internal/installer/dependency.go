package installer

import (
	"fmt"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

// resolveDependencies creates an ordered list of tools respecting dependencies
func (i *Installer) resolveDependencies() ([]*config.Tool, error) {
	var result []*config.Tool
	visited := make(map[string]bool)

	var resolve func(*config.Tool) error
	resolve = func(tool *config.Tool) error {
		if visited[tool.Name] {
			return nil
		}

		// Install dependencies first
		for _, depName := range tool.Dependencies {
			dep := i.findTool(depName)
			if dep == nil {
				return fmt.Errorf("%w: '%s' for tool '%s'",
					domain.ErrDependencyNotFound, depName, tool.Name)
			}
			if err := resolve(dep); err != nil {
				return err
			}
		}

		visited[tool.Name] = true
		result = append(result, tool)
		return nil
	}

	for idx := range i.config.Tools {
		if err := resolve(&i.config.Tools[idx]); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// findTool locates a tool by name
func (i *Installer) findTool(name string) *config.Tool {
	for idx := range i.config.Tools {
		if i.config.Tools[idx].Name == name {
			return &i.config.Tools[idx]
		}
	}
	return nil
}
