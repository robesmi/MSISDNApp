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
	LookupMSISDN(int) (*dto.NumberLookupResponse, *errs.AppError)
}

// LookupMSISDN takes a full MSISDN as a string and returns
// a response containing the MNO, country code, subscriber number
// and the country identifier in ISO 3166-1-alpha-2 format
// or an error otherwise
func (s DefaultMSISDNService) LookupMSISDN(input int) (*dto.NumberLookupResponse, *errs.AppError){
	
	//	We need to check the country code first due to different countries having
	//	subscriber numbers of different lengths

	panic("panic!!")
}