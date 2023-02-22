package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/model/dto"
	"github.com/robesmi/MSISDNApp/model/errs"
)

type MSISDNRepositoryDb struct {
	db *sqlx.DB
}

func NewMSISDNRepository(dbClient *sqlx.DB) MSISDNRepositoryDb{
	return MSISDNRepositoryDb{dbClient}
}

type MSISDNRepository interface{
	// LookupCountryCode takes a string full number and returns the respective country identifier it belongs to,
	// the country code and the country's prefix length, or an error
	LookupCountryCode(string) (*dto.CountryLookupResponse, *errs.AppError)
	// LookupMobileOperator takes a country identifier and a significant number and returns an MNO, length of carrier prefix, or an error
	LookupMobileOperator(string, string) (*dto.MobileOperatorLookupResponse, *errs.AppError)
}

//go:generate mockgen -destination=../mocks/repository/mockMSISDNRepository.go -package=repository github.com/robesmi/MSISDNApp/repository MSISDNRepository
func (repo MSISDNRepositoryDb) LookupCountryCode(fullnumber string) (*dto.CountryLookupResponse,  *errs.AppError){	
	panic("panic!!")
}

func (repo MSISDNRepositoryDb) LookupMobileOperator(significantNumber string, prefix string) (*dto.MobileOperatorLookupResponse, *errs.AppError){
	panic("panic!!")
}