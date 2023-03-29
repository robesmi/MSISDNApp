package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/robesmi/MSISDNApp/mocks/repository"
	"github.com/robesmi/MSISDNApp/model"
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
	authService = ReturnAuthService(mockUserRepo, nil)

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

func TestGetAllCountries(t *testing.T) {

	// Arrange
	teardown := setup(t)
	defer teardown()
	
	countries := []model.Country{
		{
			CountryNumberFormat: "1",
			CountryCode: "2",
			CountryIdentifier: "tt",
			CountryCodeLength: 1,
		},
		{
			CountryNumberFormat: "2",
			CountryCode: "3",
			CountryIdentifier: "te",
			CountryCodeLength: 1,
		},
	}
	mockMSISDNRepo.EXPECT().GetAllCountries().Return(&countries,nil)

	// Act
	res, err := lookupService.GetAllCountries()

	//Assert
	if err != nil{
		t.Errorf("Failed in TestGetAllCountries:\n expected = %s\n got = %s", "nil", err)
	}
	for k,v := range *res{
		if v != countries[k]{
			t.Errorf("Error in TestGetAllCountries result mismatch:\n expected = %s\n got = %s", countries[k].CountryIdentifier, v.CountryIdentifier)
		}
	}

}

func TestGetAllOperators(t *testing.T) {

	// Arrange
	teardown := setup(t)
	defer teardown()
	
	operators := []model.MobileOperator{
		{
			CountryIdentifier: "tt",
			PrefixFormat: "21",
			MNO: "test",
			PrefixLength: 2,
		},
		{
			CountryIdentifier: "ta",
			PrefixFormat: "33",
			MNO: "test1",
			PrefixLength: 3,
		},
	}
	mockMSISDNRepo.EXPECT().GetAllMobileOperators().Return(&operators,nil)

	// Act
	res, err := lookupService.GetAllMobileOperators()

	//Assert
	if err != nil{
		t.Errorf("Failed in TestGetAllOperators:\n expected = %s\n got = %s", "nil", err)
	}
	for k,v := range *res{
		if v != operators[k]{
			t.Errorf("Error in TestGetAllOperators result mismatch:\n expected = %s\n got = %s", operators[k].MNO, v.MNO)
		}
	}

}

func TestAddNewCountryValid(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()

	counReq := dto.CountryRequest{
		CountryNumberFormat: "test",
		CountryCode: "t1",
		CountryIdentifier: "tt1",
		CountryCodeLength: "2",
	}

	mockMSISDNRepo.EXPECT().AddNewCountry(counReq.CountryNumberFormat, counReq.CountryCode, counReq.CountryIdentifier, gomock.AssignableToTypeOf(2)).Return(nil)


	//Act
	err := lookupService.AddNewCountry(&counReq)

	//Arrange

	if err != nil{
		t.Errorf("Error in TestAddNewCountryValid:\n expected = %s\n got = %s", "nil", err)
	}

}

func TestAddNewOperatorValid(t *testing.T) {

	//Arrange
	teardown := setup(t)
	defer teardown()

	mobileReq := dto.OperatorRequest{
		CountryIdentifier: "test1",
		PrefixFormat: "tt1",
		MNO: "t1",
		PrefixLength: "2",
	}

	mockMSISDNRepo.EXPECT().AddNewMobileOperator(mobileReq.CountryIdentifier, mobileReq.PrefixFormat, mobileReq.MNO, gomock.AssignableToTypeOf(2))


	//Act
	err := lookupService.AddNewMobileOperator(&mobileReq)

	//Arrange

	if err != nil{
		t.Errorf("Error in TestAddNewOperatorValid:\n expected = %s\n got = %s", "nil", err)
	}
	
}