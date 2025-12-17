#!/bin/bash

# Test coverage script for StackUp

set -e

echo "Running tests with coverage..."

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Display coverage summary
go tool cover -func=coverage.out | grep total

echo ""
echo "Coverage report generated: coverage.html"
echo "Open coverage.html in your browser to view detailed coverage"
