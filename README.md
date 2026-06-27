# KBO CSV Importer & Explorer

Tools for working with Belgian KBO (Kruispuntbank van Ondernemingen) open data CSV dumps.

- **Importer** — bulk-imports KBO CSV files into SQLite, PostgreSQL, or MySQL
- **Explorer** — looks up an enterprise by number and outputs a structured, LLM-ready description with all related data

## Importer

### Usage

```bash
./kbo-csv-importer --driver=<sqlite|postgres|mysql> --dsn=<dsn> [--batchSize=<batchSize?>] [--overwrite] [--verbose] <csv-dumps-path>
```

- `--driver`: Database driver to use (`sqlite`, `postgres`, or `mysql`)
- `--dsn`: Database connection string
- `--batchSize`: Defaults to 500
- `--overwrite`: Delete existing data before import (truncate)
- `--verbose`: Show progress bar
- `<csv-dumps-path>`: Directory with CSV files

On a decent computer/server, this operation takes less than 3 minutes.

### Examples

```bash
kbo-csv-importer -driver="sqlite" --dsn="file:mydb.db?cache=shared&mode=rwc" --verbose --overwrite storage/import
kbo-csv-importer -driver="postgres" --dsn="postgres://kbo:kbopassword@localhost:5432/kbo" --verbose --overwrite storage/import
```

## Explorer

### Usage

```bash
./kbo-explorer --driver=<sqlite|postgres|mysql> --dsn=<dsn> --enterpriseNumber=<number> [--language=<FR|NL>]
```

- `--driver`: Database driver (`sqlite`, `postgres`, or `mysql`)
- `--dsn`: Database connection string
- `--enterpriseNumber`: The enterprise or establishment number to look up (e.g. `0200.065.765`)
- `--language`: Language for code descriptions (`FR` or `NL`, defaults to `FR`)

### Example

```bash
kbo-explorer --driver="sqlite" --dsn="file:mydb.db?cache=shared&mode=rwc" --enterpriseNumber="0200.065.765" --language=FR
```

### Example Output

```
# Enterprise

- **Enterprise Number**: 0200.065.765
- **Status**: Actif (AC)
- **Juridical Situation**: Situation normale (000)
- **Type of Enterprise**: Personne morale (2)
- **Juridical Form**: Société anonyme (014)
- **Juridical Form (CAC)**: Grande entreprise (2)
- **Start Date**: 1994-06-30

## Denominations / Names

- [FR] PROXIMUS (type: Dénomination sociale, entity: 0200.065.765)
- [FR] PROXIMUS (type: Dénomination commerciale, entity: 0200.065.765)
- [NL] PROXIMUS (type: Maatschappelijke benaming, entity: 0200.065.765)

## Establishments

1. Establishment 2.003.015.769 (start: 2000-01-01, enterprise: 0200.065.765)
2. Establishment 2.236.225.479 (start: 2014-03-05, enterprise: 0200.065.765)

## Addresses

### Address 1 (entity: 0200.065.765, type: Adresse du siège social)
  Boulevard du Roi Albert II 27, 1030 Bruxelles, Belgique

### Address 2 (entity: 2.003.015.769, type: Adresse de l'unité d'établissement)
  Boulevard du Roi Albert II 27, 1030 Bruxelles, Belgique

## Contacts

- Site internet: https://www.proximus.be (contact-of: ENT, entity: 0200.065.765)
- Adresse e-mail: info@proximus.com (contact-of: ENT, entity: 0200.065.765)
- Numéro de téléphone: +32 2 202 41 11 (contact-of: ENT, entity: 0200.065.765)

## Activities / NACE Codes

- [0200.065.765] Group: Activité TVA | NACE 61100: Télécommunications filaires | Version: Nace 2008
- [2.003.015.769] Group: Activité TVA | NACE 61200: Télécommunications sans fil | Version: Nace 2008
```

## Explorer HTTP server

This is a simple REST server, usefull for integrations.

Start the server:

```bash
kbo-explorer-http --driver="sqlite" --dsn="file:mydb.db?cache=shared&mode=rwc" --enterpriseNumber="0200.065.765" --language=FR --port=8080
```

Endpoint:

```bash
GET /entity/{entityNumber}
```

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

### Binaries

Download prebuilt binaries from the [releases page](https://github.com/beyto1974/kbo-csv-importer/releases).

## Build

```bash
go build -o kbo-csv-importer main.go
go build -o kbo-explorer ./private/explorer/
```

## Development

```bash
go run ./private/importer/ <options...>
go run ./private/explorer/ <options...>
```
