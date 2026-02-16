package domain

import (
	domainerror "github.com/evythrossell/account-management-api/internal/core/error"
)

type Account struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

func (a *Account) Validate() error {
	if a.DocumentNumber == "" {
		return domainerror.NewValidationError("document number is required", nil)
	}

	if len(a.DocumentNumber) < 11 || len(a.DocumentNumber) > 14 {
		return domainerror.NewValidationError("document number must be between 11 and 14 characters", nil)
	}

	return nil
}
