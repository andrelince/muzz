package pg

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

type PostgresSettings struct {
	MaxOpenConns int    `env:"POSTGRES_MAX_OPEN"`
	MaxIdleConns int    `env:"POSTGRES_MAX_IDLE"`
	DataSource   string `env:"POSTGRES_URL"`
}

func NewPostgres(config PostgresSettings) (*Postgres, error) {
	db, err := sqlx.Open("pgx", config.DataSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)

	return &Postgres{db}, nil
}

func (p Postgres) Close() error {
	return p.db.Close()
}

func (p Postgres) Raw() *sql.DB {
	return p.db.DB
}
