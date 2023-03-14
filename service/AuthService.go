package service

import (

	"github.com/google/uuid"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/repository"
	"github.com/robesmi/MSISDNApp/utils"
	"golang.org/x/crypto/bcrypt"
)

type DefaultAuthService struct {
	repository repository.UserRepository
}

func ReturnAuthService(repository repository.UserRepository) AuthService {
	return DefaultAuthService{repository: repository}
}

type AuthService interface {
	// RegisterNativeUser adds a new user to the user database using the conventional user+password combination
	RegisterNativeUser(string, string) (*dto.LoginResponse, error)
	// LoginNativeUser searches a user and confirms valid credentials, generates new access and refresh tokens
	// and returns them
	LoginNativeUser(string,string) (*dto.LoginResponse, error)
	// RegisterImportedUser adds a new user to the user database using the email received from the Identity Provider
	RegisterImportedUser(string) (*dto.LoginResponse, error)
	// LoginImportedUser searches a user via the received email from the Identity Provider, generates access and refresh
	// tokens and returns them
	LoginImportedUser(string) (*dto.LoginResponse, error)
	// RefreshTokens takes a uuid and a refresh token, checking if the token matches the database one and returning
	// a new pair of tokens if successful, or an error otherwise
	RefreshTokens(string, string) (*dto.LoginResponse, error)
	// LogOutUser finds a user via uuid and removes their refresh token in the database
	LogOutUser(string) (error)

}


func (s DefaultAuthService) RegisterNativeUser(username string, password string) (*dto.LoginResponse, error){
	
	if username == "" || password == ""{
		return nil, errs.NewInvalidCredentialsError()
	}

	_ , err := s.repository.GetUserByUsername(username)
	if _,ok := err.(*errs.UserNotFoundError); !ok {
		return nil, errs.NewUserAlreadyExistsError()
	}

	newID := uuid.NewString()
	
	accessToken , atErr := utils.CreateAccessToken("user")
	if atErr != nil{
		return nil, atErr
	}
	refreshToken, rtErr := utils.CreateRefreshToken(newID)
	if rtErr != nil {
		return nil, rtErr
	}
	encodedPassword, genErr := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if genErr != nil{
		return nil, errs.NewUnexpectedError(genErr.Error())
	}
	regErr := s.repository.RegisterNativeUser(newID, username, string(encodedPassword), "user",refreshToken)
	if regErr != nil {
		return nil, errs.NewUnexpectedError(regErr.Error())
	}

	// If successful, returns the tokens
	var response = dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}


func (s DefaultAuthService) LoginNativeUser(username string, password string) (*dto.LoginResponse, error){
	
	if username == "" || password == ""{
		return nil, errs.NewInvalidCredentialsError()
	}

	user, lookupErr := s.repository.GetUserByUsername(username)
	if lookupErr != nil {
		return nil, nil
	}
	isPasswordInvalid := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
	if isPasswordInvalid != nil{
		return nil, errs.NewInvalidCredentialsError()
	}

	// Create new tokens and update the refresh token in db
	accessToken , atErr := utils.CreateAccessToken("user")
	if atErr != nil{
		return nil, atErr
	}
	refreshToken, rtErr := utils.CreateRefreshToken(user.UUID)
	if rtErr != nil {
		return nil, rtErr
	}
	updateErr := s.repository.UpdateRefreshToken(user.UUID, refreshToken)
	if updateErr != nil{
		return nil, updateErr
	}

	var response = dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}

func (s DefaultAuthService)RegisterImportedUser(username string) (*dto.LoginResponse, error){
	
	if username == ""{
		return nil, errs.NewInvalidCredentialsError()
	}
	_ , err := s.repository.GetUserByUsername(username)
	if _, ok := err.(*errs.UserNotFoundError); ok {
		
		newID := uuid.NewString()
		
		accessToken , atErr := utils.CreateAccessToken("user")
		if atErr != nil{
			return nil, atErr
		}
		refreshToken, rtErr := utils.CreateRefreshToken(newID)
		if rtErr != nil {
			return nil, rtErr
		}
		errr := s.repository.RegisterImportedUser(newID, username, "user",refreshToken)
		if errr != nil {
			return nil, errr
		}

		var response = dto.LoginResponse{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
		}

		return &response, nil
	}else{
		return nil, errs.NewUserAlreadyExistsError()
	}
}

func (s DefaultAuthService)LoginImportedUser(username string) (*dto.LoginResponse, error){

	if username == ""{
		return nil, errs.NewInvalidCredentialsError()
	}

	user, lookupErr := s.repository.GetUserByUsername(username)
	if lookupErr != nil {
		return nil, nil
	}

	// Create new tokens and update the refresh token in db
	accessToken , atErr := utils.CreateAccessToken("user")
	if atErr != nil{
		return nil, atErr
	}
	refreshToken, rtErr := utils.CreateRefreshToken(user.UUID)
	if rtErr != nil {
		return nil, rtErr
	}
	updateErr := s.repository.UpdateRefreshToken(user.UUID, refreshToken)
	if updateErr != nil{
		return nil, updateErr
	}

	var response = dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}

func (s DefaultAuthService)RefreshTokens(id string, token string) (*dto.LoginResponse, error){
	user, err := s.repository.GetUserById(id)
	if err != nil{
		return nil, err
	}
	if user.RefreshToken != token{
		return nil, errs.NewRefreshTokenMismatch()
	}
	accessToken , atErr := utils.CreateAccessToken(user.Role)
	if atErr != nil{
		return nil, atErr
	}
	refreshToken, rtErr := utils.CreateRefreshToken(user.UUID)
	if rtErr != nil {
		return nil, rtErr
	}

	if refErr := s.repository.UpdateRefreshToken(id, refreshToken); refErr != nil{
		return nil,err
	}

	var response = dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
} 

func (s DefaultAuthService) LogOutUser(id string) (error){
	user, lookupErr := s.repository.GetUserById(id)
	if lookupErr != nil{
		return lookupErr
	}
	updateErr := s.repository.UpdateRefreshToken(user.UUID,"")
	if updateErr != nil{
		return updateErr
	}
	return nil

}