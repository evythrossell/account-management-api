package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
	"github.com/evythrossell/account-management-api/internal/core/ports"
)

type transactionService struct {
	db      *sql.DB
	accRepo ports.AccountRepository
	txRepo  ports.TransactionRepository
	opRepo  ports.OperationRepository
}

func NewTransactionService(
	db *sql.DB, 
	ar ports.AccountRepository, 
	tr ports.TransactionRepository, 
	or ports.OperationRepository,
) ports.TransactionService {
	return &transactionService{
		db: db, 
		accRepo: ar, 
		txRepo: tr, 
		opRepo: or,
	}
}

func (s *transactionService) CreateTransaction(
	ctx context.Context, 
	accountID int64, 
	operationType int16, 
	amount float64,
) (*domain.Transaction, error) {

	if accountID <= 0 {
		return nil, domainerror.NewValidationError("invalid account_id", nil)
	}
	if operationType <= 0 {
		return nil, domainerror.NewValidationError("invalid operation_type_id", nil)
	}
	if amount == 0 {
		return nil, domainerror.NewValidationError("amount cannot be zero", nil)
	}

	acc, err := s.accRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, domainerror.NewInternalError("failed to check account", err)
	}
	if acc == nil {
		return nil, domainerror.NewNotFoundError("account not found", nil)
	}

	ok, err := s.opRepo.Exists(ctx, operationType)
	if err != nil {
		return nil, domainerror.NewInternalError("failed to check operation type", err)
	}
	if !ok {
		return nil, domainerror.NewValidationError("invalid operation_type_id", nil)
	}

	normalizedAmount, err := domain.NormalizeAmount(
		domain.OperationTypeID(operationType),
		amount,
	)
	if err != nil {
		return nil, domainerror.NewValidationError(err.Error(), nil)
	}

	transaction := &domain.Transaction{
		AccountID:       accountID,
		OperationTypeID: operationType,
		Amount:          normalizedAmount,
		EventDate:       time.Now().UTC(),
	}

	saved, err := s.txRepo.Save(ctx, transaction)
	if err != nil {
		return nil, domainerror.NewInternalError("failed to persist transaction", err)
	}
	return saved, nil
}

func (s *transactionService) GetByTransactionID(
	ctx context.Context, 
	transactionID int64,
) (*domain.Transaction, error) {

	if transactionID <= 0 {
		return nil, domainerror.NewValidationError(
			"invalid transaction_id", 
			nil,
		)
	}

	tx, err := s.txRepo.FindByTransactionID(ctx, transactionID)

	if err != nil {
		return nil, domainerror.NewInternalError(
			"failed to fetch transaction", 
			err,
		)
	}

	if tx == nil {
		return nil, domainerror.NewNotFoundError(
			"transaction not found", 
			nil,
		)
	}
	return tx, nil
}
