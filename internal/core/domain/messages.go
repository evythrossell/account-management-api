package domain

const (
	ErrMsgDocumentInvalid         = "document must be between 11 and 14 digits"
	ErrMsgAccountExists           = "account with this document already exists"
	ErrMsgAccountNotFound         = "account not found"
	ErrMsgAccountIDInvalid        = "the account ID must be a valid integer"
	ErrMsgTransactionNotFound     = "transaction not found"
	ErrMsgTransactionIDInvalid    = "the transaction ID must be a valid integer"
	ErrMsgAccountIDDoesNotExist   = "account id does not exist"
	ErrMsgOperationTypeInvalid    = "invalid operation type"
	ErrMsgAmountInvalid           = "amount must be greater than zero"
	ErrMsgDatabaseError           = "database error"
	ErrMsgSaveAccountFailed       = "failed to save account"
	ErrMsgCreateTransactionFailed = "failed to create transaction"

	ErrCodeInvalidBody   = "INVALID_BODY"
	ErrCodeInvalidID     = "INVALID_ID"
	ErrCodeNotFound      = "NOT_FOUND_ERROR"
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeInternalError = "INTERNAL_SERVER_ERROR"
	ErrCodeConflict      = "CONFLICT_ERROR"

	ErrMsgInvalidBodyRequest = "invalid request body or missing required fields"
	ErrMsgUnexpectedError    = "an unexpected error occurred"
)
