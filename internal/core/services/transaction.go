package services

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	domainerror "github.com/evythrossell/account-management-api/internal/core/domain/error"
	"github.com/evythrossell/account-management-api/internal/core/ports"
)

type transactionService struct {
	accRepo ports.AccountRepository
	txRepo  ports.TransactionRepository
	opRepo  ports.OperationRepository
}

func NewTransactionService(
	ar ports.AccountRepository,
	tr ports.TransactionRepository,
	or ports.OperationRepository,
) ports.TransactionService {
	return &transactionService{
		accRepo: ar,
		txRepo:  tr,
		opRepo:  or,
	}
}

func (service *transactionService) CreateTransaction(
	ctx context.Context,
	accountID int64,
	operationTypeID int16,
	amount float64,
) (*domain.Transaction, error) {
	_, err := service.accRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	opType := domain.OperationType(operationTypeID)
	exists, err := service.opRepo.Exists(ctx, int16(opType))
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerror.ErrInvalidOperation
	}

	tx, err := domain.NewTransaction(accountID, opType, amount)
	if err != nil {
		return nil, err
	}

	return service.txRepo.Save(ctx, tx)
}

func (service *transactionService) GetByTransactionID(ctx context.Context, transactionID int64) (*domain.Transaction, error) {
	tx, err := service.txRepo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
