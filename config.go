package main

import (
	"flag"
	"fmt"
	"os"
)

func parseConfig() Config {
	flag.Usage = func() {
		fmt.Printf("Usage: %s --driver=<sqlite|postgres|mysql> --dsn=<dns> [--overwrite] [--verbose] <csv-dumps-path>\n", os.Args[0])
		flag.PrintDefaults()
	}

	dsn := flag.String("dsn", "", "DSN")
	driver := flag.String("driver", "", "Database driver")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing database")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	testRun := flag.Bool("testRun", false, "Import only first 100 rows of each table")

	flag.Parse()

	if *dsn == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: -dsn")
		flag.Usage()
		os.Exit(2)
	}

	if *driver == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: -driver")
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	dataDir := args[0]

	return Config{
		DataDir:   dataDir,
		DSN:       *dsn,
		Driver:    *driver,
		Overwrite: *overwrite,
		Verbose:   *verbose,
		TestRun:   *testRun,
	}
}
