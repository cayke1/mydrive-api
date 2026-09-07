.PHONY: help dev up down migrate logs clean seed seed-reset

# Database configuration
DB_HOST := localhost
DB_PORT := 5432
DB_USER := user
DB_PASSWORD := password
DB_NAME := mydrive
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH := ./infrastructure/postgres/migrations

help:
	@echo "MyDrive - Local Development Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev            - Start with hot reload (auto-rebuild on code changes)"
	@echo "  make dev-logs       - Show dev logs"
	@echo ""
	@echo "Production:"
	@echo "  make up             - Start containers in background"
	@echo "  make down           - Stop and remove containers"
	@echo "  make restart        - Restart all containers"
	@echo "  make logs           - Show container logs (follow mode)"
	@echo "  make clean          - Remove containers, volumes and images"
	@echo "  make ps             - Show running containers"
	@echo "  make api-logs       - Show API logs"
	@echo "  make web-logs       - Show Web logs"
	@echo ""
	@echo "Database/Migration commands:"
	@echo "  make migrate-up     - Run pending migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-force  - Force migration to specific version (usage: make migrate-force VERSION=1)"
	@echo "  make migrate-status - Show migration status"
	@echo ""
	@echo "Seed commands:"
	@echo "  make seed           - Load sample data into database"
	@echo "  make seed-reset     - Clear all tables and reload sample data"
	@echo ""

dev:
	@echo "Starting development environment with hot reload..."
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up

dev-logs:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

up:
	docker compose up -d
	@echo "✓ Containers started in background"
	@echo ""
	@echo "Services available at:"
	@echo "  - API:           http://localhost:8080"
	@echo "  - Web:           http://localhost:3000"
	@echo "  - MinIO Console: http://localhost:9001"
	@echo ""

down:
	docker compose down
	@echo "✓ Containers stopped"

restart:
	docker compose restart
	@echo "✓ Containers restarted"

logs:
	docker compose logs -f

clean:
	docker compose down -v
	@echo "✓ Containers, volumes and networks removed"

ps:
	docker compose ps

api-logs:
	docker compose logs -f api

web-logs:
	docker compose logs -f web

api-shell:
	docker compose exec api /bin/sh

db-shell:
	docker compose exec postgres psql -U user -d mydrive

redis-cli:
	docker compose exec redis redis-cli

# Database migration commands
migrate-up:
	@echo "Running migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up
	@echo "✓ Migrations completed"

migrate-down:
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1
	@echo "✓ Migration rolled back"

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make migrate-force VERSION=1"; \
	else \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION); \
		echo "✓ Migration forced to version $(VERSION)"; \
	fi

migrate-status:
	@migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

# Seed commands
seed:
	@echo "Loading sample data..."
	cat ./infrastructure/postgres/seed.sql | docker compose exec -T postgres psql -U $(DB_USER) -d $(DB_NAME)
	@echo "✓ Sample data loaded successfully"

seed-reset:
	@echo "Resetting database and reloading sample data..."
	docker compose exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) -c "TRUNCATE TABLE files CASCADE; TRUNCATE TABLE folders CASCADE; TRUNCATE TABLE users CASCADE;"
	@echo "Tables cleared"
	cat ./infrastructure/postgres/seed.sql | docker compose exec -T postgres psql -U $(DB_USER) -d $(DB_NAME)
	@echo "✓ Database reset and sample data reloaded"