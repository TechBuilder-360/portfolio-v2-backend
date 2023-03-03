package model

type Certificate struct {
	Base

	UserId     string
	Name       string
	Issuer     string
	IssuerDate string
}
