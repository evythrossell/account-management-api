package domain

import (
	"strings"
	"unicode/utf8"

	common "github.com/evythrossell/account-management-api/pkg"
)

type Account struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

func NewAccount(docNumber string) (*Account, error) {
	doc := strings.TrimSpace(docNumber)
	length := utf8.RuneCountInString(doc)

	if length < 11 || length > 14 {
		return nil, common.ErrInvalidDocument
	}

	if !isNumeric(doc) {
		return nil, common.ErrInvalidDocument
	}

	return &Account{
		DocumentNumber: doc,
	}, nil
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}