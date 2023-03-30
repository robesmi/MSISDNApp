package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/robesmi/MSISDNApp/model"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/vault"
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
	createAccessToken = func(role string, vault vault.VaultInterface) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string, vault vault.VaultInterface) (string, error) {
		return expResponse.RefreshToken, nil
	}


	encryptEmailAes256 = func(key []byte, message string) (string, error) {
		return inputEmail, nil
	}
	test := map[string]string{"":""}

	neErr := errs.NewUserNotFoundError()
	
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(nil,neErr)
	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)
	mockUserRepo.EXPECT().RegisterNativeUser(gomock.Any(),gomock.Any(),gomock.Any(),gomock.Any(),gomock.Any())

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
	testUser := model.User{
		UUID: "1",
	}

	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(&testUser,expErr)

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
	encodedPassword, _ := bcrypt.GenerateFromPassword([]byte(inputPassword),bcrypt.DefaultCost)


	expResponse := dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	createAccessToken = func(role string, vault vault.VaultInterface) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string, vault vault.VaultInterface) (string, error) {
		return expResponse.RefreshToken, nil
	}
	respUser := model.User{
		UUID: "testid",
		Username: inputEmail,
		Password: string(encodedPassword),
	}
	encryptEmailAes256 = func(key []byte, message string) (string, error) {
		return inputEmail, nil
	}

	test := map[string]string{"":""}

	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(&respUser,nil)
	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)
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
	createAccessToken = func(role string, vault vault.VaultInterface) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string, vault vault.VaultInterface) (string, error) {
		return expResponse.RefreshToken, nil
	}
	encryptEmailAes256 = func(key []byte, message string) (string, error) {
		return inputEmail, nil
	}

	test := map[string]string{"":""}


	neErr := errs.NewUserNotFoundError()
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(nil,neErr)
	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)
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
	testUser := model.User{
		UUID: "1",
	}
	encryptEmailAes256 = func(key []byte, message string) (string, error) {
		return inputEmail, nil
	}
	test := map[string]string{"":""}


	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(&testUser,nil)

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
	createAccessToken = func(role string, vault vault.VaultInterface) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string, vault vault.VaultInterface) (string, error) {
		return expResponse.RefreshToken, nil
	}
	respUser := model.User{
		UUID: "testid",
		Username: inputEmail,
	}
	encryptEmailAes256 = func(key []byte, message string) (string, error) {
		return inputEmail, nil
	}
	test := map[string]string{"":""}

	
	mockUserRepo.EXPECT().GetUserByUsername(inputEmail).Return(&respUser,nil)
	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)
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
	createAccessToken = func(role string, vault vault.VaultInterface) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string, vault vault.VaultInterface) (string, error) {
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
	createAccessToken = func(role string, vault vault.VaultInterface) (string,error) {
		return expResponse.AccessToken, nil
	}
	createRefreshToken = func(userid string, vault vault.VaultInterface) (string, error) {
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

func TestGetAllUsers(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()

	users := []model.User{
		{
			UUID: "1",
			Username: "test1",
			Password: "pw1",
			Role: "user",
		},
		{
			UUID: "2",
			Username: "test2",
			Password:  "pw2",
			Role: "admin",
		},
	}

	mockUserRepo.EXPECT().GetAllUsers().Return(&users,nil)

	//Act
	usersResponse, err := authService.GetAllUsers()
	

	//Assert
	if err != nil{
		t.Errorf("Error in TestGetAllUsers:\n expected = %s\n got = %s", "nil", err)
	}
	for k,v := range *usersResponse{
		if v != users[k]{
			t.Errorf("Error in TestGetAllUsers result mismatch:\n expected = %s\n got = %s",users[k], v)
		}
	}
}

func TestGetUserByIdValid(t *testing.T){

	//Arrange
	teardown := setup(t)
	defer teardown()

	user := model.User{
		UUID: "1",
		Username: "test1",
		Password: "pw1",
		Role: "user",
	}
	decryptEmailAes256 = func(key []byte, ciphertext string) (string, error) {
		return ciphertext, nil
	}
	test := map[string]string{"":""}

	
	mockUserRepo.EXPECT().GetUserById(user.UUID).Return(&user,nil)
	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)
	//Act

	result, err := authService.GetUserById(user.UUID)

	//Assert
	
	if err != nil{
		t.Errorf("Error in TestGetUserbyIdValid:\n expected = %s\n got = %s", "nil", err)
	}
	if result.Username != user.Username{
		t.Errorf("Error in TestGetUserbyIdValid result mismatch\n expected = %s\n got = %s", user.Username, result.Username)
	}
}

func TestGetUserByIdInvalid(t *testing.T){

	//Arrange
	teardown := setup(t)
	defer teardown()

	mockUserRepo.EXPECT().GetUserById("gibberish").Return(nil, errs.NewUserNotFoundError())

	//Act

	_, err := authService.GetUserById("gibberish")

	//Assert
	
	if err != err.(*errs.UserNotFoundError){
		t.Errorf("Error in TestGetUserbyIdInvalid:\n expected = %s\n got = %s", "nil", err)
	}

}

func TestEditUserByIdAndPassword(t *testing.T){

	//Arrange
	teardown := setup(t)
	defer teardown()

	user := model.User{
		UUID: "1",
		Username: "test1",
		Password: "pw1",
		Role: "user",
	}
	encryptEmailAes256 = func(key []byte, message string) (string, error) {
		return message, nil
	}
	test := map[string]string{"":""}

	mockUserRepo.EXPECT().EditUserById(user.UUID, user.Username, gomock.Any() ,user.Role).Return(nil)
	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)
	//Act

	err := authService.EditUserById(user.UUID, user.Username, user.Password, user.Role)

	//Assert
	
	if err != nil{
		t.Errorf("Error in TestEditUserByIdAndPassword:\n expected = %s\n got = %s", "nil", err)
	}

}

func TestEditUserByIdNoPassword(t *testing.T){

	//Arrange
	teardown := setup(t)
	defer teardown()

	user := model.User{
		UUID: "1",
		Username: "test1",
		Password: "pw1",
		Role: "user",
	}
	encryptEmailAes256 = func(key []byte, message string) (string, error) {
		return message, nil
	}
	test := map[string]string{"":""}

	mockUserRepo.EXPECT().EditUserById(user.UUID, user.Username, "",user.Role).Return(nil)
	mockVault.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(test, nil)

	//Act

	err := authService.EditUserById(user.UUID, user.Username, "", user.Role)

	//Assert
	
	if err != nil{
		t.Errorf("Error in TestEditUserByIdNoPassword:\n expected = %s\n got = %s", "nil", err)
	}

}

func TestRemoveUserById(t *testing.T){

	//Arrange
	teardown := setup(t)
	defer teardown()
	id := "1"
	mockUserRepo.EXPECT().RemoveUserById(id).Return(nil)

	//Act
	err := authService.RemoveUserById(id)

	//Assert
	if err != nil{
		t.Errorf("Error in TestRemoveuserById:\n expected = %s\n got = %s", "nil", err)
	}
}