package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uptrace/bun"

	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/schema"
)

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

type Config struct {
	DataDir   string
	Driver    string
	DSN       string
	Overwrite bool
	Verbose   bool
	TestRun   bool
}

type TableConfig struct {
	Name        string
	CSVFile     string
	Columns     []string
	DateColumns []string
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
		Columns: []string{"Variable", "Value"},
	},
	{
		Name:    "code",
		CSVFile: "code.csv",
		Columns: []string{"Category", "Code", "Language", "Description"},
	},
	{
		Name:        "enterprise",
		CSVFile:     "enterprise.csv",
		Columns:     []string{"EnterpriseNumber", "Status", "JuridicalSituation", "TypeOfEnterprise", "JuridicalForm", "JuridicalFormCac", "StartDate"},
		DateColumns: []string{"StartDate"},
	},
	{
		Name:        "establishment",
		CSVFile:     "establishment.csv",
		Columns:     []string{"EstablishmentNumber", "StartDate", "EnterpriseNumber"},
		DateColumns: []string{"StartDate"},
	},
	{
		Name:        "address",
		CSVFile:     "address.csv",
		Columns:     []string{"EntityNumber", "TypeOfAddress", "CountryNl", "CountryFr", "Zipcode", "MunicipalityNl", "MunicipalityFr", "StreetNl", "StreetFr", "HouseNumber", "Box", "ExtraAddressInfo", "DateStrikingOff"},
		DateColumns: []string{"DateStrikingOff"},
	},
	{
		Name:    "activity",
		CSVFile: "activity.csv",
		Columns: []string{"EntityNumber", "ActivityGroup", "NaceVersion", "NaceCode", "Classification"},
	},
	{
		Name:    "contact",
		CSVFile: "contact.csv",
		Columns: []string{"EntityNumber", "EntityContact", "ContactType", "Value"},
	},
	{
		Name:    "denomination",
		CSVFile: "denomination.csv",
		Columns: []string{"EntityNumber", "Language", "TypeOfDenomination", "Denomination"},
	},
}

func connectDB(config Config) (*bun.DB, error) {
	switch config.Driver {
	case "postgres", "postgresql", "pg":
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.DSN)))
		return openBunDB(sqldb, pgdialect.New(), "postgres")

	case "mysql", "mariadb":
		sqldb, err := sql.Open("mysql", config.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to open mysql database: %w", err)
		}
		return openBunDB(sqldb, mysqldialect.New(), "mysql")

	case "sqlite", "sqlite3":
		sqldb, err := sql.Open(sqliteshim.ShimName, config.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite database: %w", err)
		}
		return openBunDB(sqldb, sqlitedialect.New(), "sqlite")

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
}

func openBunDB(sqldb *sql.DB, dialect schema.Dialect, name string) (*bun.DB, error) {
	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping %s database: %w", name, err)
	}
	return bun.NewDB(sqldb, dialect), nil
}

func parseDate(dateStr string) (*string, error) {
	if dateStr == "" || dateStr == "NULL" {
		return nil, nil
	}

	if len(dateStr) == 10 && strings.Contains(dateStr, "-") {
		if t, err := time.Parse("02-01-2006", dateStr); err == nil {
			s := t.Format("2006-01-02")
			return &s, nil
		}
	}

	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		s := t.Format("2006-01-02")
		return &s, nil
	}

	s := dateStr
	return &s, nil
}

func prepareInsertQuery(table TableConfig) string {
	placeholders := strings.TrimRight(strings.Repeat("?,", len(table.Columns)), ",")
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.Name, strings.Join(table.Columns, ", "), placeholders)
}

func printRowData(columns []string, rec []string) {
	fmt.Println("  Row data:")
	for i, col := range columns {
		val := rec[i]
		if val == "" {
			val = "<NULL>"
		}
		fmt.Printf("    %s: %s\n", col, val)
	}
}

func insertRow(ctx context.Context, tx bun.Tx, table TableConfig, row []interface{}) error {
	values := make(map[string]interface{}, len(table.Columns))
	for i, col := range table.Columns {
		// lower: postgres requirement
		values[ToSnakeCase(col)] = row[i]
	}

	_, err := tx.NewInsert().
		Model(&values).
		ModelTableExpr(string(bun.Ident(table.Name))).
		Exec(ctx)

	return err
}

func importTable(db *bun.DB, table TableConfig, config Config) error {
	csvPath := filepath.Join(config.DataDir, table.CSVFile)
	f, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file %s: %w", csvPath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file %s: %w", csvPath, err)
	}
	if len(records) == 0 {
		return fmt.Errorf("no records found in %s", table.CSVFile)
	}
	// Ignore first row
	if len(records) > 1 {
		records = records[1:]
	}
	if len(records) == 0 {
		return fmt.Errorf("no data records after skipping header in %s", table.CSVFile)
	}

	dateIdx := map[int]bool{}
	for _, dc := range table.DateColumns {
		for i, c := range table.Columns {
			if c == dc {
				dateIdx[i] = true
			}
		}
	}

	if config.Verbose {
		fmt.Printf("  Progress: %s\n", progressBar(0, len(records), 20))
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	ctx := context.Background()

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement for %s: %w; %s", table.Name, err, prepareInsertQuery(table))
	}

	inserted := 0
	lastShown := -1

	for idx, rec := range records {
		// Test
		if config.TestRun && idx > 100 {
			break
		}
		if len(rec) != len(table.Columns) {
			tx.Rollback()
			fmt.Printf("\n  Field count mismatch at row %d:\n", idx+1)
			printRowData(table.Columns, rec)
			return fmt.Errorf("row %d: field count mismatch (%d != %d) in %s", idx+1, len(rec), len(table.Columns), table.CSVFile)
		}

		vals := make([]interface{}, len(rec))
		for i, v := range rec {
			if dateIdx[i] {
				if p, err := parseDate(v); err == nil {
					vals[i] = p
				} else {
					vals[i] = v
				}
			} else if v == "" || v == "NULL" {
				vals[i] = nil
			} else {
				vals[i] = strings.Trim(v, "\"")
			}
		}

		if err := insertRow(ctx, tx, table, vals); err != nil {
			panic(err)
			// }

			// if _, err := stmt.Exec(vals...); err != nil {
			tx.Rollback()
			fmt.Printf("\n  Insert failed at row %d in %s:\n", idx+1, table.CSVFile)
			fmt.Printf("    Error: %v\n", err)
			printRowData(table.Columns, rec)
			return fmt.Errorf("row %d insert failed in %s: %w", idx+1, table.CSVFile, err)
		}
		inserted++

		if config.Verbose {
			p := ((idx + 1) * 100) / len(records)
			if p != lastShown && (p == 100 || p%2 == 0) {
				fmt.Printf("\r  Progress: %s", progressBar(idx+1, len(records), 20))
				lastShown = p
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if config.Verbose {
		fmt.Println()
	}
	fmt.Printf("  Inserted: %d\n", inserted)
	return nil
}

func clearTables(ctx context.Context, db *bun.DB, dialect string, tables []TableConfig) error {
	for _, t := range tables {
		var query string

		switch dialect {
		case "sqlite":
			query = fmt.Sprintf("DELETE FROM %s", t.Name)
		case "mysql":
			query = fmt.Sprintf("TRUNCATE TABLE %s", t.Name)
		case "postgres":
			query = fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", t.Name)
		default:
			return fmt.Errorf("unsupported dialect: %s", dialect)
		}

		if _, err := db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("clearing table %s: %w", t.Name, err)
		}
	}
	return nil
}
