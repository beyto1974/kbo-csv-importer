package explorer

import "beyto1974/kbo-csv-importer/private/db"

type EntityBundle struct {
	Enterprise     *db.Enterprise
	Establishments []db.Establishment
	Activities     []db.Activity
	Addresses      []db.Address
	Contacts       []db.Contact
	Denominations  []db.Denomination
	Branches       []db.Branch
	Codes          []db.Code
	Meta           []db.Meta
}
