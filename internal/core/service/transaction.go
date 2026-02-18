package services

import (
	"context"
	"errors"

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
		if errors.Is(err, common.ErrAccountNotFound) {
			return nil, common.NewValidationError(domain.ErrMsgAccountIDDoesNotExist, err)
		}
		return nil, common.NewInternalError(domain.ErrMsgDatabaseError, err)
	}

	opType := domain.OperationType(operationTypeID)
	exists, err := service.opRepo.Exists(ctx, int16(opType))
	if err != nil {
		return nil, common.NewInternalError(domain.ErrMsgDatabaseError, err)
	}
	if !exists {
		return nil, common.NewValidationError(domain.ErrMsgOperationTypeInvalid, common.ErrInvalidOperation)
	}

	tx, err := domain.NewTransaction(accountID, opType, amount)
	if err != nil {
		if errors.Is(err, common.ErrInvalidAmount) {
			return nil, common.NewValidationError(domain.ErrMsgAmountInvalid, err)
		}
		if errors.Is(err, common.ErrInvalidOperation) {
			return nil, common.NewValidationError(domain.ErrMsgOperationTypeInvalid, err)
		}
		return nil, common.NewInternalError(domain.ErrMsgCreateTransactionFailed, err)
	}

	return service.txRepo.Save(ctx, tx)
}

func (service *transactionService) GetByTransactionID(ctx context.Context, transactionID int64) (*domain.Transaction, error) {
	tx, err := service.txRepo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, common.ErrTransactionNotFound) {
			return nil, common.NewNotFoundError(domain.ErrMsgTransactionNotFound, err)
		}
		return nil, common.NewInternalError(domain.ErrMsgDatabaseError, err)
	}

	return tx, nil
}
