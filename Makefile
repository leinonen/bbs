# Go BBS Makefile
BINARY_NAME=gobbs
BUILD_DIR=build
MAIN_FILE=main.go
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags="-s -w"

# Build info
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%FT%T%z)
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Color codes for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build clean run test install uninstall init dev help

# Default target
all: build

## help: Show this help message
help:
	@echo "Go BBS System - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /' | column -t -s ':'

## build: Build the BBS binary
build:
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## clean: Remove build artifacts and temporary files
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f host_key
	@rm -f bbs.db
	@rm -f coverage.out coverage.html
	@rm -f test_bbs_*.db
	@echo "$(GREEN)Clean complete$(NC)"

## run: Build and run the BBS server
run: build
	@echo "$(GREEN)Starting BBS server...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME)

## init: Initialize the database
init: build
	@echo "$(GREEN)Initializing database...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME) -init

## dev: Run in development mode with live reload (requires entr)
dev:
	@command -v entr >/dev/null 2>&1 || { echo "$(RED)entr is required but not installed. Install it with: apt-get install entr$(NC)" >&2; exit 1; }
	@echo "$(GREEN)Starting development mode with auto-reload...$(NC)"
	@echo "$(YELLOW)Watching for file changes...$(NC)"
	find . -name "*.go" | entr -r make run

## test: Run all tests
test:
	@echo "$(GREEN)Running all tests...$(NC)"
	$(GO) test -v ./...

## test-unit: Run unit tests only
test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GO) test -v ./domain ./repository

## test-integration: Run integration tests only
test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GO) test -v ./test

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## test-bench: Run benchmark tests
test-bench:
	@echo "$(GREEN)Running benchmark tests...$(NC)"
	$(GO) test -v -bench=. ./...

## test-race: Run tests with race detection
test-race:
	@echo "$(GREEN)Running tests with race detection...$(NC)"
	$(GO) test -v -race ./...

## deps: Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

## install: Install the binary to /usr/local/bin
install: build
	@echo "$(GREEN)Installing $(BINARY_NAME) to /usr/local/bin...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)Installation complete$(NC)"

## uninstall: Remove the binary from /usr/local/bin
uninstall:
	@echo "$(YELLOW)Uninstalling $(BINARY_NAME)...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)Uninstallation complete$(NC)"

## setup: Complete setup - download deps, build, and initialize database
setup: deps build init
	@echo "$(GREEN)Setup complete! You can now run: make run$(NC)"

## reset: Reset the BBS (clean database and keys)
reset: clean init
	@echo "$(GREEN)BBS has been reset$(NC)"

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t gobbs:$(VERSION) .

## docker-run: Run in Docker container
docker-run: docker-build
	@echo "$(GREEN)Running in Docker...$(NC)"
	docker run -p 2222:2222 gobbs:$(VERSION)

## fmt: Format Go code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GO) fmt ./...

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GO) vet ./...

## lint: Run golangci-lint (requires golangci-lint)
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "$(RED)golangci-lint is required but not installed. See: https://golangci-lint.run/usage/install/$(NC)" >&2; exit 1; }
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run

## check: Run fmt, vet, and test
check: fmt vet test
	@echo "$(GREEN)All checks passed!$(NC)"

## check-all: Run all checks including race detection and coverage
check-all: fmt vet test-coverage test-race
	@echo "$(GREEN)All comprehensive checks passed!$(NC)"

## version: Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

## backup: Backup database
backup:
	@mkdir -p backups
	@cp bbs.db backups/bbs-$(shell date +%Y%m%d-%H%M%S).db 2>/dev/null || echo "$(YELLOW)No database to backup$(NC)"
	@echo "$(GREEN)Database backed up$(NC)"

## restore: Restore database from latest backup
restore:
	@if [ -f backups/*.db ]; then \
		cp $$(ls -t backups/*.db | head -1) bbs.db; \
		echo "$(GREEN)Database restored from latest backup$(NC)"; \
	else \
		echo "$(RED)No backup found$(NC)"; \
	fi

# Create systemd service file
## service: Create systemd service file
service:
	@echo "$(GREEN)Creating systemd service file...$(NC)"
	@echo "[Unit]" > gobbs.service
	@echo "Description=Go BBS SSH Server" >> gobbs.service
	@echo "After=network.target" >> gobbs.service
	@echo "" >> gobbs.service
	@echo "[Service]" >> gobbs.service
	@echo "Type=simple" >> gobbs.service
	@echo "User=nobody" >> gobbs.service
	@echo "WorkingDirectory=/opt/gobbs" >> gobbs.service
	@echo "ExecStart=/usr/local/bin/gobbs" >> gobbs.service
	@echo "Restart=always" >> gobbs.service
	@echo "RestartSec=10" >> gobbs.service
	@echo "" >> gobbs.service
	@echo "[Install]" >> gobbs.service
	@echo "WantedBy=multi-user.target" >> gobbs.service
	@echo "$(GREEN)Service file created: gobbs.service$(NC)"
	@echo "$(YELLOW)To install: sudo cp gobbs.service /etc/systemd/system/$(NC)"
	@echo "$(YELLOW)Then: sudo systemctl daemon-reload && sudo systemctl enable gobbs$(NC)"

.DEFAULT_GOAL := help