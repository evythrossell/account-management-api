package domain

type OperationType int16

const (
	Purchase            OperationType = 1
	InstallmentPurchase OperationType = 2
	Withdrawal          OperationType = 3
	Payment             OperationType = 4
)

func (op OperationType) IsDebt() bool {
	switch op {
	case Purchase, InstallmentPurchase, Withdrawal:
		return true
	default:
		return false
	}
}

func (op OperationType) IsCredit() bool {
	return op == Payment
}

func (op OperationType) IsValid() bool {
	switch op {
	case Purchase, InstallmentPurchase, Withdrawal, Payment:
		return true
	default:
		return false
	}
}
