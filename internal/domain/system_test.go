package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemChecks(t *testing.T) {
	tests := []struct {
		name        string
		system      *System
		expectError bool
	}{
		{
			name: "Linux with Package Manager",
			system: &System{
				OS:             "linux",
				Arch:           "amd64",
				PackageManager: "apt",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.system.OS, func(t *testing.T) {
			result := tt.system.IsLinux()
			assert.Equal(t, tt.expectError, result == false)
			assert.Equal(t, tt.system.OS == "linux", result)
			assert.Equal(t, tt.system.Arch == "amd64", true)
			assert.Equal(t, tt.system.PackageManager == "apt", true)
		})
	}
}
