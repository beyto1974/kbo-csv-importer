package main

import (
	"github.com/uptrace/bun"
)

type Activity struct {
	bun.BaseModel  `bun:"table:activity,alias:a"`
	EntityNumber   string `bun:"entity_number"`
	ActivityGroup  string `bun:"activity_group"`
	NaceVersion    string `bun:"nace_version"`
	NaceCode       string `bun:"nace_code"`
	Classification string `bun:"classification"`
}

type Address struct {
	bun.BaseModel    `bun:"table:address,alias:ad"`
	EntityNumber     string       `bun:"entity_number"`
	TypeOfAddress    string       `bun:"type_of_address"`
	CountryNl        string       `bun:"country_nl"`
	CountryFr        string       `bun:"country_fr"`
	Zipcode          string       `bun:"zipcode"`
	MunicipalityNl   string       `bun:"municipality_nl"`
	MunicipalityFr   string       `bun:"municipality_fr"`
	StreetNl         string       `bun:"street_nl"`
	StreetFr         string       `bun:"street_fr"`
	HouseNumber      string       `bun:"house_number"`
	Box              string       `bun:"box"`
	ExtraAddressInfo string       `bun:"extra_address_info"`
	DateStrikingOff  bun.NullTime `bun:"date_striking_off"`
}

type Branch struct {
	bun.BaseModel    `bun:"table:branch,alias:b"`
	ID               string       `bun:"id"`
	StartDate        bun.NullTime `bun:"start_date"`
	EnterpriseNumber string       `bun:"enterprise_number"`
}

type Code struct {
	bun.BaseModel `bun:"table:code,alias:c"`
	Category      string `bun:"category"`
	Code          string `bun:"code"`
	Language      string `bun:"language"`
	Description   string `bun:"description"`
}

type Contact struct {
	bun.BaseModel `bun:"table:contact,alias:co"`
	EntityNumber  string `bun:"entity_number"`
	EntityContact string `bun:"entity_contact"`
	ContactType   string `bun:"contact_type"`
	Value         string `bun:"value"`
}

type Denomination struct {
	bun.BaseModel      `bun:"table:denomination,alias:d"`
	EntityNumber       string  `bun:"entity_number"`
	Language           string  `bun:"language"`
	TypeOfDenomination string  `bun:"type_of_denomination"`
	Denomination       *string `bun:"denomination,nullzero"`
}

type Enterprise struct {
	bun.BaseModel      `bun:"table:enterprise,alias:e"`
	EnterpriseNumber   string       `bun:"enterprise_number"`
	Status             string       `bun:"status"`
	JuridicalSituation string       `bun:"juridical_situation"`
	TypeOfEnterprise   string       `bun:"type_of_enterprise"`
	JuridicalForm      string       `bun:"juridical_form"`
	JuridicalFormCac   string       `bun:"juridical_form_cac"`
	StartDate          bun.NullTime `bun:"start_date"`
}

type Establishment struct {
	bun.BaseModel       `bun:"table:establishment,alias:es"`
	EstablishmentNumber string       `bun:"establishment_number"`
	StartDate           bun.NullTime `bun:"start_date"`
	EnterpriseNumber    string       `bun:"enterprise_number"`
}

type Meta struct {
	bun.BaseModel `bun:"table:meta,alias:m"`
	Variable      string `bun:"variable"`
	Value         string `bun:"value"`
}
