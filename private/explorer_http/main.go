package main

import (
	"beyto1974/kbo-csv-importer/private/db"
	"beyto1974/kbo-csv-importer/private/explorer"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	config := ParseConfig()
	ctx := context.Background()

	bunDB, err := db.ConnectDB(config.Driver, config.DSN)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	defer bunDB.Close()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	http.HandleFunc("/entity/{entityNumber}", func(w http.ResponseWriter, r *http.Request) {
		entityNumber := string(r.PathValue("entityNumber"))

		logger = log.New(os.Stdout, "http: ", log.LstdFlags)
		logger.Printf("Received request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		response, err := explorer.LoadEnterpriseBundle(ctx, bunDB, entityNumber, config.Language)
		if err != nil {
			fmt.Printf("Error retrieving bundle: %v\n", err)
			w.WriteHeader(404)
			fmt.Fprint(w, "custom 404")
			return
		}

		fmt.Fprint(w, explorer.FormatEnterpriseLLM(response))
	})

	logger.Printf("Server is starting on port %d...", config.Port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}
