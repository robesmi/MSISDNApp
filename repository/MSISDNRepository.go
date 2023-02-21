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
	// LookupCountryCode takes a string full number and returns the respective country identifier it belongs to,
	// the country code and the country's prefix length, or an error
	LookupCountryCode(string) (string, int,  int, *errs.AppError)
	// LookupMobileOperator takes a country identifier and a significant number and returns an MNO, length of carrier prefix, or an error
	LookupMobileOperator(string, string) (string, int, *errs.AppError)
}


func (repo MSISDNRepositoryDb) LookupCountryCode(fullnumber string) (string, int, int,  *errs.AppError){	
	panic("panic!!")
}


func (repo MSISDNRepositoryDb) LookupMobileOperator(significantNumber string, prefix string) (string, int,  *errs.AppError){
	panic("panic!!")
}
