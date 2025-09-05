PROJECT = himawari-server
MAIN_DIR = ./cmd/server

# Load variable from .env
ifneq (,$(wildcard .env))
	include .env
	export
endif

# Commands
.PHONY: run build test clean

init:
	@echo "InitDB $(PROJECT)..."
	go run $(MAIN_DIR) --init-db

run:
	@echo "Running $(PROJECT) on port $(PORT)..."
	go run $(MAIN_DIR) --port $(PORT)

build:
	@echo "Building $(PROJECT)..."
	go build -o bin/$(PROJECT) $(MAIN_DIR)

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning..."
	rm -rf bin
