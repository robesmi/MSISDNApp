package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/robesmi/MSISDNApp/model"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterNativeUserValidInput(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "verycorrectemail@nice.com"
	inputPassword := "12345Aa!"


	expResponse := dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	createAccessToken = func(role string) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string) (string, error) {
		return expResponse.RefreshToken, nil
	}


	neErr := errs.NewUserNotFoundError()
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(nil,neErr)
	mockUserRepo.EXPECT().RegisterNativeUser(gomock.Any(),inputEmail,gomock.Any(),gomock.Any(),gomock.Any())

	//Act
	resp, err := authService.RegisterNativeUser(inputEmail,inputPassword,"user")


	//Assert
	if err != nil{
		t.Errorf("Error in TestRegisterNativeUserValidInput:\n expected %s\n got = %s", "nil", err.Error())
	}
	if resp.AccessToken != expResponse.AccessToken{
		t.Errorf("Error in TestRegisterNativeUserValidInput:\n expected %s\n got = %s", expResponse.AccessToken, resp.AccessToken)
	}
	if resp.RefreshToken != expResponse.RefreshToken{
		t.Errorf("Error in TestRegisterNativeUserValidInput:\n expected %s\n got = %s", expResponse.RefreshToken, resp.RefreshToken)
	}
}

func TestRegisterNativeUserInvalidEmail(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := ""
	inputPassword := "12345Aa!"
	expErr := errs.NewInvalidCredentialsError()

	//Act
	_, err := authService.RegisterNativeUser(inputEmail,inputPassword,"user")


	//Assert
	if _,ok := err.(*errs.InvalidCredentials); !ok {
		t.Errorf("Error in TestRegisterNativeUserValidInput:\n expected = %s\n got = %s", expErr, err)
	}
}

func TestRegisterNativeUserExistingEmail(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "correctduplicateemail@goodmail.com"
	inputPassword := "12345Aa!"
	expErr := errs.NewUserAlreadyExistsError()

	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(nil,expErr)

	//Act
	_, err := authService.RegisterNativeUser(inputEmail,inputPassword, "user")


	//Assert
	if _,ok := err.(*errs.UserAlreadyExists); !ok {
		t.Errorf("Error in TestRegisterNativeUserValidInput:\n expected = %s\n got = %s", expErr, err)
	}
}

func TestLoginNativeUserCorrectInput(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "verycorrectemail@nice.com"
	inputPassword := "12345Aa!"
	encodedPassword, _ := bcrypt.GenerateFromPassword([]byte(inputPassword),bcrypt.MinCost)


	expResponse := dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	createAccessToken = func(role string) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string) (string, error) {
		return expResponse.RefreshToken, nil
	}
	respUser := model.User{
		UUID: "testid",
		Username: inputEmail,
		Password: string(encodedPassword),
	}

	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(&respUser,nil)
	mockUserRepo.EXPECT().UpdateRefreshToken(respUser.UUID,expResponse.RefreshToken)

	//Act
	resp, err := authService.LoginNativeUser(inputEmail,inputPassword)


	//Assert
	if err != nil{
		t.Errorf("Error in TestLoginNativeUserCorrectInput:\n expected %s\n got = %s", "nil", err.Error())
	}
	if resp.AccessToken != expResponse.AccessToken{
		t.Errorf("Error in TestLoginNativeUserCorrectInput:\n expected %s\n got = %s", expResponse.AccessToken, resp.AccessToken)
	}
	if resp.RefreshToken != expResponse.RefreshToken{
		t.Errorf("Error in TestLoginNativeUserCorrectInput:\n expected %s\n got = %s", expResponse.RefreshToken, resp.RefreshToken)
	}
}

func TestRegisterImportedUserValidInput(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "verycorrectemail@nice.com"


	expResponse := dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	createAccessToken = func(role string) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string) (string, error) {
		return expResponse.RefreshToken, nil
	}


	neErr := errs.NewUserNotFoundError()
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(nil,neErr)
	mockUserRepo.EXPECT().RegisterImportedUser(gomock.Any(),inputEmail,gomock.Any(),gomock.Any())

	//Act
	resp, err := authService.RegisterImportedUser(inputEmail)


	//Assert
	if err != nil{
		t.Errorf("Error in TestRegisterImportedUserValidInput:\n expected %s\n got = %s", "nil", err.Error())
	}
	if resp.AccessToken != expResponse.AccessToken{
		t.Errorf("Error in TestRegisterImportedUserValidInput:\n expected %s\n got = %s", expResponse.AccessToken, resp.AccessToken)
	}
	if resp.RefreshToken != expResponse.RefreshToken{
		t.Errorf("Error in TestRegisterImportedUserValidInput:\n expected %s\n got = %s", expResponse.RefreshToken, resp.RefreshToken)
	}
}

