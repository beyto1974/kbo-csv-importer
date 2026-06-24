package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

func parseDate(s string) (bun.NullTime, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "NULL" {
		return bun.NullTime{}, nil
	}
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02 15:04:05",
		"02/01/2006",
		"02-01-2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return bun.NullTime{Time: t}, nil
		}
	}
	return bun.NullTime{}, fmt.Errorf("invalid date: %q", s)
}

func strOrNil(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" || s == "NULL" {
		return nil
	}
	v := strings.Trim(s, `"`)
	return &v
}

func readHeaderMap(r *csv.Reader) (map[string]int, error) {
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	m := make(map[string]int, len(header))
	for i, h := range header {
		m[strings.TrimSpace(h)] = i
	}
	return m, nil
}

func field(rec []string, idx int) string {
	if idx < 0 || idx >= len(rec) {
		return ""
	}
	return rec[idx]
}

func insertBatch[T any](ctx context.Context, db bun.IDB, batch []T) error {
	if len(batch) == 0 {
		return nil
	}
	_, err := db.NewInsert().Model(&batch).Exec(ctx)
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

	headerMap, err := readHeaderMap(r)
	if err != nil {
		return fmt.Errorf("failed to read header from %s: %w", csvPath, err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	batchSize := config.BatchSize
	if batchSize <= 0 {
		batchSize = 500
	}

	inserted := 0
	rowNum := 1

	totalRecords, err := countLines(csvPath)

	//Decrement header
	totalRecords--

	lastShown := -1

	if err != nil {
		return fmt.Errorf("failed to count lines: %w", err)
	}

	if config.Verbose {
		fmt.Printf("  Progress: %s\n", progressBar(0, totalRecords, 20))
	}

	switch table.Name {
	case "activity":
		batch := make([]Activity, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}
			if len(rec) < 5 {
				return fmt.Errorf("row %d: expected 5 fields, got %d", rowNum, len(rec))
			}

			row := Activity{
				EntityNumber:   field(rec, headerMap["entity_number"]),
				ActivityGroup:  field(rec, headerMap["activity_group"]),
				NaceVersion:    field(rec, headerMap["nace_version"]),
				NaceCode:       field(rec, headerMap["nace_code"]),
				Classification: field(rec, headerMap["classification"]),
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert activity failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert activity failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "address":
		batch := make([]Address, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}

			dt, err := parseDate(field(rec, headerMap["DateStrikingOff"]))
			if err != nil {
				return fmt.Errorf("row %d: %w", rowNum, err)
			}

			row := Address{
				EntityNumber:     field(rec, headerMap["EntityNumber"]),
				TypeOfAddress:    field(rec, headerMap["TypeOfAddress"]),
				CountryNl:        field(rec, headerMap["CountryNL"]),
				CountryFr:        field(rec, headerMap["CountryFR"]),
				Zipcode:          field(rec, headerMap["Zipcode"]),
				MunicipalityNl:   field(rec, headerMap["MunicipalityNL"]),
				MunicipalityFr:   field(rec, headerMap["MunicipalityFR"]),
				StreetNl:         field(rec, headerMap["StreetNL"]),
				StreetFr:         field(rec, headerMap["StreetFR"]),
				HouseNumber:      field(rec, headerMap["HouseNumber"]),
				Box:              field(rec, headerMap["Box"]),
				ExtraAddressInfo: field(rec, headerMap["ExtraAddressInfo"]),
				DateStrikingOff:  dt,
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert address failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert address failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "branch":
		batch := make([]Branch, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}

			dt, err := parseDate(field(rec, headerMap["StartDate"]))
			if err != nil {
				return fmt.Errorf("row %d: %w", rowNum, err)
			}

			row := Branch{
				ID:               field(rec, headerMap["ID"]),
				StartDate:        dt,
				EnterpriseNumber: field(rec, headerMap["EnterpriseNumber"]),
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert branch failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert branch failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "code":
		batch := make([]Code, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}

			row := Code{
				Category:    field(rec, headerMap["Category"]),
				Code:        field(rec, headerMap["Code"]),
				Language:    field(rec, headerMap["Language"]),
				Description: field(rec, headerMap["Description"]),
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert code failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert code failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "contact":
		batch := make([]Contact, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}

			row := Contact{
				EntityNumber:  field(rec, headerMap["EntityNumber"]),
				EntityContact: field(rec, headerMap["EntityContact"]),
				ContactType:   field(rec, headerMap["ContactType"]),
				Value:         field(rec, headerMap["Value"]),
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert contact failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert contact failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "denomination":
		batch := make([]Denomination, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}

			row := Denomination{
				EntityNumber:       field(rec, headerMap["EentityNumber"]),
				Language:           field(rec, headerMap["Language"]),
				TypeOfDenomination: field(rec, headerMap["TypeOfDenomination"]),
				Denomination:       strOrNil(field(rec, headerMap["Denomination"])),
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert denomination failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert denomination failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "enterprise":
		batch := make([]Enterprise, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			dt, err := parseDate(field(rec, headerMap["StartDate"]))
			if err != nil {
				return fmt.Errorf("row %d: %w", rowNum, err)
			}

			row := Enterprise{
				EnterpriseNumber:   field(rec, headerMap["EnterpriseNumber"]),
				Status:             field(rec, headerMap["Status"]),
				JuridicalSituation: field(rec, headerMap["JuridicalSituation"]),
				TypeOfEnterprise:   field(rec, headerMap["TypeOfEnterprise"]),
				JuridicalForm:      field(rec, headerMap["JuridicalForm"]),
				JuridicalFormCac:   field(rec, headerMap["HuridicalFormCAC"]),
				StartDate:          dt,
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert enterprise failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert enterprise failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "establishment":
		batch := make([]Establishment, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}

			dt, err := parseDate(field(rec, headerMap["StartDate"]))
			if err != nil {
				return fmt.Errorf("row %d: %w", rowNum, err)
			}

			row := Establishment{
				EstablishmentNumber: field(rec, headerMap["EstablishmentNumber"]),
				StartDate:           dt,
				EnterpriseNumber:    field(rec, headerMap["EnterpriseNumber"]),
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert establishment failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert establishment failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	case "meta":
		batch := make([]Meta, 0, batchSize)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read csv row %d: %w", rowNum+1, err)
			}
			rowNum++

			if config.TestRun && inserted >= 100 {
				break
			}

			row := Meta{
				Variable: field(rec, headerMap["Variable"]),
				Value:    field(rec, headerMap["Value"]),
			}
			batch = append(batch, row)
			if len(batch) >= batchSize {
				if err := insertBatch(ctx, tx, batch); err != nil {
					return fmt.Errorf("batch insert meta failed: %w", err)
				}
				inserted += len(batch)
				batch = batch[:0]

				printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, false)
			}
		}
		if err := insertBatch(ctx, tx, batch); err != nil {
			return fmt.Errorf("final batch insert meta failed: %w", err)
		}
		inserted += len(batch)

		printProgressBar(config.Verbose, rowNum, totalRecords, lastShown, true)

	default:
		return fmt.Errorf("unknown table %q", table.Name)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if config.Verbose {
		fmt.Printf("Inserted %d rows into %s\n", inserted, table.Name)
	}
	return nil
}
