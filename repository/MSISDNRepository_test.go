package repository

import (
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

var db *sqlx.DB
var repo MSISDNRepository

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "latest", []string{"MYSQL_ROOT_PASSWORD=testpassword"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sqlx.Open("mysql", fmt.Sprintf("root:testpassword@(localhost:%s)/mysql?multiStatements=true", resource.GetPort("3306/tcp")))
		if err != nil {
			log.Println("Container not ready, waiting...")
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	} 
	dbSetup()
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func dbSetup() {
	a := "CREATE DATABASE test"
	db.Exec(a)
	firstSetup := "CREATE TABLE countries ("+
		"`country_number_format` varchar(20) NOT NULL,"+
		"`country_code` varchar(6) NOT NULL,"+
		"`country_identifier` varchar(3) NOT NULL,"+
		"`country_code_length` int NOT NULL,"+
		"PRIMARY KEY (`country_number_format`))"
	db.Exec(firstSetup)
	secondSetup := "INSERT INTO countries VALUES (\"^389[0-9]{8}$\",389,\"mk\",3),"+
	"\n(\"^350[0-9]{5}$\",350,\"gi\",3),"+
    "\n(\"^242[0-9]{9}$\",242,\"cg\",3),"+
    "\n(\"^423[0-9]{8}$\",423,\"li\",3),"+
    "\n(\"^48[0-9]{9}$\",48,\"pl\",2);"
	db.Exec(secondSetup)

	thirdSetup := "CREATE TABLE mobile_operators("+
		"`country_identifier` varchar(3) NOT NULL,"+
		"`prefix_format` varchar(60) NOT NULL,"+
		"`mno`           varchar(100) NOT NULL,"+
		"`prefix_length` int NOT NULL,"+
		"PRIMARY KEY (`country_identifier`, `prefix_format`))"
	db.Exec(thirdSetup)
	
	fourthSetup := "INSERT INTO `mobile_operators` VALUES"+
    "\n(\"mk\",\"^77[0-9]{6}$\",\"A1\",2),"+
    "\n(\"mk\",\"^71[0-9]{6}$\", \"Telekom\",2),"+
    "\n(\"li\",\"^6[0-9]{7}$\", \"Lietuvos\",1),"+
    "\n(\"pl\",\"^510[0-9]{6}$\", \"Mobile telephoOrange\",2),"+
    "\n(\"pl\",\"(?!^5329[0-9]{5}$)(?!^5366[0-9]{5}$)^53[0-7][0-9]{6}$\", \"Orange Polska S.A\",2),"+
    "\n(\"pl\",\"^53(2|8|9)[0-9]{6}$\", \"T-MOBILE POLSKA S.A.\",2),"+
    "\n(\"pl\",\"^5366[0-9]{6}$\", \"Polskie Sieci Cyfrowe Sp. z o.o.\",2);"
	db.Exec(fourthSetup)

}

func TestInvalidCountry(t *testing.T) {
	inputNumber := "11177554333"
	repo = NewMSISDNRepository(db)

	_,err := repo.LookupCountryCode(inputNumber)

	if err == nil{
		t.Error("Failed while testing invalid country")
	}
}

func TestValidCountry(t *testing.T) {
	inputNumber := "38977123456"
	repo = NewMSISDNRepository(db)

	_,err := repo.LookupCountryCode(inputNumber)

	if err != nil{
		t.Error("Failed while testing valid country")
	}
}

func TestInvalidOperator(t *testing.T){
	inputNumber, inputCi := "2312351", "mk"
	repo = NewMSISDNRepository(db)

	_,err := repo.LookupMobileOperator(inputNumber,inputCi)

	if err == nil{
		t.Error("Failed while testing invalid operator")
	}
}
func TestValidOperator(t *testing.T){
	inputNumber, inputCi := "77123456", "mk"
	repo = NewMSISDNRepository(db)

	_,err := repo.LookupMobileOperator(inputCi,inputNumber)

	if err != nil{
		t.Error("Failed while testing valid operator")
	}
}


