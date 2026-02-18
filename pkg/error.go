package pkg

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrTransactionNotFound = errors.New("transaction not found")

	ErrInvalidDocument      = errors.New("invalid document format")
	ErrAccountAlreadyExists = errors.New("account with this document already exists")

	ErrInvalidAmount    = errors.New("amount must be greater than zero")
	ErrInvalidOperation = errors.New("invalid operation type for transaction")
)

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func (e *DomainError) HTTPStatusCode() int {
	switch e.Code {
	case "VALIDATION_ERROR":
		return http.StatusBadRequest
	case "CONFLICT_ERROR":
		return http.StatusConflict
	case "NOT_FOUND_ERROR":
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func (e *DomainError) PublicMessage() string {
	if e.Message != "" {
		return e.Message
	}
	return "An unexpected error occurred"
}

var (
	ErrValidation = &DomainError{
		Code:    "VALIDATION_ERROR",
		Message: "Invalid input provided",
	}

	ErrConflict = &DomainError{
		Code:    "CONFLICT_ERROR",
		Message: "Resource conflict",
	}

	ErrNotFound = &DomainError{
		Code:    "NOT_FOUND_ERROR",
		Message: "Resource not found",
	}

	ErrInternal = &DomainError{
		Code:    "INTERNAL_ERROR",
		Message: "Internal server error",
	}
)

func NewValidationError(msg string, err error) *DomainError {
	return &DomainError{
		Code:    ErrValidation.Code,
		Message: msg,
		Err:     err,
	}
}

func NewConflictError(msg string, err error) *DomainError {
	return &DomainError{
		Code:    ErrConflict.Code,
		Message: msg,
		Err:     err,
	}
}

func NewNotFoundError(msg string, err error) *DomainError {
	return &DomainError{
		Code:    ErrNotFound.Code,
		Message: msg,
		Err:     err,
	}
}

func NewInternalError(msg string, err error) *DomainError {
	return &DomainError{
		Code:    ErrInternal.Code,
		Message: msg,
		Err:     err,
	}
}

func Is(err error, target *DomainError) bool {
	var de *DomainError
	return errors.As(err, &de) && de.Code == target.Code
}

func As(err error, target *DomainError) bool {
	var de *DomainError
	return errors.As(err, &de)
}
