package repository


import "github.com/jmoiron/sqlx"

type MSISDNRepositoryDb struct {
	db *sqlx.DB
}

func NewMSISDNRepository(dbClient *sqlx.DB) MSISDNRepositoryDb{
	return MSISDNRepositoryDb{dbClient}
}

type MSISDNRepository interface{
	LookupCountryCode(string) string
	LookupMobileOperator(string, int) string
}

// LookupCountryCode takes a int country code and returns the respective country identifier it belongs to
func (repo MSISDNRepositoryDb) LookupCountryCode(cc string) (string){
	panic("panic!!")
}

// LookupMobileOperator takes a country identifier and a NCD/NPA code and returns an MNO
func (repo MSISDNRepositoryDb) LookupMobileOperator(cc string, prefix int) (string){
	panic("panic!!")
}
