package services

import (
	"context"
	"errors"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
	"github.com/evythrossell/account-management-api/internal/core/ports"
)

type accountService struct {
	repo ports.AccountRepository
}

func NewAccountService(r ports.AccountRepository) ports.AccountService {
	return &accountService{repo: r}
}

func (s *accountService) CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	acc := &domain.Account{DocumentNumber: documentNumber}
	if err := acc.Validate(); err != nil {
		return nil, err
	}

	existingAccount, err := s.repo.FindByDocument(ctx, documentNumber)
	if err != nil {
		var de *domainerror.DomainError
		if errors.As(err, &de) {
			return nil, err
		}
		return nil, domainerror.NewInternalError("failed to check account existence", err)
	}

	if existingAccount != nil {
		return nil, domainerror.NewConflictError("document number already registered", nil)
	}

	savedAccount, err := s.repo.Save(ctx, acc)
	if err != nil {
		var de *domainerror.DomainError
		if errors.As(err, &de) {
			return nil, err
		}
		return nil, domainerror.NewInternalError("failed to save account", err)
	}

	return savedAccount, nil
}

func (s *accountService) GetAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	acc, err := s.repo.FindByDocument(ctx, documentNumber)
	if err != nil {
		var de *domainerror.DomainError
		if errors.As(err, &de) {
			return nil, err
		}
		return nil, domainerror.NewInternalError("failed to fetch account", err)
	}

	if acc == nil {
		return nil, domainerror.NewNotFoundError("account not found", nil)
	}

	return acc, nil
}
func (s *accountService) GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error) {
	acc, err := s.repo.FindByAccountID(ctx, accountID)
	if err != nil {
		var de *domainerror.DomainError
		if errors.As(err, &de) {
			return nil, err
		}
		return nil, domainerror.NewInternalError("failed to fetch account", err)
	}

	if acc == nil {
		return nil, domainerror.NewNotFoundError("account not found", nil)
	}

	return acc, nil
}