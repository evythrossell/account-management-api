package domain

import (
	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
)

type Account struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
