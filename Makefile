.PHONY: migration

MIGRATIONS_PATH := cmd/migrations/

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))
