package repository

import (
	"database/sql"
	"log"

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
//go:generate mockgen -destination=../mocks/repository/mockUserRepository.go -package=repository github.com/robesmi/MSISDNApp/repository UserRepository

type UserRepository interface {
	// Can you figure out what this does?
	GetAllUsers() (*[]model.User, error)
	// GetUserByUsername takes a username and returns a full user if found, a UserNotFoundError if no
	// such user is found, or UnexpectedError otherwise
	GetUserByUsername(string) (*model.User, error)
	// GetUserByUsername takes a uuid and returns a full user if found, a UserNotFoundError if no
	// such user is found, or UnexpectedError otherwise
	GetUserById(string) (*model.User, error)
	// RegisterNativeUser takes a UUID, username, password, role, JWT Refresh token and saves the user
	// in the db, returning an error if unsuccessful
	RegisterNativeUser(string, string, string, string, string) error
	// RegisterImporteduser takes a username, role, JWT Refresh token and saves the user
	// in the db, returning an error if unsuccessful
	RegisterImportedUser(string, string, string, string) error
	// UpdateRefreshToken takes a uuid and a refresh token and updates the user's
	// refresh token, returning an error if unsuccessful
	UpdateRefreshToken(string, string) error
}

func (db UserRepositoryDb) GetAllUsers() (*[]model.User, error){

	var allUsers []model.User
	sqlGet := "SELECT * FROM users"
	err := db.client.Select(&allUsers, sqlGet)
	if err != nil{
		log.Println("Error in GetAllUsers: " + err.Error())
		return nil, err
	}
	return &allUsers, nil

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

	sqlNewUser := "INSERT INTO users VALUES (?,?, ?,?,?)"
	_, execError := db.client.Exec(sqlNewUser, uuid, username,"", role, refresh_token)
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
