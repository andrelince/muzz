package config

import (
	"github.com/muzz/api/pkg/logger"
	"github.com/muzz/api/pkg/pg"
	"github.com/muzz/api/pkg/redis"
)

type Config struct {
	Port             string `env:"SRV_PORT" envDefault:"3000"`
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
