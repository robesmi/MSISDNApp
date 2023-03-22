package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robesmi/MSISDNApp/model"
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
	AddNewCountry(*dto.CountryRequest) (error)
	AddNewMobileOperator(*dto.OperatorRequest) (error)
	GetAllCountries() (*[]model.Country, error)
	GetAllMobileOperators() (*[]model.MobileOperator, error)
	RemoveCountry(string) (error)
	RemoveOperator(string) (error)
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

func (s DefaultMSISDNService) GetAllCountries() (*[]model.Country, error){

	resp, err := s.repo.GetAllCountries()
	if err != nil{
		return nil, err
	}
	return resp, nil
}

func (s DefaultMSISDNService) GetAllMobileOperators() (*[]model.MobileOperator, error){

	resp, err := s.repo.GetAllMobileOperators()
	if err != nil {
		return nil, err
	}
	return resp, nil
}


func (s DefaultMSISDNService) AddNewCountry( counReq *dto.CountryRequest) (error){

	codeLength, err := strconv.Atoi(counReq.CountryCodeLength)
	if err != nil{
		return err
	}
	res := s.repo.AddNewCountry(counReq.CountryNumberFormat, counReq.CountryCode, strings.ToLower(counReq.CountryIdentifier), codeLength)
	if res != nil{
		return res
	}
	return nil
}

func (s DefaultMSISDNService) AddNewMobileOperator( mobileReq *dto.OperatorRequest) (error){

	prefLength, err := strconv.Atoi(mobileReq.PrefixLength)
	if err != nil{
		return err
	}
	res := s.repo.AddNewMobileOperator(strings.ToLower(mobileReq.CountryIdentifier), mobileReq.PrefixFormat, mobileReq.MNO, prefLength)
	if res != nil{
		return res
	}
	return nil
}

func (s DefaultMSISDNService) RemoveCountry(prefix string) (error){

	err := s.repo.RemoveCountry(prefix)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultMSISDNService) RemoveOperator(prefix string) (error){

	err := s.repo.RemoveOperator(prefix)
	if err != nil {
		return err
	}
	return nil
}