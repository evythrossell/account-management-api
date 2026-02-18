package db_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	postgres "github.com/evythrossell/account-management-api/internal/adapter/storage/postgres"
	"github.com/evythrossell/account-management-api/internal/core/domain"
	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresAccountRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresAccountRepository(db)
	ctx := context.Background()

	t.Run("Save - Success", func(t *testing.T) {
		acc := &domain.Account{DocumentNumber: "123"}
		mock.ExpectQuery("INSERT INTO accounts").
			WithArgs(acc.DocumentNumber).
			WillReturnRows(sqlmock.NewRows([]string{"account_id"}).AddRow(1))

		result, err := repo.Save(ctx, acc)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.ID)
	})

	t.Run("Save - Duplicate Document", func(t *testing.T) {
		acc := &domain.Account{DocumentNumber: "123"}
		mock.ExpectQuery("INSERT INTO accounts").
			WithArgs(acc.DocumentNumber).
			WillReturnError(&pq.Error{Code: "23505"})

		result, err := repo.Save(ctx, acc)
		assert.ErrorIs(t, err, common.ErrAccountAlreadyExists)
		assert.Nil(t, result)
	})

	t.Run("Save - Generic Error", func(t *testing.T) {
		acc := &domain.Account{DocumentNumber: "123"}
		mock.ExpectQuery("INSERT INTO accounts").
			WithArgs(acc.DocumentNumber).
			WillReturnError(errors.New("db error"))

		_, err := repo.Save(ctx, acc)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("FindByDocument - Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs("123").
			WillReturnRows(sqlmock.NewRows([]string{"account_id", "document_number"}).AddRow(1, "123"))

		result, err := repo.FindByDocument(ctx, "123")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.ID)
	})

	t.Run("FindByDocument - Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs("123").
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByDocument(ctx, "123")
		assert.ErrorIs(t, err, common.ErrAccountNotFound)
		assert.Nil(t, result)
	})

	t.Run("FindByDocument - Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs("123").
			WillReturnError(errors.New("db error"))

		_, err := repo.FindByDocument(ctx, "123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "infrastructure error")
	})

	t.Run("FindByAccountID - Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"account_id", "document_number"}).AddRow(1, "123"))

		result, err := repo.FindByAccountID(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result.ID)
	})

	t.Run("FindByAccountID - Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs(int64(1)).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.FindByAccountID(ctx, 1)
		assert.ErrorIs(t, err, common.ErrAccountNotFound)
		assert.Nil(t, result)
	})

	t.Run("FindByAccountID - Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts").
			WithArgs(int64(1)).
			WillReturnError(errors.New("db error"))

		_, err := repo.FindByAccountID(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}
