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

.PHONY: all build clean test run dev deps db-create db-drop db-reset

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