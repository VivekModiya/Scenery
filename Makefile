# Makefile for Prompt2Video

.PHONY: build run test clean docker-build docker-up docker-down deps

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=prompt2video
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/api/main.go

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/api/main.go
	./$(BINARY_NAME)

# Run with hot reload using air (install: go install github.com/cosmtrek/air@latest)
dev:
	air

# Test the application
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v cmd/api/main.go

# Docker commands
docker-build:
	docker build -t prompt2video .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

# Database migrations (requires golang-migrate)
migrate-up:
	migrate -path ./migrations -database "postgres://prompt2video:password123@localhost:5432/prompt2video?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://prompt2video:password123@localhost:5432/prompt2video?sslmode=disable" down

# Install tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Security check
security:
	gosec ./...

# Generate API documentation
docs:
	swag init -g cmd/api/main.go
