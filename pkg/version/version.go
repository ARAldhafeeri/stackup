// Package version provides version information for StackUp
package version

// Version is the current version of StackUp
// This should be updated with each release
const Version = "0.1.0"

// GitCommit is the git commit hash (set during build)
var GitCommit = "unknown"

// BuildDate is the date the binary was built (set during build)
var BuildDate = "unknown"

// GetVersion returns the version string
func GetVersion() string {
	return Version
}

// GetFullVersion returns a detailed version string including git info
func GetFullVersion() string {
	if GitCommit != "unknown" {
		return Version + " (" + GitCommit + ")"
	}
	return Version
}

// GetBuildInfo returns all build information
func GetBuildInfo() string {
	return "Version: " + Version + "\n" +
		"Git Commit: " + GitCommit + "\n" +
		"Build Date: " + BuildDate
}