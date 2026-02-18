package services

import (
	"context"
	"errors"
	"log"

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
		return nil, common.NewValidationError("document must be between 11 and 14 characters", err)
	}

	savedAcc, err := service.repo.Save(ctx, acc)
	if err != nil {
		return nil, common.NewInternalError("failed to save account", err)
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
		log.Printf("[DEBUG] Service - repo.FindByAccountID error: %T - %v", err, err)
		if errors.Is(err, common.ErrAccountNotFound) {
			log.Printf("[DEBUG] Service - error is ErrAccountNotFound, wrapping as NotFoundError")
			return nil, common.NewNotFoundError("account not found", err)
		}
		log.Printf("[DEBUG] Service - error is not ErrAccountNotFound, wrapping as InternalError")
		return nil, common.NewInternalError("database error", err)
	}
	return acc, nil
}
