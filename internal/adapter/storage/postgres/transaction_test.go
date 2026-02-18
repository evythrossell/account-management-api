package db_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	postgres "github.com/evythrossell/account-management-api/internal/adapter/storage/postgres"
	"github.com/evythrossell/account-management-api/internal/core/domain"
	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresTransactionRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresTransactionRepository(db)
	ctx := context.Background()

	t.Run("Save - Success", func(t *testing.T) {
		tx := &domain.Transaction{
			AccountID:       1,
			OperationTypeID: 4,
			Amount:          123.45,
			EventDate:       time.Now(),
		}

		mock.ExpectQuery("INSERT INTO transactions").
			WithArgs(tx.AccountID, tx.OperationTypeID, tx.Amount, tx.EventDate).
			WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).AddRow(100))

		result, err := repo.Save(ctx, tx)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), result.ID)
	})

	t.Run("Save - Foreign Key Violation (23503)", func(t *testing.T) {
		tx := &domain.Transaction{AccountID: 999}
		mock.ExpectQuery("INSERT INTO transactions").
			WillReturnError(&pq.Error{Code: "23503"})

		result, err := repo.Save(ctx, tx)

		assert.ErrorIs(t, err, common.ErrAccountNotFound)
		assert.Nil(t, result)
	})

	t.Run("Save - Generic Error", func(t *testing.T) {
		tx := &domain.Transaction{AccountID: 1}
		mock.ExpectQuery("INSERT INTO transactions").
			WillReturnError(errors.New("db connection lost"))

		result, err := repo.Save(ctx, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "infrastructure error")
		assert.Nil(t, result)
	})

	t.Run("FindByTransactionID - Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM transactions").
			WithArgs(int64(100)).
			WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "account_id", "operation_type_id", "amount"}).
				AddRow(100, 1, 4, 123.45))

		result, err := repo.FindByTransactionID(ctx, 100)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), result.ID)
		assert.Equal(t, 123.45, result.Amount)
	})

	t.Run("FindByTransactionID - Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM transactions").
			WithArgs(int64(100)).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByTransactionID(ctx, 100)

		assert.ErrorIs(t, err, common.ErrTransactionNotFound)
		assert.Nil(t, result)
	})

	t.Run("FindByTransactionID - Generic Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM transactions").
			WithArgs(int64(100)).
			WillReturnError(errors.New("timeout"))

		result, err := repo.FindByTransactionID(ctx, 100)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "infrastructure error")
		assert.Nil(t, result)
	})
}
