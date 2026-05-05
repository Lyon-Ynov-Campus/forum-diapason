.PHONY: help build run test clean install-deps migrate

help:
	@echo "Forum Diapason Backend - Makefile"
	@echo ""
	@echo "Available commands:"
	@echo "  make build         Build the backend binary"
	@echo "  make run           Run the backend in development mode"
	@echo "  make test          Run all tests"
	@echo "  make test-coverage Run tests with coverage report"
	@echo "  make clean         Clean build artifacts"
	@echo "  make install-deps  Install dependencies"
	@echo "  make migrate       Run database migrations"
	@echo "  make format        Format code with go fmt"
	@echo "  make lint          Lint code"

build:
	@echo "Building backend..."
	@mkdir -p bin
	@go build -o bin/forum-api ./cmd/api

run: install-deps
	@echo "Running backend in development mode..."
	@go run ./cmd/api

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f *.db *.db-shm *.db-wal
	@rm -f coverage.out coverage.html

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

migrate:
	@echo "Running migrations..."
	@go run ./cmd/api/migrate

format:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Linting code..."
	@go vet ./...
