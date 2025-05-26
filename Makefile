# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=gophernet
MAIN_PATH=./cmd/main

# Database parameters
DB_NAME=gophernet
DB_USER=postgres
DB_PASSWORD=postgres
DB_PORT=5432

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER=docker

# Mock parameters
MOCKGEN=mockgen
MOCK_DIR=pkg/mocks

.PHONY: all build clean test run dev deps db-create db-drop db-reset docker-build docker-up docker-down docker-logs docker-clean install-mockgen generate-mocks

all: clean build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

run: build
	./$(BINARY_NAME)

dev:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Database commands
db-create:
	PGPASSWORD=$(DB_PASSWORD) createdb -h localhost -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

db-drop:
	PGPASSWORD=$(DB_PASSWORD) dropdb -h localhost -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

db-reset: db-drop db-create

# Swagger documentation
swagger:
	swag init -g cmd/main/main.go -o docs

# Docker commands
docker-build:
	$(DOCKER_COMPOSE) build

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-clean:
	$(DOCKER_COMPOSE) down -v
	$(DOCKER) system prune -f

# Mock commands
install-mockgen:
	$(GOGET) github.com/golang/mock/mockgen@latest

generate-mocks:
	@mkdir -p $(MOCK_DIR)
	$(MOCKGEN) -source=pkg/repo/burrow.go -destination=$(MOCK_DIR)/burrow_mock.go -package=mocks

# Help command
help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make clean      - Clean build files"
	@echo "  make test       - Run tests"
	@echo "  make run        - Build and run the application"
	@echo "  make dev        - Build and run in development mode"
	@echo "  make deps       - Download dependencies"
	@echo "  make db-create  - Create database"
	@echo "  make db-drop    - Drop database"
	@echo "  make db-reset   - Reset database (drop and create)"
	@echo "  make swagger    - Generate Swagger documentation"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-logs  - View Docker container logs"
	@echo "  make docker-clean - Clean up Docker resources"
	@echo "  make install-mockgen - Install mockgen tool"
	@echo "  make generate-mocks  - Generate mock files" 