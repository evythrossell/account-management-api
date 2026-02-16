package services

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/evythrossell/account-management-api/internal/core/ports"
)

type accountService struct {
	repo ports.AccountRepository
}

func NewAccountService(r ports.AccountRepository) ports.AccountService {
	return &accountService{repo: r}
}

func (s *accountService) CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	hasAccount, _ := s.repo.FindByDocument(ctx, documentNumber)
	if hasAccount != nil {
		return nil, domain.ErrAccountAlreadyExists
	}

	acc := &domain.Account{DocumentNumber: documentNumber}
	return s.repo.Save(ctx, acc)
}
