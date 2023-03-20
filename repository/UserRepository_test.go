package repository

import (
	_ "database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/model"
)

var sqlxDb *sqlx.DB
var userRepo UserRepository


func setup(t *testing.T) (sqlmock.Sqlmock){
	db, mock, err := sqlmock.New()
	if err != nil{
		t.Errorf("an error %s was not expected when opening stub db connection", err.Error())
	}
	sqlxDb = sqlx.NewDb(db,"sqlmock")
	userRepo = NewAuthRepository(sqlxDb)
	lookupRepo = NewMSISDNRepository(sqlxDb)


	return mock
}

func TestGetAllUsers(t *testing.T) {

	//Arrange
	mock := setup(t)
	expUser := model.User{
		UUID: "13",
		Username: "u1",
		Password: "p1",
		Role: "user",
		RefreshToken: "rt",
	}
	rows := mock.NewRows([]string{"id","username","password","role","refresh_token"}).
	AddRow(expUser.UUID, expUser.Username,expUser.Password,expUser.Role, expUser.RefreshToken)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	//Act
	_, getErr := userRepo.GetAllUsers()

	//Assert

	if getErr != nil{
		t.Errorf("Error in TestGetAllUsers:\n expected %s\n got %s", "nil", getErr)
	}
}

func TestGetuserByUsernameValid(t *testing.T) {

	//Arrange
	mock := setup(t)
	expUser := model.User{
		UUID: "13",
		Username: "u1",
		Password: "p1",
		Role: "user",
		RefreshToken: "rt",
	}
	rows := mock.NewRows([]string{"id","username","password","role","refresh_token"}).
	AddRow(expUser.UUID, expUser.Username,expUser.Password,expUser.Role, expUser.RefreshToken)
	mock.ExpectQuery("SELECT").WithArgs(expUser.Username).WillReturnRows(rows)
	
	//Act
	resp, getErr := userRepo.GetUserByUsername(expUser.Username)

	//Assert
	if getErr != nil{
		t.Errorf("Error in TestGetuserByUsernameValid:\n expected %s\n got %s", "nil", getErr)
	}
	if resp.Username != expUser.Username{
		t.Errorf("Error in TestGetuserByUsernameValid:\n expected %s\n got %s", expUser.Username, resp.Username)
	}
}

func TestGetUserByIdValid(t *testing.T) {

	// Arrange
	mock := setup(t)

	expUser := model.User{
		UUID: "13",
		Username: "u1",
		Password: "p1",
		Role: "user",
		RefreshToken: "rt",
	}
	rows := mock.NewRows([]string{"id","username","password","role","refresh_token"}).
	AddRow(expUser.UUID, expUser.Username,expUser.Password,expUser.Role, expUser.RefreshToken)
	mock.ExpectQuery("SELECT").WithArgs(expUser.UUID).WillReturnRows(rows)
	
	//Act
	resp , getErr := userRepo.GetUserById(expUser.UUID)


	//Assert
	if getErr != nil{
		t.Errorf("Error in TestGetUserByIdValid:\n expected %s\n got %s", "nil", getErr)
	}
	if resp.UUID != expUser.UUID{
		t.Errorf("Error in TestGetUserByIdValid:\n expected %s\n got %s", expUser.UUID, resp.UUID)
	}
}


func TestRegisterNativeUser(t *testing.T) {

	// Arrange
	mock := setup(t)

	expUser := model.User{
		UUID: "13",
		Username: "u1",
		Password: "p1",
		Role: "user",
		RefreshToken: "rt",
	}
	mock.ExpectExec("INSERT INTO users").WithArgs(expUser.UUID,expUser.Username, expUser.Password, expUser.Role, expUser.RefreshToken).
	WillReturnResult(sqlmock.NewResult(1,1))
	//Act
	insertErr := userRepo.RegisterNativeUser(expUser.UUID, expUser.Username, expUser.Password,expUser.Role, expUser.RefreshToken)



	//Assert
	if insertErr != nil{
		t.Errorf("Error in TestRegisterNativeUser:\n expected %s\n got %s", "nil", insertErr)
	}
}

func TestRegisterImportedUser(t *testing.T) {

	// Arrange
	mock := setup(t)

	expUser := model.User{
		UUID: "13",
		Username: "u1",
		Password: "p1",
		Role: "user",
		RefreshToken: "rt",
	}
	mock.ExpectExec("INSERT INTO users").WithArgs(expUser.UUID,expUser.Username, "", expUser.Role, expUser.RefreshToken).
	WillReturnResult(sqlmock.NewResult(1,1))
	//Act
	insertErr := userRepo.RegisterImportedUser(expUser.UUID, expUser.Username, expUser.Role, expUser.RefreshToken)


	//Assert
	if insertErr != nil{
		t.Errorf("Error in TestRegisterImportedUser:\n expected %s\n got %s", "nil", insertErr)
	}
}

func TestUpdateRefreshToken(t *testing.T) {

	// Arrange
	mock := setup(t)

	
	mock.ExpectExec("UPDATE users").WithArgs("rt","id").
	WillReturnResult(sqlmock.NewResult(1,1))
	//Act
	insertErr := userRepo.UpdateRefreshToken("id","rt")


	//Assert
	if insertErr != nil{
		t.Errorf("Error in TestUpdateRefreshToken:\n expected %s\n got %s", "nil", insertErr)
	}
}

func TestEditUserByIdWithPassword(t *testing.T) {

	// Arrange
	mock := setup(t)
	expUser := model.User{
		UUID: "13",
		Username: "u1",
		Password: "p1",
		Role: "user",
		RefreshToken: "rt",
	}
	
	mock.ExpectExec("UPDATE users").WithArgs(expUser.Username, expUser.Password, expUser.Role, expUser.UUID).
	WillReturnResult(sqlmock.NewResult(1,1))
	//Act
	insertErr := userRepo.EditUserById(expUser.UUID, expUser.Username, expUser.Password, expUser.Role)


	//Assert
	if insertErr != nil{
		t.Errorf("Error in TestEditUserByIdWithPassword:\n expected %s\n got %s", "nil", insertErr)
	}
}

func TestEditUserByIdNoPassword(t *testing.T) {

	// Arrange
	mock := setup(t)
	expUser := model.User{
		UUID: "13",
		Username: "u1",
		Password: "p1",
		Role: "user",
		RefreshToken: "rt",
	}
	
	mock.ExpectExec("UPDATE users").WithArgs(expUser.Username, expUser.Role, expUser.UUID).
	WillReturnResult(sqlmock.NewResult(1,1))
	//Act
	insertErr := userRepo.EditUserById(expUser.UUID, expUser.Username, "", expUser.Role)


	//Assert
	if insertErr != nil{
		t.Errorf("Error in TestEditUserByIdNoPassword:\n expected %s\n got %s", "nil", insertErr)
	}
}

func TestRemoveUserById(t *testing.T) {

	// Arrange
	mock := setup(t)
	
	mock.ExpectExec("DELETE FROM users").WithArgs("id").
	WillReturnResult(sqlmock.NewResult(1,1))
	//Act
	insertErr := userRepo.RemoveUserById("id")


	//Assert
	if insertErr != nil{
		t.Errorf("Error in TestRemoveUserById:\n expected %s\n got %s", "nil", insertErr)
	}
}