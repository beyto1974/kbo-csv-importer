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

func ParseConfig() Config {
	flag.Usage = func() {
		fmt.Printf("Usage: %s --driver=<sqlite|postgres|mysql> --dsn=<dns> --enterpriseNumber=<enterpriseNumber> --language=<language>\n", os.Args[0])
		flag.PrintDefaults()
	}

	dsn := flag.String("dsn", "", "DSN")
	driver := flag.String("driver", "", "Database driver")
	enterpriseNumber := flag.String("enterpriseNumber", "", "")
	language := flag.String("language", "", "")

	flag.Parse()

	if *driver == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: -driver")
		flag.Usage()
		os.Exit(2)
	}

	return Config{
		DSN:              *dsn,
		Driver:           *driver,
		EnterpriseNumber: *enterpriseNumber,
		Language:         *language,
	}
}
