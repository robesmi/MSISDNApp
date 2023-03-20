package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/model"
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
	LookupCountryCode(string) (*dto.CountryLookupResponse, error)
	// LookupMobileOperator takes a country identifier and a significant number and returns an MNO, length of carrier prefix, or an error
	LookupMobileOperator(string, string) (*dto.MobileOperatorLookupResponse, error)
	AddNewCountry(string, string, string, int) (error)
	AddNewMobileOperator(string, string, string, int) (error)
	GetAllCountries() (*[]model.Country, error)
	GetAllMobileOperators() (*[]model.MobileOperator, error)
}


//go:generate mockgen -destination=../mocks/repository/mockMSISDNRepository.go -package=repository github.com/robesmi/MSISDNApp/repository MSISDNRepository
func (repo MSISDNRepositoryDb) LookupCountryCode(fullnumber string) (*dto.CountryLookupResponse,  error){
	var response dto.CountryLookupResponse
	sqlQuery := "SELECT country_code, country_identifier, country_code_length FROM countries WHERE ? RLIKE country_number_format"
	err := repo.db.Get(&response, sqlQuery, fullnumber)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, errs.NewNumberNotFoundError()
		}else{
			return nil, errs.NewUnexpectedError(err.Error())
		}
	}
	return &response,nil
}

func (repo MSISDNRepositoryDb) LookupMobileOperator(ci string, significantNumber string) (*dto.MobileOperatorLookupResponse, error){
	var response dto.MobileOperatorLookupResponse
	sqlQuery := "SELECT mno, prefix_length FROM mobile_operators WHERE ? = country_identifier AND ? RLIKE prefix_format"
	err := repo.db.Get(&response, sqlQuery, ci, significantNumber)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, errs.NewNoCarriersFoundError()
		}else{
			return nil, errs.NewUnexpectedError(err.Error())
		}
	}
	return &response,nil
}

func (repo MSISDNRepositoryDb) GetAllCountries() (*[]model.Country,error){

	var response []model.Country
	sqlQuery := "SELECT * FROM countries"
	err := repo.db.Select(&response, sqlQuery)
	if err != nil{
		return nil, err
	}
	return &response,nil
}

func (repo MSISDNRepositoryDb) GetAllMobileOperators() (*[]model.MobileOperator, error){

	var response []model.MobileOperator
	sqlQuery := "SELECT * FROM mobile_operators"
	err := repo.db.Select(&response,sqlQuery)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repo MSISDNRepositoryDb) AddNewCountry(numFormat string, cc string, ci string, cLen int) (error){

	sqlAdd := "INSERT INTO countries VALUES (?,?,?,?)"
	_, err := repo.db.Exec(sqlAdd, numFormat,cc,ci,cLen)
	if err != nil{
		return err
	}
	return nil
}

func (repo MSISDNRepositoryDb) AddNewMobileOperator(ci string, prefix string, mno string, prefixLen int) (error){

	sqlAdd := "INSERT INTO mobile_operators VALUES (?,?,?,?)"
	_, err := repo.db.Exec(sqlAdd, ci, prefix, mno, prefixLen)
	if err != nil{
		return err
	}
	return nil
}