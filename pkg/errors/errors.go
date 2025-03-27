package errors

import (
	"net/http"
)

type CustomError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewError(code string, message string, status int) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Predefined errors
var (
	ErrNotFound          = NewError("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrInvalidInput      = NewError("INVALID_INPUT", "Invalid input provided", http.StatusBadRequest)
	ErrUnauthorized      = NewError("UNAUTHORIZED", "Unauthorized access", http.StatusUnauthorized)
	ErrInternalError     = NewError("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrInsufficientStock = NewError("INSUFFICIENT_STOCK", "Insufficient stock available", http.StatusConflict)
)
