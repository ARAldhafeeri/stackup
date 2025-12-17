package domain

import (
	"testing"
)

func TestPreflightChecks(t *testing.T) {
	tests := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "ErrNoInternet",
			error:    ErrNoInternet,
			expected: true,
		},
		{
			name:     "ErrNoPlatformConfig",
			error:    ErrNoPlatformConfig,
			expected: true,
		},
		{
			name:     "ErrNoInstalledMethod",
			error:    ErrNoInstallMethod,
			expected: true,
		},
		{
			name:     "ErrDependencyNotFound",
			error:    ErrDependencyNotFound,
			expected: true,
		},
		{
			name:     "ErrVerificationFailed",
			error:    ErrVerificationFailed,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPreflightError(tt.error)
			if result != tt.expected {
				t.Errorf("IsPreflightError() = %v, want %v", result, tt.expected)
			}
		})
	}
}
