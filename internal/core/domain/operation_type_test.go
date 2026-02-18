package domain_test

import (
	"testing"

	"github.com/evythrossell/account-management-api/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestOperationType_IsDebt(t *testing.T) {
	tests := []struct {
		name string
		op   domain.OperationType
		want bool
	}{
		{"Purchase is debt", domain.Purchase, true},
		{"InstallmentPurchase is debt", domain.InstallmentPurchase, true},
		{"Withdrawal is debt", domain.Withdrawal, true},
		{"Payment is not debt", domain.Payment, false},
		{"Unknown is not debt", domain.OperationType(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.op.IsDebt())
		})
	}
}

func TestOperationType_IsCredit(t *testing.T) {
	tests := []struct {
		name string
		op   domain.OperationType
		want bool
	}{
		{"Payment is credit", domain.Payment, true},
		{"Purchase is not credit", domain.Purchase, false},
		{"Unknown is not credit", domain.OperationType(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.op.IsCredit())
		})
	}
}

func TestOperationType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		op   domain.OperationType
		want bool
	}{
		{"Purchase is valid", domain.Purchase, true},
		{"InstallmentPurchase is valid", domain.InstallmentPurchase, true},
		{"Withdrawal is valid", domain.Withdrawal, true},
		{"Payment is valid", domain.Payment, true},
		{"Unknown is invalid", domain.OperationType(99), false},
		{"Zero value is invalid", domain.OperationType(0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.op.IsValid())
		})
	}
}
