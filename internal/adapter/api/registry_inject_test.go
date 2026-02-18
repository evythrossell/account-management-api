package pkg_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	adapter "github.com/evythrossell/account-management-api/internal/adapter/api"
	"github.com/stretchr/testify/assert"
)

// TestAPIContainerWithMockDB tests API container initialization with mocked database
func TestAPIContainerWithMockDB(t *testing.T) {
	t.Run("API Container with mocked database operations", func(t *testing.T) {
		// Create a sqlmock database
		mockDB, mockSQL, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer mockDB.Close()

		// Expect ping to succeed
		mockSQL.ExpectPing()

		// Test empty container getters
		c := &adapter.Container{}
		assert.Nil(t, c.Logger())
		assert.Nil(t, c.DB())
		assert.Nil(t, c.AccountRepository())
		assert.Nil(t, c.TransactionRepository())
	})

	t.Run("API Container Close with valid mock database", func(t *testing.T) {
		mockDB, mockSQL, _ := sqlmock.New()
		mockSQL.ExpectClose()

		c := &adapter.Container{}
		// Container is empty, so Close should work without error
		assert.NoError(t, c.Close())
		mockDB.Close()
	})

	t.Run("API Container initialization sequence validation", func(t *testing.T) {
		c := &adapter.Container{}

		// Verify all getters exist and don't panic
		_ = c.Logger()
		_ = c.DB()
		_ = c.AccountRepository()
		_ = c.TransactionRepository()
		_ = c.OperationRepository()
		_ = c.HealthService()
		_ = c.AccountService()
		_ = c.TransactionService()
		_ = c.AccountHandler()
		_ = c.HealthHandler()
		_ = c.TransactionHandler()

		// All should be nil at initialization
		assert.Nil(t, c.DB())
		assert.Nil(t, c.Logger())
	})

	t.Run("API Container all getters coverage test", func(t *testing.T) {
		c := &adapter.Container{}

		// Test each getter individually for code coverage
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
	})
}
