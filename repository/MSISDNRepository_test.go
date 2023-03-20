package repository

import (
	"testing"

	"github.com/robesmi/MSISDNApp/model"
)

var lookupRepo MSISDNRepository

func TestLookupCountryCode(t *testing.T) {

	//Arrange
	mock := setup(t)
	exampleCountry := model.Country{
		CountryNumberFormat: "[0-9]{1}",
		CountryCode: "1",
		CountryIdentifier: "tt",
		CountryCodeLength: 1,
	}
	rows := mock.NewRows([]string{"country_code","country_identifier","country_code_length"}).
	AddRow(exampleCountry.CountryCode, exampleCountry.CountryIdentifier, exampleCountry.CountryCodeLength)
	mock.ExpectQuery("SELECT").WithArgs("1").WillReturnRows(rows)

	//Act
	resp, getErr := lookupRepo.LookupCountryCode("1")

	//Assert

	if getErr != nil{
		t.Errorf("Error in TestLookupCountryCode:\n expected %s\n got %s", "nil", getErr)
	}
	if resp.CountryIdentifier != exampleCountry.CountryIdentifier{
		t.Errorf("Error in TestLookupCountryCode:\n expected %s\n got %s", exampleCountry.CountryIdentifier, resp.CountryIdentifier)
	}
}

func TestLookupMobileOperator(t *testing.T) {

	//Arrange
	mock := setup(t)
	exampleOperator := model.MobileOperator{
		CountryIdentifier: "tt",
		PrefixFormat: "[0-9]{1}",
		MNO: "test",
		PrefixLength: 1,
	}
	rows := mock.NewRows([]string{"mno","prefix_length"}).
	AddRow(exampleOperator.MNO, exampleOperator.PrefixLength)
	mock.ExpectQuery("SELECT").WithArgs("tt","1").WillReturnRows(rows)

	//Act
	resp, getErr := lookupRepo.LookupMobileOperator("tt","1")

	//Assert

	if getErr != nil{
		t.Errorf("Error in TestLookupMobileOperator:\n expected %s\n got %s", "nil", getErr)
	}
	if resp.MNO != exampleOperator.MNO{
		t.Errorf("Error in TestLookupMobileOperator:\n expected %s\n got %s", exampleOperator.MNO, resp.MNO)
	}
}

