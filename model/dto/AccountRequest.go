package dto

type AccountRequest struct {
	Username 	string	`form:"username"`
	Password 	string	`form:"password"`
	Role 		string	`form:"role"`
}