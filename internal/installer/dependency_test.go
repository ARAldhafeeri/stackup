package installer

import (
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
	"github.com/araldhafeeri/stackup/internal/ui"
)

func TestResolveDependencies(t *testing.T) {
	tests := []struct {
		name        string
		tools       []config.Tool
		expectError bool
		validate    func(*testing.T, []*config.Tool)
	}{
		{
			name: "No Dependencies",
			tools: []config.Tool{
				{Name: "tool-a", Version: "1.0"},
				{Name: "tool-b", Version: "1.0"},
			},
			expectError: false,
			validate: func(t *testing.T, result []*config.Tool) {
				if len(result) != 2 {
					t.Errorf("len(result) = %d, want 2", len(result))
				}
			},
		},
		{
			name: "Simple Dependency Chain",
			tools: []config.Tool{
				{Name: "tool-c", Version: "1.0", Dependencies: []string{"tool-a"}},
				{Name: "tool-a", Version: "1.0"},
				{Name: "tool-b", Version: "1.0", Dependencies: []string{"tool-a"}},
			},
			expectError: false,
			validate: func(t *testing.T, result []*config.Tool) {
				if len(result) != 3 {
					t.Errorf("len(result) = %d, want 3", len(result))
				}
				// tool-a should be first
				if result[0].Name != "tool-a" {
					t.Errorf("First tool = %q, want %q", result[0].Name, "tool-a")
				}
				// Verify dependencies come before dependents
				toolIndex := make(map[string]int)
				for i, tool := range result {
					toolIndex[tool.Name] = i
				}
				for i, tool := range result {
					for _, dep := range tool.Dependencies {
						if toolIndex[dep] >= i {
							t.Errorf("Tool %s at index %d depends on %s at index %d (dependency should come first)",
								tool.Name, i, dep, toolIndex[dep])
						}
					}
				}
			},
		},
		{
			name: "Multiple Dependencies",
			tools: []config.Tool{
				{Name: "tool-a", Version: "1.0"},
				{Name: "tool-b", Version: "1.0"},
				{Name: "tool-c", Version: "1.0", Dependencies: []string{"tool-a", "tool-b"}},
			},
			expectError: false,
			validate: func(t *testing.T, result []*config.Tool) {
				if len(result) != 3 {
					t.Errorf("len(result) = %d, want 3", len(result))
				}
				// tool-c should be last
				if result[len(result)-1].Name != "tool-c" {
					t.Errorf("Last tool = %q, want %q", result[len(result)-1].Name, "tool-c")
				}
			},
		},
		{
			name: "Missing Dependency",
			tools: []config.Tool{
				{Name: "tool-a", Version: "1.0", Dependencies: []string{"nonexistent"}},
			},
			expectError: true,
		},
		{
			name: "Deep Dependency Chain",
			tools: []config.Tool{
				{Name: "tool-d", Version: "1.0", Dependencies: []string{"tool-c"}},
				{Name: "tool-c", Version: "1.0", Dependencies: []string{"tool-b"}},
				{Name: "tool-b", Version: "1.0", Dependencies: []string{"tool-a"}},
				{Name: "tool-a", Version: "1.0"},
			},
			expectError: false,
			validate: func(t *testing.T, result []*config.Tool) {
				if len(result) != 4 {
					t.Errorf("len(result) = %d, want 4", len(result))
				}
				// Check order: a, b, c, d
				expectedOrder := []string{"tool-a", "tool-b", "tool-c", "tool-d"}
				for i, expected := range expectedOrder {
					if result[i].Name != expected {
						t.Errorf("result[%d].Name = %q, want %q", i, result[i].Name, expected)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{Tools: tt.tools}
			sys := &domain.System{OS: "linux", PackageManager: "apt"}
			console := ui.NewConsole()
			installer := New(cfg, sys, console)

			result, err := installer.resolveDependencies()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestFindTool(t *testing.T) {
	cfg := &config.Config{
		Tools: []config.Tool{
			{Name: "git", Version: "latest"},
			{Name: "docker", Version: "latest"},
			{Name: "node", Version: "20.x"},
		},
	}

	sys := &domain.System{OS: "linux"}
	console := ui.NewConsole()
	installer := New(cfg, sys, console)

	tests := []struct {
		name     string
		toolName string
		found    bool
	}{
		{"Existing Tool - git", "git", true},
		{"Existing Tool - docker", "docker", true},
		{"Existing Tool - node", "node", true},
		{"Non-existent Tool", "nonexistent", false},
		{"Empty Name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := installer.findTool(tt.toolName)
			found := tool != nil

			if found != tt.found {
				t.Errorf("findTool(%q) found = %v, want %v", tt.toolName, found, tt.found)
			}

			if found && tool.Name != tt.toolName {
				t.Errorf("Found tool name = %q, want %q", tool.Name, tt.toolName)
			}
		})
	}
}

// Benchmark dependency resolution
func BenchmarkResolveDependencies(b *testing.B) {
	cfg := &config.Config{
		Tools: []config.Tool{
			{Name: "tool-a", Version: "1.0"},
			{Name: "tool-b", Version: "1.0", Dependencies: []string{"tool-a"}},
			{Name: "tool-c", Version: "1.0", Dependencies: []string{"tool-b"}},
			{Name: "tool-d", Version: "1.0", Dependencies: []string{"tool-a", "tool-c"}},
			{Name: "tool-e", Version: "1.0", Dependencies: []string{"tool-d"}},
		},
	}

	sys := &domain.System{OS: "linux"}
	console := ui.NewConsole()
	installer := New(cfg, sys, console)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = installer.resolveDependencies()
	}
}
