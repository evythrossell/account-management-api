package pkg_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	adapter "github.com/evythrossell/account-management-api/internal/adapter/api"
	infrastructure "github.com/evythrossell/account-management-api/internal/infrastructure"
	pkg "github.com/evythrossell/account-management-api/pkg"
	"github.com/stretchr/testify/assert"
)

type mockLoggerAPI struct{}

func (m *mockLoggerAPI) Info(msg string, fields ...pkg.Field)  {}
func (m *mockLoggerAPI) Error(msg string, fields ...pkg.Field) {}
func (m *mockLoggerAPI) Debug(msg string, fields ...pkg.Field) {}
func (m *mockLoggerAPI) Warn(msg string, fields ...pkg.Field)  {}
func (m *mockLoggerAPI) Fatal(msg string, fields ...pkg.Field) {}

func TestAPIContainer(t *testing.T) {
	cfg := &infrastructure.Config{
		DatabaseURL: "postgres://user:pass@localhost:5432/db?sslmode=disable",
	}

	testLogger := &mockLoggerAPI{}

	t.Run("New - Success", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
		mock.ExpectPing()

		c, err := adapter.New(cfg, testLogger)

		if err == nil {
			assert.NotNil(t, c.DB())
			assert.NotNil(t, c.Logger())
			assert.NotNil(t, c.AccountRepository())
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

	t.Run("New - API initialization", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
		mock.ExpectPing()

		c, err := adapter.New(cfg, testLogger)
		if err == nil {
			assert.NotNil(t, c.Logger())
			assert.NotNil(t, c.DB())
			assert.NotNil(t, c.AccountRepository())
			assert.NoError(t, c.Close())
		}
		db.Close()
	})

	t.Run("Close - Nil DB", func(t *testing.T) {
		c := &adapter.Container{}
		assert.NoError(t, c.Close())
	})

	t.Run("New - Connection Error", func(t *testing.T) {
		invalidCfg := &infrastructure.Config{DatabaseURL: "invalid"}
		_, err := adapter.New(invalidCfg, testLogger)
		assert.Error(t, err)
	})

	t.Run("New - Ping Error", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
		mock.ExpectPing().WillReturnError(errors.New("connection refused"))

		invalidCfg := &infrastructure.Config{DatabaseURL: "postgres://invalid"}
		_, err := adapter.New(invalidCfg, testLogger)
		assert.Error(t, err)
		db.Close()
	})
}

func TestAPIContainerErrorHandling(t *testing.T) {
	testLogger := &mockLoggerAPI{}

	t.Run("New - Error Handling (nil config)", func(t *testing.T) {
		_, err := adapter.New(nil, testLogger)
		assert.Error(t, err)
	})

	t.Run("New - Error Handling (nil logger)", func(t *testing.T) {
		cfg := &infrastructure.Config{DatabaseURL: "postgres://user:pass@localhost:5432/db?sslmode=disable"}
		_, err := adapter.New(cfg, nil)
		assert.Error(t, err)
	})
}

func TestAPIContainerClose(t *testing.T) {
	t.Run("Close - Nil DB", func(t *testing.T) {
		c := &adapter.Container{}
		assert.NoError(t, c.Close())
	})
}

func TestAPIContainerGetters(t *testing.T) {
	c := &adapter.Container{}

	assert.Nil(t, c.Logger())
	assert.Nil(t, c.DB())
	assert.Nil(t, c.AccountRepository())
	assert.Nil(t, c.TransactionRepository())
	assert.Nil(t, c.OperationRepository())
	assert.Nil(t, c.HealthService())
	assert.Nil(t, c.AccountService())
	assert.Nil(t, c.TransactionService())
	assert.Nil(t, c.AccountHandler())
	assert.Nil(t, c.HealthHandler())
	assert.Nil(t, c.TransactionHandler())
}
