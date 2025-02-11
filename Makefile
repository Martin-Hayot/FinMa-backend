# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	
	
	@go build -o main.exe main.go

# Run the application
run:
	@go run main.go


# Create DB container
docker-run:
	docker compose up -d

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi


# Test the application
test:
	@echo "Testing..."
	@go test ./... -v


# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
	@air


.PHONY: all build run test clean watch
