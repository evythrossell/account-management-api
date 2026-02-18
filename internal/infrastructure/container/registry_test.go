package container_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	infrastructure "github.com/evythrossell/account-management-api/internal/infrastructure"
	"github.com/evythrossell/account-management-api/internal/infrastructure/container"
	pkg "github.com/evythrossell/account-management-api/pkg"
	"github.com/stretchr/testify/assert"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...pkg.Field)  {}
func (m *mockLogger) Error(msg string, fields ...pkg.Field) {}
func (m *mockLogger) Debug(msg string, fields ...pkg.Field) {}
func (m *mockLogger) Warn(msg string, fields ...pkg.Field)  {} // Adicionado
func (m *mockLogger) Fatal(msg string, fields ...pkg.Field) {} // Adicionado para resolver o erro

func TestContainer(t *testing.T) {
	cfg := &infrastructure.Config{
		DatabaseURL: "postgres://user:pass@localhost:5432/db?sslmode=disable",
	}

	testLogger := &mockLogger{}

	t.Run("New - Success", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
		mock.ExpectPing()

		c, err := container.New(cfg, testLogger)

		if err == nil {
			assert.NotNil(t, c.DB())
			assert.NotNil(t, c.Logger())
			assert.NotNil(t, c.AccountRepository())
			assert.NotNil(t, c.TransactionRepository())
			assert.NotNil(t, c.OperationRepository())
			assert.NotNil(t, c.AccountService())
			assert.NotNil(t, c.TransactionService())
			assert.NotNil(t, c.HealthService())
			assert.NotNil(t, c.AccountHandler())
			assert.NotNil(t, c.HealthHandler())
			assert.NotNil(t, c.TransactionHandler())
			assert.NoError(t, c.Close())
		}
		db.Close()
	})

	t.Run("Close - Nil DB", func(t *testing.T) {
		c := &container.Container{}
		assert.NoError(t, c.Close())
	})

	t.Run("New - Connection Error", func(t *testing.T) {
		invalidCfg := &infrastructure.Config{DatabaseURL: "invalid"}
		_, err := container.New(invalidCfg, testLogger)
		assert.Error(t, err)
	})
}
