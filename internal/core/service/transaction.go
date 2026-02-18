package services

import (
	"context"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/evythrossell/account-management-api/internal/core/port"
	common "github.com/evythrossell/account-management-api/pkg"
)

type transactionService struct {
	accRepo port.AccountRepository
	txRepo  port.TransactionRepository
	opRepo  port.OperationRepository
}

func NewTransactionService(
	ar port.AccountRepository,
	tr port.TransactionRepository,
	or port.OperationRepository,
) port.TransactionService {
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
		return nil, common.ErrInvalidOperation
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
