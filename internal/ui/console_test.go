package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/domain"
)

func TestNewConsole(t *testing.T) {
	console := NewConsole()
	if console == nil {
		t.Fatal("NewConsole returned nil")
	}

	if console.verbose {
		t.Error("Default verbose should be false")
	}

	if console.noColor {
		t.Error("Default noColor should be false")
	}
}

func TestNewConsoleWithOptions(t *testing.T) {
	console := NewConsoleWithOptions(true, true)
	if console == nil {
		t.Fatal("NewConsoleWithOptions returned nil")
	}

	if !console.verbose {
		t.Error("Expected verbose to be true")
	}

	if !console.noColor {
		t.Error("Expected noColor to be true")
	}
}

func TestPrintHeader(t *testing.T) {
	console := NewConsole()
	sys := &domain.System{
		OS:             "linux",
		Arch:           "amd64",
		PackageManager: "apt",
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	console.PrintHeader("1.0.0", sys, "test-profile")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "StackUp") {
		t.Error("Output should contain 'StackUp'")
	}

	if !strings.Contains(output, "1.0.0") {
		t.Error("Output should contain version")
	}

	if !strings.Contains(output, "linux") {
		t.Error("Output should contain OS")
	}

	if !strings.Contains(output, "test-profile") {
		t.Error("Output should contain profile")
	}
}

func TestPrintToolHeader(t *testing.T) {
	console := NewConsole()
	tool := &config.Tool{
		Name:        "git",
		DisplayName: "Git VCS",
		Description: "Version control system",
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	console.PrintToolHeader(1, 3, tool)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "[1/3]") {
		t.Error("Output should contain progress")
	}

	if !strings.Contains(output, "Git VCS") {
		t.Error("Output should contain display name")
	}

	if !strings.Contains(output, "Version control system") {
		t.Error("Output should contain description")
	}
}

func TestVerboseOutput(t *testing.T) {
	tests := []struct {
		name          string
		verbose       bool
		shouldContain bool
	}{
		{
			name:          "Verbose Enabled",
			verbose:       true,
			shouldContain: true,
		},
		{
			name:          "Verbose Disabled",
			verbose:       false,
			shouldContain: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			console := NewConsoleWithOptions(tt.verbose, false)

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			console.Verbose("test message")

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			contains := strings.Contains(output, "test message")
			if contains != tt.shouldContain {
				t.Errorf("Verbose output contains message = %v, want %v", contains, tt.shouldContain)
			}
		})
	}
}

func TestGetPackageManagerDisplay(t *testing.T) {
	console := NewConsole()

	tests := []struct {
		name     string
		pm       string
		expected string
	}{
		{
			name:     "With Package Manager",
			pm:       "apt",
			expected: "apt",
		},
		{
			name:     "Without Package Manager",
			pm:       "",
			expected: "None",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := console.getPackageManagerDisplay(tt.pm)
			if result != tt.expected {
				t.Errorf("getPackageManagerDisplay(%q) = %q, want %q", tt.pm, result, tt.expected)
			}
		})
	}
}
