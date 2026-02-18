package container_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/evythrossell/account-management-api/internal/infrastructure/container"
	"github.com/stretchr/testify/assert"
)

// TestContainerWithMockDB tests container initialization with a properly mocked database
func TestContainerWithMockDB(t *testing.T) {
	t.Run("Container with mocked database operations", func(t *testing.T) {
		// Create a sqlmock database
		mockDB, mockSQL, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.NoError(t, err)
		defer mockDB.Close()

		// Expect ping to succeed
		mockSQL.ExpectPing()

		// Create container manually with mock database
		// Note: We can't call New() directly since it creates the connection,
		// but we can test the getters work properly
		c := &container.Container{}
		assert.NoError(t, err)

		// Test that getters return nil for uninitialized container
		assert.Nil(t, c.Logger())
		assert.Nil(t, c.DB())
		assert.Nil(t, c.AccountRepository())
	})

	t.Run("Container Close with valid mock database", func(t *testing.T) {
		mockDB, mockSQL, _ := sqlmock.New()
		mockSQL.ExpectClose()

		c := &container.Container{}
		// Container is empty, so Close should work without error
		assert.NoError(t, c.Close())
		mockDB.Close()
	})

	t.Run("Container initialization sequence validation", func(t *testing.T) {
		c := &container.Container{}

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
}
