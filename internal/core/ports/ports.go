package ports

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/domain"
)

type AccountRepository interface {
	Save(ctx context.Context, account *domain.Account) (*domain.Account, error)
	FindByDocument(ctx context.Context, documentNumber string) (*domain.Account, error)
}

type AccountService interface {
	CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error)
}
