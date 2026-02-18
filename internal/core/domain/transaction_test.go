package domain_test

import (
	"testing"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	common "github.com/evythrossell/account-management-api/pkg"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name           string
		accountID      int64
		opType         domain.OperationType
		amount         float64
		expectedAmount float64
		expectedErr    error
	}{
		{
			name:           "Success - Purchase should be negative",
			accountID:      1,
			opType:         domain.Purchase,
			amount:         100.50,
			expectedAmount: -100.50,
			expectedErr:    nil,
		},
		{
			name:           "Success - Payment should be positive",
			accountID:      1,
			opType:         domain.Payment,
			amount:         50.00,
			expectedAmount: 50.00,
			expectedErr:    nil,
		},
		{
			name:           "Success - Withdrawal should be negative",
			accountID:      1,
			opType:         domain.Withdrawal,
			amount:         20.00,
			expectedAmount: -20.00,
			expectedErr:    nil,
		},
		{
			name:        "Error - Amount zero",
			accountID:   1,
			opType:      domain.Purchase,
			amount:      0,
			expectedErr: common.ErrInvalidAmount,
		},
		{
			name:        "Error - Amount negative",
			accountID:   1,
			opType:      domain.Purchase,
			amount:      -10.00,
			expectedErr: common.ErrInvalidAmount,
		},
		{
			name:        "Error - Invalid Operation Type",
			accountID:   1,
			opType:      domain.OperationType(99),
			amount:      100.00,
			expectedErr: common.ErrInvalidOperation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := domain.NewTransaction(tt.accountID, tt.opType, tt.amount)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, tx)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tx)
				assert.Equal(t, tt.accountID, tx.AccountID)
				assert.Equal(t, tt.opType, tx.OperationTypeID)
				assert.Equal(t, tt.expectedAmount, tx.Amount)
				assert.NotZero(t, tx.EventDate)
			}
		})
	}
}
