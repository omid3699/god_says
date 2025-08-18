.PHONY: build test test-coverage test-race  clean run run-server 


# Build variables
CLI_BINARY := ./bin/godsays
CLI_PATH := ./cmd/main.go
# Build flags
BUILD_FLAGS := -ldflags="-s -w"


help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


build: ## Build the binary
	@echo "Building God Says Binary"
	go build $(BUILD_FLAGS) -o $(CLI_BINARY) $(CLI_PATH)


clean: ## Clean built binaries
	@echo "Cleaning binaries..."
	rm -f $(CLI_BINARY)


test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -race -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -cover ./...


run: build ## Run the God Says
	./$(CLI_BINARY)

run-server: build ## Run God Says server
	./$(CLI_BINARY) -http

install: build ## Install binaries to GOPATH/bin
	@echo "Installing binaries..."
	go install $(CLI_PATH)
	go install $(SERVER_PATH)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

lint: fmt vet ## Run formatting and vetting


.DEFAULT_GOAL := help
