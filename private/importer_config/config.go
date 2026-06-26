package importer_config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	DataDir   string
	Driver    string
	DSN       string
	BatchSize int
	Overwrite bool
	Verbose   bool
	TestRun   bool
}

func require(name, value string) error {
	if value == "" {
		return fmt.Errorf("missing required flag: -%s", name)
	}
	return nil
}

func validateConfig(cfg Config) error {
	for _, check := range []error{
		require("dsn", cfg.DSN),
		require("driver", cfg.Driver),
	} {
		if check != nil {
			return check
		}
	}
	return nil
}

func ParseConfig() Config {
	flag.Usage = func() {
		fmt.Printf("Usage: %s --driver=<sqlite|postgres|mysql> --dsn=<dns> [--overwrite] [--verbose] <csv-dumps-path>\n", os.Args[0])
		flag.PrintDefaults()
	}

	dsn := flag.String("dsn", "", "DSN")
	driver := flag.String("driver", "", "Database driver")
	batchSize := flag.Int("batchSize", 500, "Insert batch size")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing database")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	testRun := flag.Bool("testRun", false, "Import only first 100 rows of each table")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	dataDir := args[0]

	cfg := Config{
		DataDir:   dataDir,
		DSN:       *dsn,
		BatchSize: *batchSize,
		Driver:    *driver,
		Overwrite: *overwrite,
		Verbose:   *verbose,
		TestRun:   *testRun,
	}

	if err := validateConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}

	return cfg
}
