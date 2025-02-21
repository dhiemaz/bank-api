package migration

import (
	"database/sql"
	"github.com/dhiemaz/bank-api/config"
	db "github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
)

func RunMigration() {
	config := config.GetConfig()
	conn := db.InitDatabase(config)
	migrateDB(conn, "")
}

func migrateDB(conn *sql.DB, migrationURL string) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}

	if migrationURL == "" {
		migrationURL = "file://infrastructure/db/migration"
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
