package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Driver   string
	DSN      string
	Language string
	Port     int
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
		fmt.Printf("Usage: %s --driver=<sqlite|postgres|mysql> --dsn=<dns> --language=<language>\n", os.Args[0])
		flag.PrintDefaults()
	}

	dsn := flag.String("dsn", "", "DSN")
	driver := flag.String("driver", "", "Database driver")
	language := flag.String("language", "", "FR, NL or DE")
	port := flag.Int("port", 8080, "HTTP port")

	flag.Parse()

	cfg := Config{
		DSN:      *dsn,
		Driver:   *driver,
		Language: *language,
		Port:     *port,
	}

	if err := validateConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}

	return cfg
}
