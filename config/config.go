package config

import (
    "os"
)

// Config holds all configuration for the application
type Config struct {
    // Server settings
    ServerPort string

    // Database settings
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
}

// New returns a new Config struct with values from environment variables
func New() *Config {
    return &Config{
        ServerPort: getEnv("SERVER_PORT", "8080"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "medicalquest"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
    }
}

// Helper function to get environment variables with a fallback value
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
