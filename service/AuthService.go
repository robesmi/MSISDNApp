package service

import (
	"log"

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

func ReturnAuthService() AuthService {
	return DefaultAuthService{}
}

type AuthService interface {
	// RegisterNativeUser adds a new user to the user database using the conventional user+password combination
	RegisterNativeUser(string, string) (*dto.LoginResponse, *errs.AppError)
	// Login native user searches users via username and password combination, generates new access and refresh tokens
	LoginNativeUser(string,string) (*dto.LoginResponse, *errs.AppError)
	RegisterImportedUser(string) (*dto.LoginResponse, *errs.AppError)
	LoginImportedUser(string) (*dto.LoginResponse, *errs.AppError)

}


func (s DefaultAuthService) RegisterNativeUser(username string, password string) (*dto.LoginResponse, *errs.AppError){
	
	//Checks if user exists and returns error if so
	_ , err := s.repository.GetUserByUsername(username)
	if err != errs.UserNotFound() {
		return nil, errs.UserAlreadyExists()
	}

	//Creates new id, encrypted password and tokens and registers to db
	newID := uuid.NewString()
	
	accessToken , errrr := utils.CreateAccessToken("user")
	if errrr != nil{
		return nil, errrr
	}
	refreshToken, err := utils.CreateRefreshToken(newID)
	if err != nil {
		return nil, err
	}
	encodedPassword, erro := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if erro != nil{
		return nil, errs.UnexpectedError(erro.Error())
	}
	errr := s.repository.RegisterNativeUser(newID, username, string(encodedPassword), "user",refreshToken)
	if errr != nil {
		return nil, errr
	}

	// If successful, returns the tokens
	var response = dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}


func (s DefaultAuthService) LoginNativeUser(username string, password string) (*dto.LoginResponse, *errs.AppError){
	
	// Password check
	user, lookupErr := s.repository.GetUserByUsername(username)
	if lookupErr != nil {
		return nil, lookupErr
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil{
		log.Println(err)
	}
	isPasswordInvalid := bcrypt.CompareHashAndPassword(encryptedPassword,[]byte(user.Password))
	if isPasswordInvalid != nil{
		return nil, errs.InvalidCredentials()
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

func (s DefaultAuthService)RegisterImportedUser(username string) (*dto.LoginResponse, *errs.AppError){
	
	_ , err := s.repository.GetUserByUsername(username)
	if err != errs.UserNotFound() {
		return nil, errs.UserAlreadyExists()
	}

	//Creates new id, encrypted password and tokens and registers to db
	newID := uuid.NewString()
	
	accessToken , errrr := utils.CreateAccessToken("user")
	if errrr != nil{
		return nil, errrr
	}
	refreshToken, err := utils.CreateRefreshToken(newID)
	if err != nil {
		return nil, err
	}
	errr := s.repository.RegisterImportedUser(newID, username, "user",refreshToken)
	if errr != nil {
		return nil, errr
	}

	// If successful, returns the login response
	var response = dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}

func (s DefaultAuthService)LoginImportedUser(username string) (*dto.LoginResponse, *errs.AppError){
	// Password check
	user, lookupErr := s.repository.GetUserByUsername(username)
	if lookupErr != nil {
		return nil, lookupErr
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

func (s DefaultAuthService) LogOutUser(id string) (*errs.AppError){
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