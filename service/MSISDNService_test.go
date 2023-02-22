package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/robesmi/MSISDNApp/mocks/repository"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
)

var mockRepo *repository.MockMSISDNRepository
var service MSISDNService

func setup(t *testing.T) func(){


	ctrl := gomock.NewController(t)
	mockRepo = repository.NewMockMSISDNRepository(ctrl)
	service = NewMSISDNService(mockRepo)

	return func(){
		service = nil
		ctrl.Finish()
	}
}

func TestNonExistantCountryNumber(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	// Arrange
	input := "6934567890"
	expErr := errs.NumberNotFoundError("Number not found")

	mockRepo.EXPECT().LookupCountryCode(input).Return(nil, expErr)

	// Act
	_, err := service.LookupMSISDN(input)

	//Assert
	if err.Code != expErr.Code{
		t.Error("Failed while testing non existant country number")
	}

}

func TestNonExistantOperatorNumber(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	// Arrange
	input := "38942123456"
	nextInput := "42123456"
	expErr := errs.NoCarriersFound("Invalid operator number")

	expCountryResponse := dto.CountryLookupResponse{
		CountryCode: 389,
		CountryIdentifier: "mk",
		CountryCodeLength: 3,
	}

	gomock.InOrder(
		mockRepo.EXPECT().LookupCountryCode(input).Return(&expCountryResponse, nil),
		mockRepo.EXPECT().LookupMobileOperator(expCountryResponse.CountryCode, nextInput).Return(nil, expErr),
	)

	// Act
	_, err := service.LookupMSISDN(input)

	//Assert
	if err.Code != expErr.Code{
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
		CountryCode: 389,
		CountryIdentifier: "mk",
		CountryCodeLength: 3,
	}
	expMOResponse := dto.MobileOperatorLookupResponse{
		MNO: "A1",
		PrefixLength: 2,
	}
	expFunctionResponse := dto.NumberLookupResponse{
		MNO: "A1",
		CC: 389,
		SN: "123456",
		CI: "mk",
	} 

	gomock.InOrder(
		mockRepo.EXPECT().LookupCountryCode(input).Return(&expCountryResponse, nil),
		mockRepo.EXPECT().LookupMobileOperator(expCountryResponse.CountryIdentifier, secondInput).Return(&expMOResponse, nil),
	)

	// Act
	response, err := service.LookupMSISDN(input)

	//Assert
	if err != nil || response == nil || !response.Compare(expFunctionResponse){
		t.Error("Failed while testing valid number")
	}
}

