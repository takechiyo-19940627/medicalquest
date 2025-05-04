.PHONY: up down build run test lint generate migrate-dev migrate-prod

# Docker related commands
up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose build

# Development commands
run:
	go run cmd/api/main.go

test:
	go test -v ./...

lint:
	go vet ./...
	golangci-lint run ./...

# Code generation
generate:
	go generate ./...

# Database migration commands
migrate-dev:
	atlas migrate apply \
		--dir file://migrations \
		--url "postgres://postgres:postgres@localhost:5432/medicalquest?sslmode=disable"

migrate-prod:
	atlas migrate apply \
		--dir file://migrations \
		--url "${DATABASE_URL}"