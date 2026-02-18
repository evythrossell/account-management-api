package services

import (
	"context"
	"errors"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	domainerror "github.com/evythrossell/account-management-api/internal/core/domain/error"
	"github.com/evythrossell/account-management-api/internal/core/ports"
)

type accountService struct {
	repo ports.AccountRepository
}

func NewAccountService(repo ports.AccountRepository) ports.AccountService {
	return &accountService{repo: repo}
}

func (service *accountService) CreateAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	acc, err := domain.NewAccount(documentNumber)
	if err != nil {
		return nil, err
	}

	_, err = service.repo.FindByDocument(ctx, documentNumber)
	if err == nil {
		return nil, domainerror.ErrAccountAlreadyExists
	}

	if !errors.Is(err, domainerror.ErrAccountNotFound) {
		return nil, err
	}

	return service.repo.Save(ctx, acc)
}

func (service *accountService) GetAccount(ctx context.Context, documentNumber string) (*domain.Account, error) {
	acc, err := service.repo.FindByDocument(ctx, documentNumber)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (service *accountService) GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error) {
	acc, err := service.repo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
