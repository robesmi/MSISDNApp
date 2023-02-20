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
	LookupCountryCode(int) (string, *errs.AppError)
	LookupMobileOperator(string, int) (string, *errs.AppError)
}

// LookupCountryCode takes a int country code and returns the respective country identifier it belongs to
func (repo MSISDNRepositoryDb) LookupCountryCode(cc int) (string, *errs.AppError){	
	panic("panic!!")
}

// LookupMobileOperator takes a country identifier and a NCD/NPA code and returns an MNO
func (repo MSISDNRepositoryDb) LookupMobileOperator(cc string, prefix int) (string, *errs.AppError){
	panic("panic!!")
}
