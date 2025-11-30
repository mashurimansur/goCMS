.PHONY: help build run test test-coverage lint clean install-tools

# Variables
APP_NAME=goCMS
MAIN_PATH=./cmd/server
BINARY_PATH=./tmp/main
GO=go
GOFLAGS=-v

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Display this help screen
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

install-tools: ## Install development tools
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install gotest.tools/gotestsum@v1.12.3
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

build: ## Build the application
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	mkdir -p tmp
	$(GO) build $(GOFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)✓ Build successful: $(BINARY_PATH)$(NC)"

run: build ## Build and run the application
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	$(BINARY_PATH)

test: ## Run unit tests
	@echo "$(GREEN)Running unit tests...$(NC)"
	export GOTOOLCHAIN=go1.25.3+auto && \
	chmod +x .github/scripts/unit_test.sh && \
	./.github/scripts/unit_test.sh
	@echo "$(GREEN)✓ Tests completed$(NC)"

test-coverage: test ## Run tests and display coverage report
	@echo "$(GREEN)Coverage Report:$(NC)"
	@if [ -f code_coverage_results ]; then \
		cat code_coverage_results; \
	fi

lint: ## Run golangci-lint
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run --timeout=5m --out-format=colored-line-number ./...
	@echo "$(GREEN)✓ Linting completed$(NC)"

clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning up...$(NC)"
	rm -rf $(BINARY_PATH)
	rm -f coverage.out coverage.out.temp report.xml code_coverage_results
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

check: lint test ## Run linter and tests

all: clean build test lint ## Clean, build, test and lint

.DEFAULT_GOAL := help
