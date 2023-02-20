package errs

import "net/http"

type AppError struct {
	Message string
	Code    int
}

func (e AppError) AsMessage() *AppError {
	return &AppError{
		Message: e.Message,
	}
}

func UnexpectedError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func InvalidNumberError(message string) *AppError{
	return &AppError{
		Message: message,
		Code: http.StatusBadRequest,
	}
}