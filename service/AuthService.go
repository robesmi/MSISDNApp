package service

import (
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/repository"
)

type DefaultAuthService struct {
	Repository repository.UserRepository
}

func ReturnAuthService() AuthService {
	return DefaultAuthService{}
}

type AuthService interface {
	RegisterNativeUser(string, string)
	LoginNativeUser(string,string) (*dto.LoginResponse, *errs.AppError)
}

// RegisterNativeUser adds a new user to the user database using the conventional user+password combination
//TODO: Add new user with bcrypt password and generate access and refresh tokens
func (s DefaultAuthService) RegisterNativeUser(username string, password string) {
	panic("panic!!")
}

// Login native user searches users via username and password combination, generates new access and refresh token
func (s DefaultAuthService) LoginNativeUser(username string, password string) (*dto.LoginResponse, *errs.AppError){

	panic("panic!!")
}
