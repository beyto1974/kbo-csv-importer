package db

const CreateTableQueries = `
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

-- activity
CREATE INDEX idx_activity_entity_number
ON activity (entity_number);

CREATE INDEX idx_activity_nace_code
ON activity (nace_code);

CREATE INDEX idx_activity_group_code
ON activity (activity_group, nace_version, nace_code);

-- address
CREATE INDEX idx_address_entity_number
ON address (entity_number);

CREATE INDEX idx_address_zipcode
ON address (zipcode);

CREATE INDEX idx_address_municipality_nl
ON address (municipality_nl);

CREATE INDEX idx_address_municipality_fr
ON address (municipality_fr);

CREATE INDEX idx_address_type_entity
ON address (type_of_address, entity_number);

-- branch
CREATE INDEX idx_branch_id
ON branch (id);

CREATE INDEX idx_branch_enterprise_number
ON branch (enterprise_number);

CREATE INDEX idx_branch_start_date
ON branch (start_date);

-- code
CREATE INDEX idx_code_category_code_lang
ON code (category, code, language);

-- contact
CREATE INDEX idx_contact_entity_number
ON contact (entity_number);

CREATE INDEX idx_contact_type_value
ON contact (contact_type, value);

-- denomination
CREATE INDEX idx_denomination_entity_lang
ON denomination (entity_number, language);

CREATE INDEX idx_denomination_type
ON denomination (type_of_denomination);

-- enterprise
CREATE INDEX idx_enterprise_number
ON enterprise (enterprise_number);

CREATE INDEX idx_enterprise_status
ON enterprise (status);

CREATE INDEX idx_enterprise_juridical_form
ON enterprise (juridical_form);

CREATE INDEX idx_enterprise_start_date
ON enterprise (start_date);

-- establishment
CREATE INDEX idx_establishment_number
ON establishment (establishment_number);

CREATE INDEX idx_establishment_enterprise_number
ON establishment (enterprise_number);

CREATE INDEX idx_establishment_start_date
ON establishment (start_date);

-- meta
CREATE UNIQUE INDEX idx_meta_variable
ON meta (variable);
`

type TableConfig struct {
	Name    string
	CSVFile string
}

func GetTableMap() map[string]TableConfig {
	tableMap := map[string]TableConfig{}
	for _, t := range tables {
		tableMap[t.Name] = t
	}
	return tableMap
}

func GetImportOrder() []string {
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
