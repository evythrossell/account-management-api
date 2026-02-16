package domain

import (
	"errors"
	"math"
	"time"
)

type Transaction struct {
	TransactionID   int64     `json:"transaction_id"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int16     `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"-"`
}

func NormalizeAmount(operationTypeID OperationTypeID, amount float64) (float64, error) {

	if amount == 0 {
		return 0, errors.New("amount cannot be zero")
	}

	if operationTypeID.IsDebt() {
		return -math.Abs(amount), nil
	}
	if operationTypeID.IsCredit() {
		return math.Abs(amount), nil
	}

	return 0, errors.New("invalid operation type")
}