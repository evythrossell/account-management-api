package config_test

import (
	"os"
	"testing"

	"github.com/evythrossell/account-management-api/internal/infrastructure"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("Success - Load with all variables", func(t *testing.T) {
		os.Setenv("POSTGRES_USER", "user")
		os.Setenv("POSTGRES_PASSWORD", "pass")
		os.Setenv("POSTGRES_HOST", "localhost")
		os.Setenv("POSTGRES_DB", "db")
		os.Setenv("POSTGRES_PORT", "5432")
		os.Setenv("PORT", "9090")

		defer func() {
			os.Unsetenv("POSTGRES_USER")
			os.Unsetenv("POSTGRES_PASSWORD")
			os.Unsetenv("POSTGRES_HOST")
			os.Unsetenv("POSTGRES_DB")
			os.Unsetenv("POSTGRES_PORT")
			os.Unsetenv("PORT")
		}()

		cfg, err := config.Load()

		assert.NoError(t, err)
		assert.Equal(t, "user", cfg.DBUser)
		assert.Equal(t, "9090", cfg.ServerPort)
		assert.Contains(t, cfg.DatabaseURL, "postgres://user:pass@localhost:5432/db")
	})

	t.Run("Error - Missing required variable", func(t *testing.T) {
		os.Clearenv()

		cfg, err := config.Load()

		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "missing required environment variable")
	})

	t.Run("Success - Use default values", func(t *testing.T) {
		os.Setenv("POSTGRES_USER", "user")
		os.Setenv("POSTGRES_PASSWORD", "pass")
		os.Setenv("POSTGRES_DB", "db")

		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("PORT")

		defer os.Clearenv()

		cfg, err := config.Load()

		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, "5432", cfg.DBPort)
		assert.Equal(t, "8080", cfg.ServerPort)
	})
}
