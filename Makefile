.PHONY: infra infra-down docker docker-down logs run

# === Local development ===
# Start only infrastructure (postgres, redis, kafka)
infra:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d
	@echo ""
	@echo "Infrastructure is ready!"
	@echo "Run your app:  make run"

# Stop infrastructure
infra-down:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

# Run Go app locally (loads .env + .env.local)
run:
	@set -a && . ./.env && . ./.env.local && set +a && go run cmd/app/main.go

# === Full stack in Docker ===
# Start everything (infra + app container)
docker:
	docker compose up -d

docker-down:
	docker compose down

# === Utilities ===
logs:
	docker compose logs -f

logs-kafka:
	docker compose logs -f kafka