package pg

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

type PgMigrationSettings struct {
	MigrationPath string
	DataSource    string
}

type PgMigration struct {
	db       *sql.DB
	settings PgMigrationSettings
	logger   *logrus.Logger
}

func (m PgMigration) Up() error {
	defer func() {
		if err := m.Close(); err != nil {
			m.logger.Fatal(context.Background(), "goose: failed to close DB: %v\n", nil)
		}
	}()

	return goose.Up(m.db, m.settings.MigrationPath)
}

func (m PgMigration) Close() error {
	return m.db.Close()
}

func NewPgMigration(logger *logrus.Logger, settings PgMigrationSettings) (PgMigration, error) {
	db, err := goose.OpenDBWithDriver("postgres", settings.DataSource)
	if err != nil {
		return PgMigration{}, err
	}

	return PgMigration{db: db, settings: settings, logger: logger}, nil
}
