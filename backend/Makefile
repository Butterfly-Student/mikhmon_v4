.PHONY: build run test clean docker-build docker-run migrate

# Build the application
build:
	go build -o bin/api ./cmd/api/main.go

# Run the application
run:
	go run ./cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t mikhmon-api:latest .

# Run Docker container
docker-run:
	docker-compose up -d

# Stop Docker containers
docker-stop:
	docker-compose down

# Run database migrations
migrate:
	go run ./cmd/migrate/main.go

# Generate Swagger docs
swagger:
	swag init -g ./cmd/api/main.go

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy
