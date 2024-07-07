package config

import (
	"github.com/muzz/api/pkg/logger"
	"github.com/muzz/api/pkg/pg"
	"github.com/muzz/api/pkg/redis"
)

type Config struct {
	Port             string `env:"SRV_PORT" envDefault:"3000"`
	MigrationPath    string `env:"MIGRATION_PATH"`
	PostgresSettings pg.PostgresSettings
	RedisSettings    redis.RedisSettings
	LoggerSettings   logger.Settings
}

func NewPostgresSettings(config Config) pg.PostgresSettings {
	return config.PostgresSettings
}

func NewRedisSettings(config Config) redis.RedisSettings {
	return config.RedisSettings
}

func NewLoggerSettings(config Config) logger.Settings {
	return config.LoggerSettings
}

func NewMigrationSettings(config Config) pg.PgMigrationSettings {
	return pg.PgMigrationSettings{
		MigrationPath: config.MigrationPath,
		DataSource:    config.PostgresSettings.DataSource,
	}
}
