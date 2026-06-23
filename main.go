package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := parseConfig()

	if config.Verbose {
		fmt.Println("KBO CSV Importer")

		fmt.Printf("\nConfiguration:\n")
		fmt.Printf("  Data Dir:   %s\n", config.DataDir)
		fmt.Printf("  Overwrite:  %v\n", config.Overwrite)
		fmt.Printf("  Verbose:    %v\n\n", config.Verbose)
	}

	db, err := connectDB(config)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	ctx := context.Background()

	if config.Overwrite {
		clearTables(ctx, db, config.Driver, tables)
	}

	if _, err := db.ExecContext(context.Background(), createTableQueries); err != nil {
		panic(err)
	}

	// Import tables
	tableMap := map[string]TableConfig{}
	for _, t := range tables {
		tableMap[t.Name] = t
	}

	order := getImportOrder()
	totalTables, successfulTables := 0, 0

	for _, name := range order {
		t, ok := tableMap[name]
		if !ok {
			continue
		}
		totalTables++
		fmt.Printf("Importing table: %s (%s)\n", t.Name, t.CSVFile)
		if err := importTable(db, t, config); err != nil {
			fmt.Printf("  ERROR: %v\n", err)
			fmt.Println("\nImport failed on error. Exiting.")
			os.Exit(1)
		}
		successfulTables++
		fmt.Printf("  SUCCESS\n")
	}

	fmt.Printf("\nSummary: %d/%d tables imported successfully\n", successfulTables, totalTables)
}
