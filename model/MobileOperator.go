package model

type MobileOperator struct {
	// CountryIdentifier is a ISO 3166-1-alpha-2 format of the country
	// the MSISDN belongs to 
	CountryIdentifier string
	// Prefix is a regex to distinguish the NCD and format of numbers
	// for each MNO
	Prefix string
	// MNO is the name of the MNO the regex applies to
	MNO string
	// PrefixLength is the length of the MNO's carrier code, used to
	// trim away the unneeded carrier code to isolate the subscriber number
	PrefixLength int

}