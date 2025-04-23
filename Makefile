SHELL := /bin/bash

DATABASE_URL=postgres://blog:blog@localhost/blog?sslmode=disable
MIGRATIONS_DIR=./migrations

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

help:
	@echo "Makefile for project tasks:"
	@echo "  migrations      	- Create a new migration file"
	@echo "  migrate-up         - Apply all pending database migrations"
	@echo "  migrate-down       - Rollback the last database migration"
	@echo "  swagger-gen        - Generate Swagger documentation"
	@echo "  swagger-clean      - Remove all generated Swagger documentation"
	@echo "  help               - Show this help message"





