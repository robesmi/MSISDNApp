package dto

type MobileOperatorLookupResponse struct {
	MNO string	`db:"mno"`
	PrefixLength int	`db:"prefix_length"`
}