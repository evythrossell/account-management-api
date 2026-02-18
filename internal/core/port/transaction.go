package port

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/domain"
)

type TransactionRepository interface {
	Save(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	FindByTransactionID(ctx context.Context, transactionID int64) (*domain.Transaction, error)
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, accountID int64, operationType int16, amount float64) (*domain.Transaction, error)
	GetByTransactionID(ctx context.Context, transactionID int64) (*domain.Transaction, error)
}
