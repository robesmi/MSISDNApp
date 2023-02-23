package dto

type NumberLookupResponse struct{
	MNO string	`json:"MNO identifier" db:"mno"`
	CC string	`json:"Country Code" db:"country_code"`
	SN string	`json:"Subscriber Number""`
	CI string	`json:"Country Identifier" db:"country_identifier"`
}

func (r NumberLookupResponse) Compare(a NumberLookupResponse) bool {
	return r.MNO == a.MNO && r.CC == a.CC && r.SN == a.SN && r.CI == a.CI
}