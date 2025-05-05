#!/bin/sh

# Make the script executable with:
# chmod +x scripts/generate-ent.sh

# Create ent directory if it doesn't exist
mkdir -p internal/ent

# Generate ent code from schema
go run -mod=mod entgo.io/ent/cmd/ent generate ./internal/ent/schema