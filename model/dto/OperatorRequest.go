package dto

type OperatorRequest struct {
	CountryIdentifier	string	`form:"countryidentifier"`
	PrefixFormat		string	`form:"prefixformat"`
	MNO					string	`form:"mno"`
	PrefixLength		string	`form:"prefixlength"`
}