.PHONY: up down start stop

up:
	@docker compose up

down:
	@docker compose down

# To rebuild the images before starting
start:
	@docker compose up --build

# To remove named volumes (postgres, redis)
stop:
	@docker compose down -v