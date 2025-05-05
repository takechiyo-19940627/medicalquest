# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands
- Build: `docker-compose up -d --build`
- Lint: `go vet ./...` and `golangci-lint run ./...`
- Test: `go test -v ./...`
- Single test: `go test -v ./path/to/package -run TestName`
- Generate Ent code: `./scripts/generate-ent.sh`
- Run migrations: `make migrate-dev`

## Code Style Guidelines
- Language: Japanese for comments, English for code
- Indentation: 4 spaces
- Line length: Max 100 characters
- Naming: 
  - camelCase for variables and functions
  - PascalCase for classes
  - snake_case for database fields
- Imports: Group by type, alphabetize within groups
- Error handling: Use try/catch blocks with specific error types
- Types: Use strong typing where available
- Database: Follow the ERD in ARCHITECTURE.md
- Documentation: Document all public functions with clear descriptions

## Project Structure
- Question/Choice model relationship as defined in docs/ARCHITECTURE.md