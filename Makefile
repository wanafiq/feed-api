SHELL := /bin/bash

DATABASE_URL=postgres://blog:blog@localhost/blog?sslmode=disable
MIGRATIONS_DIR=./migrations
TMP_DIR=./tmp
BIN_DIR=./bin

.PHONY: help
all: help

.PHONY: migration migrate-up migrate-down
migration:
	@read -p "Enter migration name: " name; \
    migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $$name
migrate-up:
	@migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) up
migrate-down:
	@migrate -database $(DATABASE_URL) -path $(MIGRATIONS_DIR) down

.PHONY: clean build-dev build
clean:
	rm -rf $(TMP_DIR)
	rm -rf $(BIN_DIR)
build-dev: clean
	mkdir -p $(TMP_DIR)
	go build -o $(TMP_DIR)/main ./cmd/main.go
build: clean
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/main ./cmd/main.go

help:
	@echo "Makefile for project tasks:"
	@echo "  migrations      	- Create a new migration file"
	@echo "  migrate-up         - Apply all pending database migrations"
	@echo "  migrate-down       - Rollback the last database migration"
	@echo "  swagger-gen        - Generate Swagger documentation"
	@echo "  swagger-clean      - Remove all generated Swagger documentation"
	@echo "  build-dev      	- Build app for development"
	@echo "  build      		- Build app for production"
	@echo "  help               - Show this help message"





