package errors

import (
	"net/http"
)

type ApplicationError interface {
	Error() string
	Resolve() (code int, message string)
}

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) Resolve() (code int, message string) {
	return e.Code, e.Message
}

func UnauthorizedError(message string) AppError {
	return AppError{Code: http.StatusUnauthorized, Message: message}
}

func InternalServerError(message string) AppError {
	return AppError{Code: http.StatusInternalServerError, Message: message}
}

func ValidationError(message string) AppError {
	return AppError{Code: 400, Message: message}
}
