package domain

import "time"

type Transaction struct {
	TransactionID   int64     `json:"transaction_id"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int16     `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}
