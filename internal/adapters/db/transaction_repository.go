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

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Save(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	stmt := `INSERT INTO transactions (account_id, operation_type_id, amount, event_date) 
			VALUES ($1, $2, $3, $4) RETURNING transaction_id`

	err := r.db.QueryRowContext(ctx, stmt, transaction.AccountID, transaction.OperationTypeID, transaction.Amount, transaction.EventDate).
		Scan(&transaction.TransactionID)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return nil, domainerror.NewConflictError("account id already registered", err)
			}
		}
		return nil, fmt.Errorf("save transaction: %w", err)
	}
	return transaction, nil
}
