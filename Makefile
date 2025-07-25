# go-shared-libs Makefile
.PHONY: help test lint clean tag release verify deps

# Variables
LIBRARY_NAME := go-shared-libs
SHELL := /bin/bash

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

# Default target
.DEFAULT_GOAL := help

# Test all packages
test:
	@printf "%b\n" "$(CYAN)Running tests for all packages...$(RESET)"
	@go test -v ./...

# Test with coverage
test-coverage:
	@printf "%b\n" "$(CYAN)Running tests with coverage...$(RESET)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@printf "%b\n" "$(GREEN)Coverage report generated: coverage.html$(RESET)"

# Test specific package
test-%:
	@printf "%b\n" "$(CYAN)Running tests for $* package...$(RESET)"
	@go test -v ./$*

# Run linter (requires golangci-lint)
lint:
	@printf "%b\n" "$(CYAN)Running linter...$(RESET)"
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		printf "%b\n" "$(YELLOW)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(RESET)"; \
	fi

# Download dependencies
deps:
	@printf "%b\n" "$(CYAN)Downloading dependencies...$(RESET)"
	@go mod download
	@go mod tidy
	@printf "%b\n" "$(GREEN)✓ Dependencies updated$(RESET)"

# Verify all packages build
verify:
	@printf "%b\n" "$(CYAN)Verifying all packages build...$(RESET)"
	@if go build ./...; then \
		printf "%b\n" "$(GREEN)✓ All packages build successfully$(RESET)"; \
	else \
		printf "%b\n" "$(RED)✗ Some packages failed to build$(RESET)"; \
		exit 1; \
	fi

# List all packages
list-packages:
	@printf "%b\n" "$(CYAN)Available packages:$(RESET)"
	@find . -name "*.go" -not -path "./.*" | xargs dirname | sort -u | sed 's|^./||' | grep -v '^$$' | sed 's/^/  - /'

# Create and push a new tag
tag:
	@if [ -z "$(VERSION)" ]; then \
		printf "%b\n" "$(YELLOW)Usage: make tag VERSION=v1.0.0$(RESET)"; \
		exit 1; \
	fi
	@printf "%b\n" "$(CYAN)Creating tag $(VERSION)...$(RESET)"
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@printf "%b\n" "$(GREEN)✓ Tag $(VERSION) created and pushed$(RESET)"

# Show current version (latest tag)
version:
	@printf "%b\n" "$(CYAN)Current version:$(RESET)"
	@git describe --tags --abbrev=0 2>/dev/null || echo "No tags found"

# Release preparation
release-check:
	@printf "%b\n" "$(CYAN)Checking release readiness...$(RESET)"
	@printf "%b\n" "$(CYAN)1. Running tests...$(RESET)"
	@make test
	@printf "%b\n" "$(CYAN)2. Running linter...$(RESET)"
	@make lint
	@printf "%b\n" "$(CYAN)3. Verifying build...$(RESET)"
	@make verify
	@printf "%b\n" "$(GREEN)✓ Ready for release$(RESET)"

# Clean build artifacts
clean:
	@printf "%b\n" "$(CYAN)Cleaning up...$(RESET)"
	@rm -f coverage.out coverage.html
	@go clean -cache
	@printf "%b\n" "$(GREEN)Clean complete$(RESET)"

# Documentation generation (requires godoc)
docs:
	@printf "%b\n" "$(CYAN)Starting documentation server...$(RESET)"
	@if command -v godoc > /dev/null 2>&1; then \
		printf "%b\n" "$(GREEN)Documentation available at: http://localhost:6060/pkg/$(RESET)"; \
		godoc -http=:6060; \
	else \
		printf "%b\n" "$(YELLOW)godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest$(RESET)"; \
	fi

# Help
help:
	@printf "%b\n" "$(CYAN)$(LIBRARY_NAME) - Available Commands:$(RESET)"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run all tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make test-<package>  - Run tests for specific package"
	@echo "  make lint            - Run linter (requires golangci-lint)"
	@echo "  make verify          - Verify all packages build"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps            - Download and tidy dependencies"
	@echo ""
	@echo "Release:"
	@echo "  make tag VERSION=v1.0.0  - Create and push new tag"
	@echo "  make version             - Show current version"
	@echo "  make release-check       - Check if ready for release"
	@echo ""
	@echo "Utilities:"
	@echo "  make list-packages   - List all available packages"
	@echo "  make docs            - Start documentation server"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make help            - Show this help"
