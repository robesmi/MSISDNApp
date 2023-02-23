package dto

type CountryLookupResponse struct {
	CountryCode string	`db:"country_code"`
	CountryIdentifier string	`db:"country_identifier"`
	CountryCodeLength int		`db:"country_code_length"`
}