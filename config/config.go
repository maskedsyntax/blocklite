package config

import (
	"os"
)

// Configuration settings for the application
type Config struct {
	Port string
}

// Load the configuration from environment variables or defaults
func LoadConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"), // Default to port 8080 if not set
	}
}

// Retrieve the value of the environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
