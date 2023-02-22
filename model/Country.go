package model

import "github.com/robesmi/MSISDNApp/model/dto"

type Country struct {
	// CountryNumberFormat is a regex representing the country code and format
	// of each country, identified by the first 1-4 digits
	CountryNumberFormat string	`db:"country_number_format"`
	// CountryCode is an int with the country code needed to dial the country
	CountryCode string			`db:"country_code"`
	// CountryIdentifier is a ISO 3166-1-alpha-2 format of the country
	// the MSISDN belongs to
	CountryIdentifier string	`db:"country_identifier"`
	// CountryCodeLength is the length of the CountryCode, used to
	// trim away the unneeded country code in following queries
	CountryCodeLength int		`db:"country_code_length"`
}

func (c *Country) toDto() dto.CountryLookupResponse{ 
	return dto.CountryLookupResponse{
		CountryCode: c.CountryCode,
		CountryIdentifier: c.CountryIdentifier,
		CountryCodeLength: c.CountryCodeLength,
	}
}