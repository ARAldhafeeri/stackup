// Package main is the entry point for the StackUp CLI application
package main

import (
	"fmt"
	"os"

	"github.com/araldhafeeri/stackup/internal/config"
	"github.com/araldhafeeri/stackup/internal/installer"
	"github.com/araldhafeeri/stackup/internal/platform"
	"github.com/araldhafeeri/stackup/internal/ui"
	"github.com/araldhafeeri/stackup/pkg/version"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("StackUp v%s\n", version.Version)
		os.Exit(0)
	case "example":
		fmt.Print(config.ExampleConfig)
		os.Exit(0)
	case "install":
		if err := runInstall(); err != nil {
			fmt.Fprintf(os.Stderr, "Installation failed: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func runInstall() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("config file required\nUsage: stackup install <config.yaml>")
	}

	configPath := os.Args[2]

	// Load configuration
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Detect system
	sys := platform.Detect()

	// Create console UI
	console := ui.NewConsole()

	// Create and run installer
	inst := installer.New(cfg, sys, console)
	return inst.Run()
}

func printUsage() {
	fmt.Printf("StackUp v%s - Stack your dev tools effortlessly\n\n", version.Version)
	fmt.Println("Usage: stackup <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  install <config.yaml>    Install tools from config file")
	fmt.Println("  version                  Show version")
	fmt.Println("  example                  Show example config")
	fmt.Println("\nFor more information, visit: https://github.com/araldhafeeri/stackup")
}
