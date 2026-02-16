package domain

type OperationTypeID int16

const (
	OperationTypePurchase OperationTypeID = 1
	OperationTypeInstallmentPurchase OperationTypeID = 2
	OperationTypeWithdrawal OperationTypeID = 3
	OperationTypePayment OperationTypeID = 4
)

func (op OperationTypeID) IsDebt() bool {
	switch op {
		case OperationTypePurchase,
			 OperationTypeInstallmentPurchase,
			 OperationTypeWithdrawal:
				return true
		default:
			return false
	}
}

func (op OperationTypeID) IsCredit() bool {
	return op == OperationTypePayment
}