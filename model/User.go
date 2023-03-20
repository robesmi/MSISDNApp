package model

type User struct {
	UUID string			`db:"id" form:"id"`
	Username string		`db:"username" form:"username"`
	Password string		`db:"password" form:"password"`
	Role string			`db:"role" form:"role"`
	RefreshToken string	`db:"refresh_token" form:"ref_token"`
}