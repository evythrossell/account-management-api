package errors

import (
	"errors"
	"fmt"
	"net/http"
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
	case ErrValidation.Code:
		return http.StatusBadRequest
	case ErrConflict.Code:
		return http.StatusConflict
	case ErrNotFound.Code:
		return http.StatusNotFound
	case ErrInternal.Code:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (e *DomainError) PublicMessage() string {
	switch e.Code {
	case ErrValidation.Code:
		return "Invalid request format"
	case ErrConflict.Code:
		return "Document number already registered"
	case ErrNotFound.Code:
		return "Resource not found"
	case ErrInternal.Code:
		return "Internal server error"
	default:
		return "An unexpected error occurred"
	}
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
