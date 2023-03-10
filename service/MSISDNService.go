package service

import (
	"fmt"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/repository"
)

type DefaultMSISDNService struct {
	repo repository.MSISDNRepository
}

func NewMSISDNService(repository repository.MSISDNRepository) DefaultMSISDNService{
	return DefaultMSISDNService{repository}
}

type MSISDNService interface {
	LookupMSISDN(string) (*dto.NumberLookupResponse, error)
}

// LookupMSISDN takes a full MSISDN as a string and returns
// a response containing the MNO, country code, subscriber number
// and the country identifier in ISO 3166-1-alpha-2 format
// or an error otherwise
//go:generate mockgen -destination=../mocks/service/mockMSISDNService.go -package=service github.com/robesmi/MSISDNApp/service MSISDNService
func (s DefaultMSISDNService) LookupMSISDN(input string) (*dto.NumberLookupResponse, error){
	
	countryResponse, err := s.repo.LookupCountryCode(input)
	if err != nil {
		return nil, err
	}
	significantNumber := fmt.Sprint(input[countryResponse.CountryCodeLength:])
	
	mnoResponse, err := s.repo.LookupMobileOperator(countryResponse.CountryIdentifier, significantNumber)
	if err != nil{
		return nil, err
	}

	subscriberNumber := fmt.Sprint(significantNumber[mnoResponse.PrefixLength:])

	var response = dto.NumberLookupResponse{
		MNO: mnoResponse.MNO,
		SN: subscriberNumber,
		CI: countryResponse.CountryIdentifier,
		CC: countryResponse.CountryCode,
	}
	return &response, nil
}
