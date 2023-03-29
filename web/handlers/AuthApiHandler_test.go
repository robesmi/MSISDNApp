package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/vault"
)

func TestNativeRegisterCallCorrectInput(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@goodmail.com",
		Password: "12345Aa!",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/register", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	var resp = dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	mockAuthService.EXPECT().RegisterNativeUser(jsonReq.Username, jsonReq.Password,"user").Return(&resp,nil)

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusOK{
		t.Errorf("Error in TestNativeRegisterCallCorrectInput:\n expected = %d\n got = %d", http.StatusOK, recorder.Code)
	}
	co1, co2 := recorder.Result().Cookies()[0].Value, recorder.Result().Cookies()[1].Value
	if co1 != resp.AccessToken{
		t.Errorf("Error in TestNativeRegisterCallCorrectInput:\n expected = %s\n got = %s", resp.AccessToken,co1)
	}
	if co2 != resp.RefreshToken{
		t.Errorf("Error in TestNativeRegisterCallCorrectInput:\n expected = %s\n got = %s", resp.RefreshToken,co2)
	}
}

func TestNativeRegisterCallBadEmail(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@reallybadmail",
		Password: "12345Aa!",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/register", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Errorf("Error in TestNativeRegisterCallBadEmail:\n expected = %d\n got = %d", http.StatusBadRequest, recorder.Code)
	}

}


func TestNativeRegisterCallBadPassword(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@goodmail.com",
		Password: "lol",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/register", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Errorf("Error in TestNativeRegisterCallBadPassword:\n expected = %d\n got = %d", http.StatusBadRequest, recorder.Code)
	}

}

func TestNativeRegisterCallDuplicateEmail(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@goodmail.com",
		Password: "12345Aa!",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/register", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	err := errs.NewUserAlreadyExistsError()
	mockAuthService.EXPECT().RegisterNativeUser(jsonReq.Username, jsonReq.Password, "user").Return(nil,err)

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Errorf("Error in TestNativeRegisterCallDuplicateEmail:\n expected = %d\n got = %d", http.StatusBadRequest, recorder.Code)
	}

}

func TestNativeLoginCallCorrectInput(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@goodmail.com",
		Password: "12345Aa!",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/login", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	var resp = dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	mockAuthService.EXPECT().LoginNativeUser(jsonReq.Username, jsonReq.Password).Return(&resp,nil)

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusOK{
		t.Errorf("Error in TestNativeLoginCallCorrectInput:\n expected = %d\n got = %d", http.StatusOK, recorder.Code)
	}
	co1, co2 := recorder.Result().Cookies()[0].Value, recorder.Result().Cookies()[1].Value
	if co1 != resp.AccessToken{
		t.Errorf("Error in TestNativeLoginCallCorrectInput:\n expected = %s\n got = %s", resp.AccessToken,co1)
	}
	if co2 != resp.RefreshToken{
		t.Errorf("Error in TestNativeLoginCallCorrectInput:\n expected = %s\n got = %s", resp.RefreshToken,co2)
	}
}

func TestNativeLoginCallBadEmail(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@reallybadmail",
		Password: "12345Aa!",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/login", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Errorf("Error in TestNativeLoginCallBadEmail:\n expected = %d\n got = %d", http.StatusBadRequest, recorder.Code)
	}

}

func TestNativeLoginCallBadPassword(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@goodmail.com",
		Password: "lol",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/login", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Errorf("Error in TestNativeLoginCallBadPassword:\n expected = %d\n got = %d", http.StatusBadRequest, recorder.Code)
	}

}

func TestNativeLoginCallInvalidCredentials(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := LoginForm{
		Username: "test@goodmail.com",
		Password: "12345Aa!",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/login", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	err := errs.NewInvalidCredentialsError()
	mockAuthService.EXPECT().LoginNativeUser(jsonReq.Username, jsonReq.Password).Return(nil,err)

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusBadRequest{
		t.Errorf("Error in TestNativeLoginCallInvalidCredentials:\n expected = %d\n got = %d", http.StatusOK, recorder.Code)
	}

}

func TestRefreshAccessTokenCall(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := RefreshRequest{
		RefreshToken: "",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	var resp = dto.LoginResponse{
		AccessToken: "test1",
		RefreshToken: "test2",
	}
	validateAccessToken = func(vault *vault.Vault,s string)(jwt.MapClaims, error){
		return jwt.MapClaims{
			"id":"testid",
		},nil
	}

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/refresh", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	mockAuthService.EXPECT().RefreshTokens("testid",jsonReq.RefreshToken).Return(&resp,nil)

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusOK{
		t.Errorf("Error in TestRefreshAccessTokenCall:\n expected = %d\n got = %d", http.StatusOK, recorder.Code)
	}
	co1, co2 := recorder.Result().Cookies()[0].Value, recorder.Result().Cookies()[1].Value
	if co1 != resp.AccessToken{
		t.Errorf("Error in TestRefreshAccessTokenCall:\n expected = %s\n got = %s", resp.AccessToken,co1)
	}
	if co2 != resp.RefreshToken{
		t.Errorf("Error in TestRefreshAccessTokenCall:\n expected = %s\n got = %s", resp.RefreshToken,co2)
	}

}

func TestLogOutCall(t *testing.T) {

	//Arrange
	recorder := httptest.NewRecorder()
	teardown := setup(t,recorder)
	defer teardown()

	jsonReq := RefreshRequest{
		RefreshToken: "",
	}
	jsonVal, _ := json.Marshal(jsonReq)

	
	validateRefreshToken = func(vault *vault.Vault,s string)(jwt.MapClaims, error){
		return jwt.MapClaims{
			"id":"testid",
		},nil
	}

	//Act
	req := httptest.NewRequest(http.MethodPost, "/service/api/logout", bytes.NewBuffer(jsonVal))
	req.Header.Set("Content-Type", "application/json")

	mockAuthService.EXPECT().LogOutUser("testid").Return(nil)

	router.ServeHTTP(recorder,req)
	
	//Assert
	if recorder.Code != http.StatusOK{
		t.Errorf("Error in TestLogOutCall:\n expected = %d\n got = %d", http.StatusOK, recorder.Code)
	}
	co1, co2 := recorder.Result().Cookies()[0].Value, recorder.Result().Cookies()[1].Value
	if co1 != ""{
		t.Errorf("Error in TestLogOutCall:\n expected = %s\n got = %s", "",co1)
	}
	if co2 != ""{
		t.Errorf("Error in TestLogOutCall:\n expected = %s\n got = %s", "",co2)
	}

}