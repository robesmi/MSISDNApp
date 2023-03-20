package dto

type CountryRequest struct {
	CountryNumberFormat	string	`form:"countryformat"`
	CountryCode			string	`form:"countrycode"`
	CountryIdentifier	string	`form:"countryidentifier"`
	CountryCodeLength	string	`form:"countrycodelength"`
}