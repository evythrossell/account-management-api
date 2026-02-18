package domain

import (
	"strings"

	"github.com/evythrossell/account-management-api/internal/core/common"
)

type Account struct {
	ID             int64
	DocumentNumber string
}

func NewAccount(documentNumber string) (*Account, error) {
	doc := strings.TrimSpace(documentNumber)

	if !isValidDocument(doc) {
		return nil, common.ErrInvalidDocument
	}

	return &Account{
		DocumentNumber: doc,
	}, nil
}

func isValidDocument(document string) bool {
	length := len(document)

	if length != 11 && length != 14 {
		return false
	}

	for _, char := range document {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
