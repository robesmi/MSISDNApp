package errs

type UserNotFoundError struct {
	Message string
}

func (u UserNotFoundError) Error() string{
	return u.Message
}

func NewUserNotFoundError() *UserNotFoundError{
	return &UserNotFoundError{
		Message: "User not found",
	}
}

type UserAlreadyExists struct {
	Message string
}

func (u UserAlreadyExists) Error() string {
	return u.Message
}


func NewUserAlreadyExistsError() *UserAlreadyExists{
	return &UserAlreadyExists{
		Message: "User already exists",
	}
}

type UnexpectedError struct{
	Message string
}

func (u UnexpectedError) Error() string{
	return u.Message
}

func NewUnexpectedError(err string) *UnexpectedError{
	return &UnexpectedError{
		Message: "Unexpected error " + err,
	}
}

type TokenError struct{
	Message string
}

func(u TokenError) Error() string{
	return u.Message
}

func NewTokenError(err string) *TokenError{
	return &TokenError{
		Message: "Token error" + err,
	}
}

type InvalidCredentials struct{
	Message string
}

func(u InvalidCredentials) Error() string{
	return u.Message
}

func NewInvalidCredentialsError() *InvalidCredentials{
	return &InvalidCredentials{
		Message: "Invalid credentials",
	}
}

type MalformedTokenError struct{
	Message string
}

func(u MalformedTokenError) Error() string{
	return u.Message
}

func NewMalformedTokenError() *MalformedTokenError{
	return &MalformedTokenError{
		Message: "This is not a token",
	}
}

type ExpiredTokenError struct{
	Message string
}

func(u ExpiredTokenError) Error() string{
	return u.Message
}

func NewExpiredTokenError() *ExpiredTokenError{
	return &ExpiredTokenError{
		Message: "Token is expired",
	}
}

type NumberNotFoundError struct{
	Message string
}

func(u NumberNotFoundError) Error() string{
	return u.Message
}

func NewNumberNotFoundError() *NumberNotFoundError{
	return &NumberNotFoundError{
		Message: "Country not found or invalid number entered",
	}
}

type NoCarriersFoundError struct{
	Message string
}

func(u NoCarriersFoundError) Error() string{
	return u.Message
}

func NewNoCarriersFoundError() *NoCarriersFoundError{
	return &NoCarriersFoundError{
		Message: "Carrier not found or invalid number entered",
	}
}

type RefreshTokenMismatch struct{
	Message string
}

func(u RefreshTokenMismatch) Error() string{
	return u.Message
}

func NewRefreshTokenMismatch() *RefreshTokenMismatch{
	return &RefreshTokenMismatch{
		Message: "Please log in again",
	}
}

type TokenValidationError struct{
	Message string
}

func(u TokenValidationError) Error() string{
	return u.Message
}

func NewTokenValidationError(msg string) *TokenValidationError{
	return &TokenValidationError{
		Message: "Token validation error: " + msg,
	}
}

type EncryptionError struct{
	Message string
}

func(u EncryptionError) Error() string{
	return u.Message
}

func NewEncryptionError(msg string) *EncryptionError{
	return &EncryptionError{
		Message: "Encryption error: " + msg,
	}
}
