# KBO CSV Importer

Imports Belgian KBO CSV dumps into SQLite, PostgreSQL, and MySQL.

## Usage

```bash
./kbo-importer --driver=<sqlite|postgres|mysql> --dsn=<dsn> [--batchSize=<batchSize?>] [--overwrite] [--verbose] <csv-dumps-path>
```

- `--driver`: Database driver to use (`sqlite`, `postgres`, or `mysql`)
- `--dsn`: Database connection string
- `--batchSize`: Defaults to 500
- `--overwrite`: Delete existing data before import (truncate)
- `--verbose`: Show progress bar
- `<sqlite-db-path>`: Kept for compatibility if your command still accepts it
- `<csv-dumps-path>`: Directory with CSV files

Example usage:

```bash
kbo-csv-importer -driver="sqlite" --dsn="file:mydb.db?cache=shared&mode=rwc" --verbose --overwrite storage/import
kbo-csv-importer -driver="postgres" --dsn="postgres://kbo:kbopassword@localhost:5432/kbo" --verbose --overwrite storage/import
```

### Binaries

Download prebuilt binaries from the [releases page](https://github.com/beyto1974/kbo-csv-importer/releases).

## Required CSV Files

```
meta.csv, code.csv, enterprise.csv, establishment.csv,
address.csv, activity.csv, contact.csv, denomination.csv
```

## Features

- Foreign key constraints enabled + validated
- Fails on any error (shows row data on FK violation)
- Date parsing: `YYYY-MM-DD` or `DD-MM-YYYY` → `YYYY-MM-DD`
- Primary keys on all tables
- NULL handling for empty strings

## Build

```bash
go build -o kbo-importer main.go
```