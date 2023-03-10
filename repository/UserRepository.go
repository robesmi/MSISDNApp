package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/robesmi/MSISDNApp/model"
	"github.com/robesmi/MSISDNApp/model/errs"
)

type UserRepositoryDb struct {
	client *sqlx.DB
}

func NewAuthRepository(client *sqlx.DB) UserRepositoryDb {
	return UserRepositoryDb{client}
}

type UserRepository interface {
	GetUserByUsername(string) (*model.User, error)
	GetUserById(string) (*model.User, error)
	// RegisterNativeUser takes a UUID, username, password, role, JWT Access and Refresh tokens and saves the user
	// in the db, returning an error if unsuccessful
	RegisterNativeUser(string, string, string, string, string) error
	// RegisterNativeUser takes a username, role, JWT Access and Refresh tokens and saves the user
	// in the db, returning an error if unsuccessful
	RegisterImportedUser(string, string, string, string) error
	UpdateRefreshToken(string, string) error
}

func (db UserRepositoryDb) GetUserByUsername(username string) (*model.User, error){
	var user model.User
	sqlFind := "SELECT id, username, password, role, refresh_token FROM users WHERE username = ?"
	err := db.client.Get(&user, sqlFind, username)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, errs.NewUserNotFoundError()
		}else{
			return nil, errs.NewUnexpectedError(err.Error())
		}
	}
	return &user, nil

}

func (db UserRepositoryDb) GetUserById(id string) (*model.User, error){
	var user model.User
	sqlFind := "SELECT id, username, password, role, refresh_token FROM users WHERE id = ?"
	err := db.client.Get(&user, sqlFind, id)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, errs.NewUserNotFoundError()
		}else{
			return nil, errs.NewUnexpectedError(err.Error())
		}
	}
	return &user, nil

}

func (db UserRepositoryDb) RegisterNativeUser(uuid string, username string, password string, role string, refresh_token string) (error){
	
	sqlNewUser := "INSERT INTO users VALUES (?,?,?,?,?)"
	_, execError := db.client.Exec(sqlNewUser, uuid, username, password, role, refresh_token)
	if execError != nil{
		return errs.NewUnexpectedError(execError.Error())
	}

	return nil
}


func (db UserRepositoryDb) RegisterImportedUser(uuid string, username string, role string, refresh_token string)  error{

	sqlNewUser := "INSERT INTO users VALUES (?,?,?,?)"
	_, execError := db.client.Exec(sqlNewUser, uuid, username, role, refresh_token)
	if execError != nil{
		return errs.NewUnexpectedError(execError.Error())
	}
	return nil
}


func (db UserRepositoryDb) UpdateRefreshToken(uuid string, refreshToken string) error{
	sqlRefresh := "UPDATE users SET refresh_token = ? WHERE id = ?"
	_, refreshErr := db.client.Exec(sqlRefresh, refreshToken, uuid)
	if refreshErr != nil{
		return errs.NewUnexpectedError(refreshErr.Error())
	}
	return nil
}
