package domain

import "errors"

var (
	// ErrNoInternet indicates no internet connectivity
	ErrNoInternet = errors.New("no internet connectivity")

	// ErrNoPlatformConfig indicates no configuration for the current platform
	ErrNoPlatformConfig = errors.New("no configuration for current platform")

	// ErrNoInstallMethod indicates no installation method available
	ErrNoInstallMethod = errors.New("no installation method available")

	// ErrDependencyNotFound indicates a dependency was not found
	ErrDependencyNotFound = errors.New("dependency not found")

	// ErrVerificationFailed indicates installation verification failed
	ErrVerificationFailed = errors.New("verification failed")
)