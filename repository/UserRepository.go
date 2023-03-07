package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/model"
	"github.com/robesmi/MSISDNApp/model/errs"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryDb struct {
	client *sqlx.DB
}

func NewAuthRepository(client *sqlx.DB) UserRepositoryDb {
	return UserRepositoryDb{client}
}

type UserRepository interface {
	GetUserByEmail(string)
	// RegisterNativeUser inserts a new local user with username, role and refresh token
	RegisterNativeUser(string, string) (*model.Login, *errs.AppError)
	// GetNativeUser will fetch a local user with username, role and refresh token
	GetNativeUser(string,string) (*model.Login, *errs.AppError)
	RegisterImportedUser(string) (*model.Login, *errs.AppError)
	GetImportedUser(string) (*model.Login, *errs.AppError)
}

func (db UserRepositoryDb) GetUserByEmail(username string) (*model.User, *errs.AppError){
	var user model.User
	sqlFind := "SELECT id, username, password, refresh_token FROM users WHERE username = ?"
	err := db.client.Get(&user, sqlFind, username)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, errs.UserNotFound()
		}else{
			return nil, errs.UnexpectedError(fmt.Sprintf("Error in function GetUserByEmail: %s",err.Error()))
		}
	}
	return &user, nil

}

func (db UserRepositoryDb) RegisterNativeUser(username string, password string) (*model.Login, *errs.AppError){

	panic("Panic!!!")
}

func (db UserRepositoryDb) GetNativeUser(username string, password string) (*model.Login, *errs.AppError){
	var login model.Login
	sqlFind := "SELECT role, username, access_token, refresh_token FROM users WHERE username = ? AND password = ?"
	passwordCheck, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil{
		log.Println(err)
	}
	queryErr := db.client.Get(&login, sqlFind, username, string(passwordCheck))
	if queryErr != nil {
		if queryErr == sql.ErrNoRows{
			return nil, errs.InvalidInputError("Invalid username or password")
		}else {
			return nil, errs.UnexpectedError("Unexpected database error")
		}
	}
	return &login, nil
}

func (db UserRepositoryDb) RegisterImportedUser(username string, password string) (*model.Login, *errs.AppError){

	panic("panic!!!")
}

func (db UserRepositoryDb) GetImportedUser(username string){
	panic("panic!!")
}