package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/model/errs"
)

type MSISDNRepositoryDb struct {
	db *sqlx.DB
}

func NewMSISDNRepository(dbClient *sqlx.DB) MSISDNRepositoryDb{
	return MSISDNRepositoryDb{dbClient}
}

type MSISDNRepository interface{
	LookupCountryCode(string) (string, int, *errs.AppError)
	LookupMobileOperator(string, string) (string, int, *errs.AppError)
}

// LookupCountryCode takes a string country code and returns the respective country identifier it belongs to, the country's prefix length, or an error
func (repo MSISDNRepositoryDb) LookupCountryCode(cc string) (string, int, *errs.AppError){	
	panic("panic!!")
}

// LookupMobileOperator takes a country identifier and a NCD/NPA code and returns an MNO
func (repo MSISDNRepositoryDb) LookupMobileOperator(cc string, prefix string) (string, int,  *errs.AppError){
	panic("panic!!")
}
