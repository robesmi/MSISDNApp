package model

type User struct {
	UUID string			`db:"id"`
	Username string		`db:"username"`
	Password string		`db:"password"`
	Role string			`db:"role"`
	RefreshToken string	`db:"refresh_token"`
}