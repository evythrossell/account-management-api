package db

import (
	"context"
	"database/sql"

	"github.com/evythrossell/account-management-api/internal/core/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Save(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	stmt := `INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id`

	err := r.db.QueryRowContext(ctx, stmt, account.DocumentNumber).Scan(&account.ID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) FindByDocument(ctx context.Context, documentNumber string) (*domain.Account, error) {
	stmt := `SELECT account_id, document_number FROM accounts WHERE document_number = $1`

	var acc domain.Account
	err := r.db.QueryRowContext(ctx, stmt, documentNumber).Scan(&acc.ID, &acc.DocumentNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}
