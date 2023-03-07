package dto

type LoginResponse struct {
	Role string
	AccessToken string
	RefreshToken string
}