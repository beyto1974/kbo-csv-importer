package explorer

import (
	"beyto1974/kbo-csv-importer/private/db"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func LoadEnterpriseBundle(ctx context.Context, bunDB *bun.DB, inputNumber string, language string) (*EntityBundle, error) {
	bundle := &EntityBundle{}

	var enterpriseNumber string

	// 1) Try enterprise first
	err := bunDB.NewSelect().
		Model((*db.Enterprise)(nil)).
		Column("enterprise_number").
		Where("enterprise_number = ?", inputNumber).
		Scan(ctx, &enterpriseNumber)

	if err != nil {
		// 2) Fallback: input might be an establishment_number
		err = bunDB.NewSelect().
			Model((*db.Establishment)(nil)).
			Column("enterprise_number").
			Where("establishment_number = ?", inputNumber).
			Scan(ctx, &enterpriseNumber)
		if err != nil {
			// Returns also error if not found
			if fmt.Sprint(err) == "sql: no rows in result set" {
				return nil, fmt.Errorf("Enterprise number of establishment number not found")
			}

			return nil, err
		}
	}

	// 3) Load enterprise
	bundle.Enterprise = new(db.Enterprise)

	err = bunDB.NewSelect().
		Model(bundle.Enterprise).
		ColumnExpr("e.*").
		ColumnExpr("s.code AS status_code__code").
		ColumnExpr("s.description AS status_code__description").
		ColumnExpr("js.code AS juridical_situation_code__code").
		ColumnExpr("js.description AS juridical_situation_code__description").
		ColumnExpr("toe.code AS type_of_enterprise_code__code").
		ColumnExpr("toe.description AS type_of_enterprise_code__description").
		ColumnExpr("jf.code AS juridical_form_code__code").
		ColumnExpr("jf.description AS juridical_form_code__description").
		ColumnExpr("jfc.code AS juridical_form_cac_code__code").
		ColumnExpr("jfc.description AS juridical_form_cac_code__description").
		Join(`LEFT JOIN code AS s
		ON s.category = ? AND s.code = e.status AND s.language = ?`, "Status", language).
		Join(`LEFT JOIN code AS js
		ON js.category = ? AND js.code = e.juridical_situation AND js.language = ?`, "JuridicalSituation", language).
		Join(`LEFT JOIN code AS toe
		ON toe.category = ? AND toe.code = e.type_of_enterprise AND toe.language = ?`, "TypeOfEnterprise", language).
		Join(`LEFT JOIN code AS jf
		ON jf.category = ? AND jf.code = e.juridical_form AND jf.language = ?`, "JuridicalForm", language).
		Join(`LEFT JOIN code AS jfc
		ON jfc.category = ? AND jfc.code = e.juridical_form_cac AND jfc.language = ?`, "JuridicalFormCac", language).
		Where("e.enterprise_number = ?", enterpriseNumber).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	// 4) Load all establishments for enterprise
	err = bunDB.NewSelect().
		Model(&bundle.Establishments).
		Where("enterprise_number = ?", enterpriseNumber).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 5) Build targets = enterprise + all establishments
	targets := make([]string, 0, len(bundle.Establishments)+1)
	targets = append(targets, enterpriseNumber)
	for _, es := range bundle.Establishments {
		targets = append(targets, es.EstablishmentNumber)
	}

	// 6) Load all related tables using IN (...)
	err = bunDB.NewSelect().
		Model(&bundle.Activities).
		ColumnExpr("a.*").
		ColumnExpr("ag.code AS activity_group_code__code").
		ColumnExpr("ag.description AS activity_group_code__description").
		ColumnExpr("n.description AS nace_code_code__description").
		ColumnExpr("cl.code AS classification_code__code").
		ColumnExpr("cl.description AS classification_code__description").
		Join(`LEFT JOIN code AS ag
		ON ag.category = ? AND ag.code = a.activity_group AND ag.language = ?`, "ActivityGroup", language).
		Join(`
			LEFT JOIN code AS n
			ON n.category = CONCAT('Nace', a.nace_version)
			AND n.code = a.nace_code
			AND n.language = ?
		`, language).
		Join(`LEFT JOIN code AS cl
		ON cl.category = ? AND cl.code = a.classification AND cl.language = ?`, "Classification", language).
		Where("a.entity_number IN (?)", bun.List(targets)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	err = bunDB.NewSelect().
		Model(&bundle.Addresses).
		ColumnExpr("ad.*").
		ColumnExpr("toa.code AS type_of_address_code__code").
		ColumnExpr("toa.description AS type_of_address_code__description").
		Join(`LEFT JOIN code AS toa
		ON toa.category = ? AND toa.code = ad.type_of_address AND toa.language = ?`, "TypeOfAddress", language).
		Where("ad.entity_number IN (?)", bun.List(targets)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	err = bunDB.NewSelect().
		Model(&bundle.Contacts).
		ColumnExpr("co.*").
		ColumnExpr("ct.code AS contact_type_code__code").
		ColumnExpr("ct.description AS contact_type_code__description").
		Join(`LEFT JOIN code AS ct
		ON ct.category = ? AND ct.code = co.contact_type AND ct.language = ?`, "ContactType", language).
		Where("co.entity_number IN (?)", bun.List(targets)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	err = bunDB.NewSelect().
		Model(&bundle.Denominations).
		ColumnExpr("d.*").
		ColumnExpr("tod.code AS type_of_denomination_code__code").
		ColumnExpr("tod.description AS type_of_denomination_code__description").
		Join(`LEFT JOIN code AS tod
		ON tod.category = ? AND tod.code = d.type_of_denomination AND tod.language = ?`, "TypeOfDenomination", language).
		Where("d.entity_number IN (?)", bun.List(targets)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 7) Branches use enterprise_number
	err = bunDB.NewSelect().
		Model(&bundle.Branches).
		Where("enterprise_number = ?", enterpriseNumber).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	// 8) Codes + meta if needed
	err = bunDB.NewSelect().
		Model(&bundle.Codes).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	err = bunDB.NewSelect().
		Model(&bundle.Meta).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return bundle, nil
}
