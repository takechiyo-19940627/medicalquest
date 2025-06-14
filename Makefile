.PHONY: up down build run test lint generate migrate-dev migrate-prod

# Docker related commands
up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose build

# Database migration commands
migrate-hash:
	atlas migrate hash --dir "file://migrations"
migrate-apply:
	atlas migrate apply --env dev
migrate-diff:
	atlas migrate diff ${MIGRATION_NAME} \
		--env dev \
		--dev-url "docker://postgres/15/dev?search_path=public"
migrate-status:
	atlas migrate status \
		--env dev \
		--url "postgres://postgres:postgres@localhost:5432/medicalquest?sslmode=disable"
