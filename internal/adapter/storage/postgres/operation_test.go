package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	postgres "github.com/evythrossell/account-management-api/internal/adapter/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPostgresOperationRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPostgresOperationRepository(db)
	ctx := context.Background()

	t.Run("Exists - True", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(int16(1)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repo.Exists(ctx, 1)

		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists - False", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(int16(99)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := repo.Exists(ctx, 99)

		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists - Infrastructure Error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS`).
			WithArgs(int16(1)).
			WillReturnError(errors.New("connection failed"))

		exists, err := repo.Exists(ctx, 1)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "infrastructure error")
	})
}
