# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=gophernet
MAIN_PATH=./cmd/main

# Mock parameters
MOCKGEN=mockgen
MOCK_DIR=pkg/mocks

.PHONY: all deps install-mockgen generate-mocks migrate run

# Main setup command that handles everything
all: deps install-mockgen generate-mocks migrate build
	@echo "Setup complete!"

# Dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy
	
# Migration commands
migrate:
	$(GOCMD) run -mod=mod entgo.io/ent/cmd/ent generate ./pkg/db/ent/schema

# Mock commands
install-mockgen:
	$(GOGET) github.com/golang/mock/mockgen@latest

generate-mocks:
	@mkdir -p $(MOCK_DIR)
	$(MOCKGEN) -source=pkg/repo/burrow.go -destination=$(MOCK_DIR)/burrow_mock.go -package=mocks

# Run the application
run: build
	./$(BINARY_NAME)

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
