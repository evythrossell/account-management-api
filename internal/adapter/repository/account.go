package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/evythrossell/account-management-api/internal/core/common"
	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/lib/pq"
)

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (p *PostgresAccountRepository) Save(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	stmt := `INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id`

	err := p.db.QueryRowContext(ctx, stmt, account.DocumentNumber).Scan(&account.ID)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("%w: %v", common.ErrAccountAlreadyExists, err)
			}
		}
		return nil, fmt.Errorf("save account: %w", err)
	}
	return account, nil
}

func (p *PostgresAccountRepository) FindByDocument(ctx context.Context, documentNumber string) (*domain.Account, error) {
	stmt := `SELECT account_id, document_number FROM accounts WHERE document_number = $1`

	var acc domain.Account
	err := p.db.QueryRowContext(ctx, stmt, documentNumber).Scan(&acc.ID, &acc.DocumentNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrAccountNotFound
		}
		return nil, fmt.Errorf("infrastructure error: find account by document: %w", err)
	}

	return &acc, nil
}

func (p *PostgresAccountRepository) FindByAccountID(ctx context.Context, accountID int64) (*domain.Account, error) {
	stmt := `SELECT account_id, document_number FROM accounts WHERE account_id = $1`
	var acc domain.Account

	err := p.db.QueryRowContext(ctx, stmt, accountID).Scan(&acc.ID, &acc.DocumentNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.ErrAccountNotFound
		}
		return nil, fmt.Errorf("infrastructure error: find account by id: %w", err)
	}

	return &acc, nil
}
