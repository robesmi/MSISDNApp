package model

import "github.com/robesmi/MSISDNApp/model/dto"

type MobileOperator struct {
	// CountryIdentifier is a ISO 3166-1-alpha-2 format of the country
	// the MSISDN belongs to
	CountryIdentifier string	`db:"country_identifier"`
	// Prefix is a regex to distinguish the NCD and format of numbers
	// for each MNO
	PrefixFormat string			`db:"prefix_format"`
	// MNO is the name of the MNO the regex applies to
	MNO string					`db:"mno"`
	// PrefixLength is the length of the MNO's carrier code, used to
	// trim away the unneeded carrier code to isolate the subscriber number
	PrefixLength int			`db:"prefix_length"`
}

func (m *MobileOperator) toDto() dto.MobileOperatorLookupResponse{
	return dto.MobileOperatorLookupResponse{
		MNO: m.MNO,
		PrefixLength: m.PrefixLength,
	}
}