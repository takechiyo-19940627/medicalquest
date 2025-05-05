package database

import (
    "fmt"
    
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/sql"
    
    "github.com/medicalquest/internal/config"
    "github.com/medicalquest/internal/ent"
)

// NewClient creates a new ent client for interacting with the database
func NewClient(cfg *config.Config) (*ent.Client, error) {
    // Create database connection string
    connectionString := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
    )
    
    // Connect to PostgreSQL with the connection string
    drv, err := sql.Open(dialect.Postgres, connectionString)
    if err != nil {
        return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
    }
    
    // Create an ent client
    return ent.NewClient(ent.Driver(drv)), nil
}