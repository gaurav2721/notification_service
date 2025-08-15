.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Build the application
build:
	go build -o bin/notification-service main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Build Docker image
docker-build:
	docker build -t notification-service .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env notification-service

# Start services with Docker Compose
docker-compose-up:
	docker-compose up --build

# Stop services with Docker Compose
docker-compose-down:
	docker-compose down

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate documentation
docs:
	godoc -http=:6060

# Create environment file from example
env:
	cp env.example .env

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  clean              - Clean build artifacts"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-compose-up  - Start services with Docker Compose"
	@echo "  docker-compose-down- Stop services with Docker Compose"
	@echo "  deps               - Install dependencies"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code"
	@echo "  docs               - Generate documentation"
	@echo "  env                - Create environment file from example"
	@echo "  help               - Show this help message" 