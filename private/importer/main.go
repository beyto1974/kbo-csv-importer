package main

import (
	"beyto1974/kbo-csv-importer/private/db"
	"beyto1974/kbo-csv-importer/private/importer_config"
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := importer_config.ParseConfig()

	if config.Verbose {
		fmt.Println("KBO CSV Importer")

		fmt.Printf("\nConfiguration:\n")
		fmt.Printf("  Data Dir:   %s\n", config.DataDir)
		fmt.Printf("  Overwrite:  %v\n", config.Overwrite)
		fmt.Printf("  Verbose:    %v\n\n", config.Verbose)
	}

	bunDB, err := db.ConnectDB(config.Driver, config.DSN)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	defer bunDB.Close()

	ctx := context.Background()

	if config.Overwrite {
		db.ClearTables(ctx, bunDB, config.Driver)
	}

	if _, err := bunDB.ExecContext(context.Background(), db.CreateTableQueries); err != nil {
		panic(err)
	}

	// Import tables
	tableMap := db.GetTableMap()

	order := db.GetImportOrder()
	totalTables, successfulTables := 0, 0

	for _, name := range order {
		t, ok := tableMap[name]
		if !ok {
			continue
		}
		totalTables++
		fmt.Printf("Importing table: %s (%s)\n", t.Name, t.CSVFile)
		if err := importTable(bunDB, t, config); err != nil {
			fmt.Printf("  ERROR: %v\n", err)
			fmt.Println("\nImport failed on error. Exiting.")
			os.Exit(1)
		}
		successfulTables++
		fmt.Printf("  SUCCESS\n")
	}

	fmt.Printf("\nSummary: %d/%d tables imported successfully\n", successfulTables, totalTables)
}
