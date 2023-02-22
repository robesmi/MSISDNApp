package dto

type NumberLookupResponse struct{
	MNO string	`json:"MNO identifier"`
	CC int		`json:"Country Code"`
	SN string	`json:"Subscriber Number"`
	CI string	`json:"Country Identifier"`
}

func (r *NumberLookupResponse) Compare(a NumberLookupResponse) bool {
	return r.MNO == a.MNO && r.CC == a.CC && r.SN == a.SN && r.CI == a.CI
}