func TestRegisterImportedUserDuplicateUser(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "duplicateemail@nice.com"
	
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(nil,nil)
	//Act
	_, err := authService.RegisterImportedUser(inputEmail)


	//Assert
	if _,ok := err.(*errs.UserAlreadyExists); !ok{
		t.Errorf("Error in TestRegisterImportedUserDuplicateUser:\n expected %s\n got = %s", errs.NewUserAlreadyExistsError(), err.Error())
	}

}

func TestLoginImportedUserValidInput(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "validmail@nice.com"

	expResponse := dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	createAccessToken = func(role string) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string) (string, error) {
		return expResponse.RefreshToken, nil
	}
	respUser := model.User{
		UUID: "testid",
		Username: inputEmail,
	}
	
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(&respUser,nil)
	mockUserRepo.EXPECT().UpdateRefreshToken(respUser.UUID,expResponse.RefreshToken).Return(nil)
	//Act
	resp, err := authService.LoginImportedUser(inputEmail)


	//Assert
	if err != nil{
		t.Errorf("Error in TestLoginImportedUserValidInput:\n expected %s\n got = %s", "nil", err.Error())
	}
	if resp.AccessToken != expResponse.AccessToken{
		t.Errorf("Error in TestLoginImportedUserValidInput:\n expected %s\n got = %s", expResponse.AccessToken, resp.AccessToken)
	}
	if resp.RefreshToken != expResponse.RefreshToken{
		t.Errorf("Error in TestLoginImportedUserValidInput:\n expected %s\n got = %s", expResponse.RefreshToken, resp.RefreshToken)
	}
}

func TestRefreshTokensValid(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "validmail@nice.com"

	expResponse := dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	createAccessToken = func(role string) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string) (string, error) {
		return expResponse.RefreshToken, nil
	}
	respUser := model.User{
		UUID: "testid",
		Username: inputEmail,
		RefreshToken: expResponse.RefreshToken,
	}
	
	mockUserRepo.EXPECT().GetUserById(respUser.UUID).Return(&respUser,nil)
	mockUserRepo.EXPECT().UpdateRefreshToken(respUser.UUID,expResponse.RefreshToken).Return(nil)
	//Act
	resp, err := authService.RefreshTokens(respUser.UUID,respUser.RefreshToken)


	//Assert
	if err != nil{
		t.Errorf("Error in TestRefreshTokensValid:\n expected %s\n got = %s", "nil", err.Error())
	}
	if resp.AccessToken != expResponse.AccessToken{
		t.Errorf("Error in TestRefreshTokensValid:\n expected %s\n got = %s", expResponse.AccessToken, resp.AccessToken)
	}
	if resp.RefreshToken != expResponse.RefreshToken{
		t.Errorf("Error in TestRefreshTokensValid:\n expected %s\n got = %s", expResponse.RefreshToken, resp.RefreshToken)
	}
}

func TestRefreshTokensInvalid(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()
	
	inputEmail := "validmail@nice.com"

	expResponse := dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	createAccessToken = func(role string) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string) (string, error) {
		return expResponse.RefreshToken, nil
	}
	respUser := model.User{
		UUID: "testid",
		Username: inputEmail,
		RefreshToken: expResponse.RefreshToken,
	}
	
	mockUserRepo.EXPECT().GetUserById(respUser.UUID).Return(&respUser,nil)
	//Act
	_, err := authService.RefreshTokens(respUser.UUID,"bad")


	//Assert
	if _,ok := err.(*errs.RefreshTokenMismatch); !ok{
		t.Errorf("Error in TestRefreshTokensInvalid:\n expected %s\n got = %s", errs.NewRefreshTokenMismatch(), err)
	}

}

func TestLogOutUser(t *testing.T){

	//Arrange
	teardown := setup(t)
	defer teardown()

	respUser := model.User{
		UUID: "uuid",
	}

	mockUserRepo.EXPECT().GetUserById(respUser.UUID).Return(&respUser,nil)
	mockUserRepo.EXPECT().UpdateRefreshToken(respUser.UUID,"").Return(nil)

	//Act
	err := authService.LogOutUser(respUser.UUID)


	//Assert
	if err != nil{
		t.Errorf("Error in TestLogOutuser:\n expected = %s\n got = %s", "nil",err)
	}
}