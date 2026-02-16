package domain

import "errors"

var ErrAccountAlreadyExists = errors.New("Document number already registered.")

type Account struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
