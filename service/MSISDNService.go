package service

import (
	"fmt"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/robesmi/MSISDNApp/repository"
)

type DefaultMSISDNService struct {
	repo repository.MSISDNRepository
}

func NewMSISDNService(repository repository.MSISDNRepositoryDb) DefaultMSISDNService{
	return DefaultMSISDNService{repository}
}

type MSISDNService interface {
	LookupMSISDN(string) (*dto.NumberLookupResponse, *errs.AppError)
}

// LookupMSISDN takes a full MSISDN as a string and returns
// a response containing the MNO, country code, subscriber number
// and the country identifier in ISO 3166-1-alpha-2 format
// or an error otherwise
//go:generate mockgen -destination=../mocks/service/mockMSISDNService.go -package=service github.com/robesmi/MSISDNApp/service MSISDNService
func (s DefaultMSISDNService) LookupMSISDN(input string) (*dto.NumberLookupResponse, *errs.AppError){
	
	ci,cc, ccLength, err := s.repo.LookupCountryCode(input)
	if err != nil {
		return nil, err
	}
	significantNumber := fmt.Sprint(input[ccLength:])
	
	mno, carrierLength, err := s.repo.LookupMobileOperator(ci, significantNumber)
	if err != nil{
		return nil, err
	}

	subscriberNumber := fmt.Sprint(significantNumber[carrierLength:])

	var response = dto.NumberLookupResponse{
		MNO: mno,
		SN: subscriberNumber,
		CI: ci,
		CC: cc,
	}
	return &response, nil
}
