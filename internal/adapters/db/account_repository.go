package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
	"github.com/lib/pq"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Save(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	stmt := `INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id`

	err := r.db.QueryRowContext(ctx, stmt, account.DocumentNumber).Scan(&account.ID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, domainerror.NewConflictError("document number already registered", err)
			}
		}
		return nil, fmt.Errorf("save account: %w", err)
	}
	return account, nil
}

func (r *accountRepository) FindByDocument(ctx context.Context, documentNumber string) (*domain.Account, error) {
	stmt := `SELECT account_id, document_number FROM accounts WHERE document_number = $1`

	var acc domain.Account
	err := r.db.QueryRowContext(ctx, stmt, documentNumber).Scan(&acc.ID, &acc.DocumentNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find account by document: %w", err)
	}
	return &acc, nil
}

func (r *accountRepository) FindByAccountID(ctx context.Context, accountID int64) (*domain.Account, error) {
	var acc domain.Account
	row := r.db.QueryRowContext(ctx, "SELECT account_id, document_number FROM accounts WHERE account_id=$1", accountID)
	if err := row.Scan(&acc.ID, &acc.DocumentNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find account by id: %w", err)
	}
	return &acc, nil
}
