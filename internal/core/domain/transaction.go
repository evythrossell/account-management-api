package domain

import (
	"math"
	"time"

	common "github.com/evythrossell/account-management-api/pkg"
)

type Transaction struct {
	ID              int64         `json:"transaction_id"`
	AccountID       int64         `json:"account_id"`
	OperationTypeID OperationType `json:"operation_type_id"`
	Amount          float64       `json:"amount"`
	EventDate       time.Time     `json:"-"`
}

func NewTransaction(accountID int64, opType OperationType, amount float64) (*Transaction, error) {
	if amount <= 0 {
		return nil, common.ErrInvalidAmount
	}

	if !opType.IsValid() {
		return nil, common.ErrInvalidOperation
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
