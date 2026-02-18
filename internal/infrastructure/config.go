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
		ServerPort: getEnv("PORT", "8080"),
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", ""),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName:     getEnv("POSTGRES_DB", ""),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.DatabaseURL = buildDatabaseURL(cfg)

	return cfg, nil
}

func (c *Config) validate() error {
	requiredVars := map[string]string{
		"POSTGRES_USER":     c.DBUser,
		"POSTGRES_PASSWORD": c.DBPassword,
		"POSTGRES_HOST":     c.DBHost,
		"POSTGRES_DB":       c.DBName,
	}

	for key, value := range requiredVars {
		if value == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
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
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
