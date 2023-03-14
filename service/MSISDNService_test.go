package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/robesmi/MSISDNApp/mocks/repository"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
)

var mockMSISDNRepo *repository.MockMSISDNRepository
var mockUserRepo *repository.MockUserRepository
var lookupService MSISDNService
var authService AuthService

func setup(t *testing.T) func(){


	ctrl := gomock.NewController(t)
	mockMSISDNRepo = repository.NewMockMSISDNRepository(ctrl)
	mockUserRepo = repository.NewMockUserRepository(ctrl)
	lookupService = NewMSISDNService(mockMSISDNRepo)
	authService = ReturnAuthService(mockUserRepo)

	return func(){
		lookupService = nil
		authService = nil
		ctrl.Finish()
	}
}

func TestNonExistantCountryNumber(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	// Arrange
	input := "6934567890"
	expErr := errs.NewNumberNotFoundError()

	mockMSISDNRepo.EXPECT().LookupCountryCode(input).Return(nil, expErr)

	// Act
	_, err := lookupService.LookupMSISDN(input)

	//Assert
	if !errors.Is(err, expErr){
		t.Error("Failed while testing non existant country number")
	}

}

func TestNonExistantOperatorNumber(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	// Arrange
	input := "38942123456"
	nextInput := "42123456"
	expErr := errs.NewNoCarriersFoundError()

	expCountryResponse := dto.CountryLookupResponse{
		CountryCode: "389",
		CountryIdentifier: "mk",
		CountryCodeLength: 3,
	}

	gomock.InOrder(
		mockMSISDNRepo.EXPECT().LookupCountryCode(input).Return(&expCountryResponse, nil),
		mockMSISDNRepo.EXPECT().LookupMobileOperator(expCountryResponse.CountryIdentifier, nextInput).Return(nil, expErr),
	)

	// Act
	_, err := lookupService.LookupMSISDN(input)

	//Assert
	if !errors.Is(err, expErr){
		t.Error("Failed while testing non existant carrier number")
	}
}

func TestValidNumber(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	// Arrange
	input := "38977123456"
	secondInput := "77123456"

	expCountryResponse := dto.CountryLookupResponse{
		CountryCode: "389",
		CountryIdentifier: "mk",
		CountryCodeLength: 3,
	}
	expMOResponse := dto.MobileOperatorLookupResponse{
		MNO: "A1",
		PrefixLength: 2,
	}
	expFunctionResponse := dto.NumberLookupResponse{
		MNO: "A1",
		CC: "389",
		SN: "123456",
		CI: "mk",
	} 

	gomock.InOrder(
		mockMSISDNRepo.EXPECT().LookupCountryCode(input).Return(&expCountryResponse, nil),
		mockMSISDNRepo.EXPECT().LookupMobileOperator(expCountryResponse.CountryIdentifier, secondInput).Return(&expMOResponse, nil),
	)

	// Act
	response, err := lookupService.LookupMSISDN(input)

	//Assert
	if err != nil || response == nil || !response.Compare(expFunctionResponse){
		t.Error("Failed while testing valid number")
	}
}

