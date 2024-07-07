package redis

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	*redis.Client
}

type RedisSettings struct {
	DataSource string `env:"REDIS_URL"`
}

func NewRedis(settings RedisSettings) (*Redis, error) {
	opt, err := redis.ParseURL(settings.DataSource)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	return &Redis{client}, nil
}
