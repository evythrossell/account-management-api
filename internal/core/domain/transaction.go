package domain

import (
	"math"
	"time"

	domainerror "github.com/evythrossell/account-management-api/internal/core/domain/error"
)

type Transaction struct {
	ID              int64
	AccountID       int64
	OperationTypeID OperationType
	Amount          float64
	EventDate       time.Time
}

func NewTransaction(accountID int64, opType OperationType, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, domainerror.ErrInvalidAmount
	}

	if !opType.IsValid() {
		return nil, domainerror.ErrInvalidOperation
	}

	normalizedAmount := math.Abs(amount)
	if opType.IsDebt() {
		normalizedAmount = -normalizedAmount
	}

	return &Transaction{
		AccountID:       accountID,
		OperationTypeID: opType,
		Amount:          normalizedAmount,
		EventDate:       time.Now(),
	}, nil
}
