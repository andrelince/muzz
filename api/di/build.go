package di

import (
	"net/http"

	"github.com/muzz/api/config"
	"github.com/muzz/api/pkg/env"
	"github.com/muzz/api/pkg/logger"
	"github.com/muzz/api/pkg/pg"
	"github.com/muzz/api/pkg/redis"
	"github.com/muzz/api/repository"
	"github.com/muzz/api/rest"
	"github.com/muzz/api/rest/middleware"
	"github.com/muzz/api/service"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func buildConfig(c *dig.Container) error {
	if err := c.Provide(func() (config.Config, error) {
		return env.New(config.Config{})
	}); err != nil {
		return err
	}

	if err := c.Provide(config.NewPostgresSettings); err != nil {
		return err
	}

	if err := c.Provide(config.NewMigrationSettings); err != nil {
		return err
	}

	if err := c.Provide(config.NewRedisSettings); err != nil {
		return err
	}

	if err := c.Provide(config.NewLoggerSettings); err != nil {
		return err
	}

	if err := c.Provide(logger.New); err != nil {
		return err
	}

	if err := c.Provide(func(l *logrus.Logger, c pg.PostgresSettings) (*pg.Postgres, error) {
		p, err := pg.NewPostgres(c)
		if err != nil {
			l.Error("failed to open database connection")
		}
		return p, err
	}); err != nil {
		return err
	}

	if err := c.Provide(pg.NewPgMigration); err != nil {
		return err
	}

	if err := c.Provide(func(l *logrus.Logger, c redis.RedisSettings) (*redis.Redis, error) {
		r, err := redis.NewRedis(c)
		if err != nil {
			l.Error("failed to open cache connection")
		}
		return r, err
	}); err != nil {
		return err
	}

	if err := c.Provide(http.NewServeMux); err != nil {
		return err
	}

	if err := c.Provide(func(l *logrus.Logger, r *redis.Redis, config config.Config) repository.AuthConnector {
		return repository.NewAuthRepo(l, r, config.SecretKey)
	}); err != nil {
		return err
	}

	if err := c.Provide(func(l *logrus.Logger, p *pg.Postgres) repository.UserConnector {
		return repository.NewUserRepo(l, p)
	}); err != nil {
		return err
	}

	if err := c.Provide(func(r repository.UserConnector, a repository.AuthConnector) service.UserConnector {
		return service.NewUserService(r, a)
	}); err != nil {
		return err
	}

	if err := c.Provide(func(r repository.AuthConnector) service.AuthConnector {
		return service.NewAuthService(r)
	}); err != nil {
		return err
	}

	if err := c.Provide(func(s service.AuthConnector) middleware.AuthMiddleware {
		return middleware.NewAuthHandler(s)
	}); err != nil {
		return err
	}

	if err := c.Provide(rest.NewHandler); err != nil {
		return err
	}

	return nil
}
