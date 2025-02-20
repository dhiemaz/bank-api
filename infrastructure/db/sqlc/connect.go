package db

import (
	"database/sql"
	"github.com/dhiemaz/bank-api/config"
	"log"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func InitDatabase(config *config.Config) *sql.DB {
	conn, err := sql.Open(config.Database.Driver, config.Database.URL)
	if err != nil {
		log.Fatalf("cannot open connection to db, err: %s", err)
	}
	
	return conn
}

func migrateDB(conn *sql.DB, migrationURL string) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}

	if migrationURL == "" {
		migrationURL = "file://db/migration"
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		"postgres", driver)

	if err != nil {
		return err
	}

	m.Up()

	if err != nil {
		return err
	}

	return nil
}
