.PHONY: build test clean run lint

# Variables
BINARY_NAME=godora
GO=go
GOFMT=gofmt
GOLINT=$(shell go env GOPATH)/bin/golangci-lint

# Build
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/api

# Test
test:
	$(GO) test -v -cover ./...

# Clean
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

# Run
run:
	$(GO) run ./cmd/api

# Lint
lint: $(GOLINT)
	$(GOLINT) run

# Format
fmt:
	$(GOFMT) -w .

# Install dependencies
deps: $(GOLINT)
	$(GO) mod download

# Install golangci-lint
$(GOLINT):
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# All-in-one development setup
dev: deps fmt lint test build

# Help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make run      - Run the application"
	@echo "  make lint     - Run linter"
	@echo "  make fmt      - Format code"
	@echo "  make deps     - Install dependencies"
	@echo "  make dev      - Setup development environment"
	@echo "  make help     - Show this help message" 