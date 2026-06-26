package main

import (
	"beyto1974/kbo-csv-importer/private/db"
	"beyto1974/kbo-csv-importer/private/explorer_config"
	"context"
	"fmt"
	"os"
)

func main() {
	config := explorer_config.ParseConfig()

	bunDB, err := db.ConnectDB(config.Driver, config.DSN)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	defer bunDB.Close()

	ctx := context.Background()

	response, err := LoadEnterpriseBundle(ctx, bunDB, config.EnterpriseNumber, config.Language)
	if err != nil {
		fmt.Printf("Error retrieving bundle: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(FormatEnterpriseLLM(response))
}
