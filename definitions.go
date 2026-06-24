package main

const createTableQueries = `
CREATE TABLE IF NOT EXISTS activity (
    entity_number      VARCHAR(50),
    activity_group     VARCHAR(100),
    nace_version       VARCHAR(20),
    nace_code          VARCHAR(20),
    classification     VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS address (
    entity_number      VARCHAR(50),
    type_of_address    VARCHAR(50),
    country_nl         VARCHAR(100),
    country_fr         VARCHAR(100),
    zipcode            VARCHAR(20),
    municipality_nl    VARCHAR(150),
    municipality_fr    VARCHAR(150),
    street_nl          VARCHAR(255),
    street_fr          VARCHAR(255),
    house_number       VARCHAR(50),
    box                VARCHAR(50),
    extra_address_info VARCHAR(255),
    date_striking_off  DATE
);

CREATE TABLE IF NOT EXISTS branch (
    id               VARCHAR(50),
    start_date       DATE,
    enterprise_number VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS code (
    category     VARCHAR(100),
    code         VARCHAR(100),
    language     VARCHAR(10),
    description  VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS contact (
    entity_number  VARCHAR(50),
    entity_contact VARCHAR(100),
    contact_type   VARCHAR(100),
    value          VARCHAR(500)
);

CREATE TABLE IF NOT EXISTS denomination (
    entity_number      VARCHAR(50),
    language          VARCHAR(10),
    type_of_denomination VARCHAR(100),
    denomination      VARCHAR(320) NULL
);

CREATE TABLE IF NOT EXISTS enterprise (
    enterprise_number   VARCHAR(50),
    status              VARCHAR(100),
    juridical_situation VARCHAR(100),
    type_of_enterprise  VARCHAR(100),
    juridical_form      VARCHAR(100),
    juridical_form_cac  VARCHAR(100),
    start_date          DATE
);

CREATE TABLE IF NOT EXISTS establishment (
    establishment_number VARCHAR(50),
    start_date           DATE,
    enterprise_number    VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS meta (
    variable VARCHAR(255),
    value    VARCHAR(100)
);
`

type TableConfig struct {
	Name    string
	CSVFile string
}

func getImportOrder() []string {
	return []string{
		"meta",
		"code",
		"enterprise",
		"establishment",
		"address",
		"activity",
		"contact",
		"denomination",
	}
}

var tables = []TableConfig{
	{
		Name:    "meta",
		CSVFile: "meta.csv",
	},
	{
		Name:    "code",
		CSVFile: "code.csv",
	},
	{
		Name:    "enterprise",
		CSVFile: "enterprise.csv",
	},
	{
		Name:    "establishment",
		CSVFile: "establishment.csv",
	},
	{
		Name:    "address",
		CSVFile: "address.csv",
	},
	{
		Name:    "activity",
		CSVFile: "activity.csv",
	},
	{
		Name:    "contact",
		CSVFile: "contact.csv",
	},
	{
		Name:    "denomination",
		CSVFile: "denomination.csv",
	},
}
