package services

import (
	"context"
	"errors"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/evythrossell/account-management-api/internal/core/port"
	common "github.com/evythrossell/account-management-api/pkg"
)

type accountService struct {
	repo port.AccountRepository
}

func NewAccountService(repo port.AccountRepository) port.AccountService {
	return &accountService{repo: repo}
}

func (service *accountService) CreateAccount(ctx context.Context, docNumber string) (*domain.Account, error) {
	acc, err := domain.NewAccount(docNumber)
	if err != nil {
		return nil, common.NewValidationError(domain.ErrMsgDocumentInvalid, err)
	}

	savedAcc, err := service.repo.Save(ctx, acc)
	if err != nil {
		if errors.Is(err, common.ErrAccountAlreadyExists) {
			return nil, common.NewConflictError(domain.ErrMsgAccountExists, err)
		}
		return nil, common.NewInternalError(domain.ErrMsgSaveAccountFailed, err)
	}

	return savedAcc, nil
}

func (service *accountService) GetAccountByDocument(ctx context.Context, documentNumber string) (*domain.Account, error) {
	acc, err := service.repo.FindByDocument(ctx, documentNumber)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *accountService) GetAccountByID(ctx context.Context, id int64) (*domain.Account, error) {
	acc, err := s.repo.FindByAccountID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrAccountNotFound) {
			return nil, common.NewNotFoundError(domain.ErrMsgAccountNotFound, err)
		}
		return nil, common.NewInternalError(domain.ErrMsgDatabaseError, err)
	}
	return acc, nil
}
