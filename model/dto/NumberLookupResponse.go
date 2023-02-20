package dto

type NumberLookupResponse struct{
	MNO string	`json:"MNO identifier"`
	CC int		`json:"Country Code"`
	SN string	`json:"Subscriber Number"`
	CI string	`json:"Country Identifier"`
}