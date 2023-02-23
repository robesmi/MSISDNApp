package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
	var response dto.CountryLookupResponse
	sqlQuery := "SELECT country_code, country_identifier, country_code_length FROM countries WHERE ? RLIKE country_number_format"
	err := repo.db.Get(&response, sqlQuery, fullnumber)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, errs.NumberNotFoundError("Country not found")
		}else{
			return nil, errs.UnexpectedError("Unexpected database error")
		}
	}
	return &response,nil
}

func (repo MSISDNRepositoryDb) LookupMobileOperator(ci string, significantNumber string) (*dto.MobileOperatorLookupResponse, *errs.AppError){
	var response dto.MobileOperatorLookupResponse
	sqlQuery := "SELECT mno, prefix_length FROM mobile_operators WHERE ? = country_identifier AND ? RLIKE prefix_format"
	err := repo.db.Get(&response, sqlQuery, ci, significantNumber)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, errs.NoCarriersFound("Carrier not found")
		}else{
			return nil, errs.UnexpectedError("Unexpected database error")
		}
	}
	return &response,nil
}