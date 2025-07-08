.PHONY: migration, lint

MIGRATIONS_PATH := cmd/migrations/

lint:
	@golangci-lint fmt && golangci-lint run

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))
