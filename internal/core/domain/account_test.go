package domain_test

import (
	"strings"
	"testing"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name           string
		documentNumber string
		expectedErr    error
		expectNil      bool
	}{
		{
			name:           "Valid CPF (11 digits)",
			documentNumber: "12345678901",
			expectedErr:    nil,
			expectNil:      false,
		},
		{
			name:           "Valid CNPJ (14 digits)",
			documentNumber: "12345678901234",
			expectedErr:    nil,
			expectNil:      false,
		},
		{
			name:           "Document with spaces (should trim)",
			documentNumber: " 12345678901 ",
			expectedErr:    nil,
			expectNil:      false,
		},
		{
			name:           "Invalid length (too short)",
			documentNumber: "123",
			expectedErr:    common.ErrInvalidDocument,
			expectNil:      true,
		},
		{
			name:           "Invalid length (middle size)",
			documentNumber: "123456789012",
			expectedErr:    common.ErrInvalidDocument,
			expectNil:      true,
		},
		{
			name:           "Invalid characters (letters)",
			documentNumber: "1234567890a",
			expectedErr:    common.ErrInvalidDocument,
			expectNil:      true,
		},
		{
			name:           "Empty document",
			documentNumber: "",
			expectedErr:    common.ErrInvalidDocument,
			expectNil:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := domain.NewAccount(tt.documentNumber)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				if tt.expectNil {
					assert.Nil(t, acc)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, acc)

				expectedClean := strings.TrimSpace(tt.documentNumber)
				assert.Equal(t, expectedClean, acc.DocumentNumber)
			}
		})
	}
}
