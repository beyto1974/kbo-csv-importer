package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/schema"
)

func ConnectDB(driver string, dsn string) (*bun.DB, error) {
	switch driver {
	case "postgres", "postgresql", "pg":
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		return openBunDB(sqldb, pgdialect.New(), "postgres")

	case "mysql", "mariadb":
		sqldb, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open mysql database: %w", err)
		}
		return openBunDB(sqldb, mysqldialect.New(), "mysql")

	case "sqlite", "sqlite3":
		sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite database: %w", err)
		}
		return openBunDB(sqldb, sqlitedialect.New(), "sqlite")

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func openBunDB(sqldb *sql.DB, dialect schema.Dialect, name string) (*bun.DB, error) {
	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping %s database: %w", name, err)
	}
	return bun.NewDB(sqldb, dialect), nil
}

func ClearTables(ctx context.Context, db *bun.DB, dialect string) error {
	for _, t := range tables {
		var query string

		switch dialect {
		case "sqlite":
			query = fmt.Sprintf("DELETE FROM %s", t.Name)
		case "mysql":
			query = fmt.Sprintf("TRUNCATE TABLE %s", t.Name)
		case "postgres":
			query = fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY", t.Name)
		default:
			return fmt.Errorf("unsupported dialect: %s", dialect)
		}

		if _, err := db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("clearing table %s: %w", t.Name, err)
		}
	}
	return nil
}
