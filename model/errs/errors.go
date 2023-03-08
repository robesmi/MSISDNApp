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

func NumberNotFoundError(message string) *AppError{
	return &AppError{
		Message: message,
		Code: http.StatusBadRequest,
	}
}

func InvalidInputError(message string) *AppError{
	return &AppError{
		Message: message,
		Code: http.StatusBadRequest,
	}
}
func NoCarriersFound(message string) *AppError{
	return &AppError{
		Message: message,
		Code: http.StatusBadRequest,
	}
}

func UserNotFound() *AppError{
	return &AppError{
		Message: "User not found",
		Code: http.StatusBadRequest,
	}
}

func UserAlreadyExists() *AppError{
	return &AppError{
		Message: "User already exists",
		Code: http.StatusBadRequest,
	}
}

func InvalidCredentials() *AppError{
	return &AppError{
		Message: "Invalid username or password",
		Code: http.StatusBadRequest,
	}
}

func TokenError(message string) *AppError{
	return &AppError{
		Message: "Error with token: " + message,
		Code: http.StatusInternalServerError,
	}
}
func MalformedToken() *AppError{
	return &AppError{
		Message: "Not a token",
		Code: http.StatusInternalServerError,
	}
}
func ExpiredToken() *AppError{
	return &AppError{
		Message: "Token is expired",
		Code: http.StatusInternalServerError,
	}
}

