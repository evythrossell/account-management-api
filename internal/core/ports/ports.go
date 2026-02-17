package ports

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/domain"
)

type AccountRepository interface {
	Save(ctx context.Context, account *domain.Account) (*domain.Account, error)
	FindByDocument(ctx context.Context, documentNumber string) (*domain.Account, error)
	FindByAccountID(ctx context.Context, accountId int64) (*domain.Account, error)
}

type TransactionRepository interface {
	Save(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	FindByTransactionID(ctx context.Context, transactionID int64) (*domain.Transaction, error)
}

type OperationRepository interface {
	Exists(ctx context.Context, operationType int16) (bool, error)
}

type AccountService interface {
	CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error)
	GetAccount(ctx context.Context, documentNumber string) (*domain.Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error)
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, accountID int64, operationType int16, amount float64) (*domain.Transaction, error)
	GetByTransactionID(ctx context.Context, transactionID int64) (*domain.Transaction, error)
}
