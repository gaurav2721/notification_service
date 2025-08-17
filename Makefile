.PHONY: build run test docker-build docker-run docker-exec docker-shell

# Build the application
build:
	go build -o bin/notification-service main.go

# Run the application
run:
	go run main.go

# Run the application with debug output
run-debug:
	LOG_LEVEL=debug go run main.go

# Run tests
test:
	go test -v ./...

# Build Docker image
docker-build:
	docker build -t notification-service .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env notification-service

# Exec into running Docker container to view output files
docker-exec:
	docker exec -it $$(docker ps -q --filter ancestor=notification-service) sh -c "echo '=== Mock Service Output Files ===' && ls -la output/ 2>/dev/null || echo 'Output directory not found, creating it...' && mkdir -p output && echo '=== Email notifications ===' && cat output/email.txt 2>/dev/null || echo 'No email notifications yet' && echo '=== Slack notifications ===' && cat output/slack.txt 2>/dev/null || echo 'No slack notifications yet' && echo '=== APNS notifications ===' && cat output/apns.txt 2>/dev/null || echo 'No APNS notifications yet' && echo '=== FCM notifications ===' && cat output/fcm.txt 2>/dev/null || echo 'No FCM notifications yet'"

# Exec into running Docker container shell
docker-shell:
	docker exec -it $$(docker ps -q --filter ancestor=notification-service) /bin/sh

# Help
help:
	@echo "Available commands:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  run-debug          - Run the application with debug output"
	@echo "  test               - Run tests"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-exec        - Exec into running container and view output files"
	@echo "  docker-shell       - Exec into running container shell"
	@echo "  help               - Show this help message" 