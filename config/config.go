package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	Environment string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:  getEnv("PORT", "8080"),
		DBHost:      getEnv("POSTGRES_HOST", "localhost"),
		DBPort:      getEnv("POSTGRES_PORT", "5432"),
		DBUser:      getEnv("POSTGRES_USER", "postgres"),
		DBPassword:  getEnv("POSTGRES_PASSWORD", ""),
		DBName:      getEnv("POSTGRES_DB", "accountmanagementapi"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	cfg.DatabaseURL = buildDatabaseURL(cfg)

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.DBUser == "" {
		return fmt.Errorf("POSTGRES_USER is required")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is required")
	}
	if c.DBHost == "" {
		return fmt.Errorf("POSTGRES_HOST is required")
	}
	if c.DBPort == "" {
		return fmt.Errorf("POSTGRES_PORT is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("POSTGRES_DB is required")
	}
	if c.ServerPort == "" {
		return fmt.Errorf("PORT is required")
	}

	return nil
}

func buildDatabaseURL(c *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
