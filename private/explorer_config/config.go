package explorer_config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Driver           string
	DSN              string
	EnterpriseNumber string
	Language         string
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
		require("enterpriseNumber", cfg.EnterpriseNumber),
		require("language", cfg.Language),
	} {
		if check != nil {
			return check
		}
	}
	return nil
}

func ParseConfig() Config {
	flag.Usage = func() {
		fmt.Printf("Usage: %s --driver=<sqlite|postgres|mysql> --dsn=<dns> --enterpriseNumber=<enterpriseNumber> --language=<language>\n", os.Args[0])
		flag.PrintDefaults()
	}

	dsn := flag.String("dsn", "", "DSN")
	driver := flag.String("driver", "", "Database driver")
	enterpriseNumber := flag.String("enterpriseNumber", "", "")
	language := flag.String("language", "", "FR, NL or DE")

	flag.Parse()

	cfg := Config{
		DSN:              *dsn,
		Driver:           *driver,
		EnterpriseNumber: *enterpriseNumber,
		Language:         *language,
	}

	if err := validateConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}

	return cfg
}
