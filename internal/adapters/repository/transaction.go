package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	domainerror "github.com/evythrossell/account-management-api/internal/core/domain/error"
	"github.com/lib/pq"
)

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (p *PostgresTransactionRepository) Save(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	stmt := `INSERT INTO transactions (account_id, operation_type_id, amount, event_date) 
			VALUES ($1, $2, $3, $4) RETURNING transaction_id`

	err := p.db.QueryRowContext(ctx, stmt,
		transaction.AccountID,
		transaction.OperationTypeID,
		transaction.Amount,
		transaction.EventDate,
	).Scan(&transaction.ID)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return nil, fmt.Errorf("%w: %v", domainerror.ErrAccountNotFound, err)
			}
		}
		return nil, fmt.Errorf("infrastructure error: failed to save transaction: %w", err)
	}
	return transaction, nil
}

func (p *PostgresTransactionRepository) FindByTransactionID(ctx context.Context, transactionID int64) (*domain.Transaction, error) {
	stmt := `SELECT transaction_id, account_id, operation_type_id, amount FROM transactions WHERE transaction_id = $1`

	var tx domain.Transaction
	err := p.db.QueryRowContext(ctx, stmt, transactionID).Scan(
		&tx.ID,
		&tx.AccountID,
		&tx.OperationTypeID,
		&tx.Amount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerror.ErrTransactionNotFound
		}
		return nil, fmt.Errorf("infrastructure error: failed to find transaction: %w", err)
	}
	return &tx, nil
}
