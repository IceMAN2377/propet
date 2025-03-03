.PHONY: build clean run docker-build docker-run docker-down fmt vet

# Main commands
all: clean build

# Build both services
build: build-api build-task

# Build API Service
build-api:
	@echo "Building API Service..."
	cd ApiService && go build -o bin/api-service main.go

# Build Task Service
build-task:
	@echo "Building Task Service..."
	cd TaskService && go build -o bin/task-service ./cmd/myapp

# Clean up
clean:
	@echo "Cleaning up..."
	rm -rf ApiService/bin TaskService/bin
	mkdir -p ApiService/bin TaskService/bin

# Run services locally (for development)
run-api:
	@echo "Running API Service..."
	cd ApiService && go run main.go

run-task:
	@echo "Running Task Service..."
	cd TaskService && go run cmd/myapp/main.go

# Format code
fmt:
	@echo "Formatting code..."
	cd ApiService && go fmt ./...
	cd TaskService && go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	cd ApiService && go vet ./...
	cd TaskService && go vet ./...

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-run:
	@echo "Starting services with Docker..."
	docker-compose up

docker-background:
	@echo "Starting services with Docker in background..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-clean:
	@echo "Stopping Docker services and removing volumes..."
	docker-compose down -v

# Check if everything is working
check: docker-build docker-background
	@echo "Waiting for services to start..."
	@sleep 5
	@echo "Making test request to API Gateway..."
	@curl -X GET http://localhost:8000/tasks/1 || echo "Request failed"
	@echo "\nMaking direct request to Task Service..."
	@curl -X GET http://localhost:8080/tasks/1 || echo "Request failed"
	@echo "\nChecking docker logs..."
	@docker-compose logs --tail=10

# Help
help:
	@echo "Available commands:"
	@echo "  make build          - Build both services"
	@echo "  make clean          - Clean up binaries"
	@echo "  make fmt            - Format Go code"
	@echo "  make vet            - Run Go vet on code"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-run     - Start services with Docker"
	@echo "  make docker-down    - Stop Docker services"
	@echo "  make check          - Build, run and check if services are working"
	@echo "  make help           - Show this help" 