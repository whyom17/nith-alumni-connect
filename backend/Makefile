DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/alumni_connect?sslmode=disable
MIGRATIONS_DIR ?= db/migrations

.PHONY: help run test lint sqlc migrate-up migrate-down migrate-status migrate-new docker-db docker-down

help:
	@echo "Available commands:"
	@echo "  make run            - Run the API server"
	@echo "  make test           - Run unit and integration tests with race detector"
	@echo "  make lint           - Run linter (golangci-lint and go vet)"
	@echo "  make sqlc           - Generate Go code from SQL queries via sqlc"
	@echo "  make migrate-up     - Run pending goose migrations"
	@echo "  make migrate-down   - Roll back the latest migration"
	@echo "  make migrate-status - Show current migration status"
	@echo "  make migrate-new name=your_migration - Create a new SQL migration"
	@echo "  make docker-db      - Start local Postgres container"
	@echo "  make docker-down    - Stop local Postgres container"

run:
	go run cmd/api/main.go

test:
	go test -race -v ./...

lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipped"; \
	fi

sqlc:
	sqlc generate

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

migrate-new:
	@if [ -z "$(name)" ]; then echo "Error: 'name' is required (e.g., make migrate-new name=add_users)"; exit 1; fi
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

docker-db:
	docker compose up -d db

docker-down:
	docker compose down
