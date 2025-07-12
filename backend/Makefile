.PHONY: migration lint build migrate-up migrate-down start stop
include .env

MIGRATIONS_PATH := cmd/migrations/
MIGRATE_CMD = migrate -path=$(MIGRATIONS_PATH) -database=$(PGSQL_URL)

run:
	@$(MAKE) lint && go run cmd/api/main.go

build:
	@$(MAKE) lint && go build -o bin/order-service cmd/api/main.go

lint:
	@golangci-lint fmt && golangci-lint run

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@$(MIGRATE_CMD) up

migrate-down:
	@$(MIGRATE_CMD) down

start:
	@docker compose up --build

stop:
	@docker compose down