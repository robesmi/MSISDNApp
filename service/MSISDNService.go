package service

import "github.com/robesmi/MSISDNApp/repository"
import "github.com/robesmi/MSISDNApp/model/dto"
import "github.com/robesmi/MSISDNApp/model/errs"

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
	
	//	We need to check the country code first due to different countries having
	//	MSISDN numbers of different lengths. Country code is 1-3 digits, NA being an exception

	// If it starts with 1, it is a NA number and we grab 3 more digits to see exact country in region

	// Otherwise, query with 2 digits and 3 digits and either receive a country identifier or error
	// If 2 digit query yielded result, continue with that result, otherwise if 3 digits query yielded result proceed with that result
	// If error, no countries match and return error


	panic("panic!!")
}