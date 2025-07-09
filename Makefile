.PHONY: migration, lint, build

MIGRATIONS_PATH := cmd/migrations/

lint:
	@golangci-lint fmt && golangci-lint run

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

build:
	@go build -o bin/order-service cmd/api/main.go
