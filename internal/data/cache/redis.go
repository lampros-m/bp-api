package cache

import (
	"bestprice/bestprice-api/internal/config"
	"bestprice/bestprice-api/internal/data/services"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type RedisHanlder struct {
	Client    *redis.Client
	Db        int
	Connected bool
	RedisServices
}

type RedisServices struct {
	Cache ICacheServices
}

func NewRedisHandler(conf *config.Config) (*RedisHanlder, error) {
	redisOptions := &redis.Options{
		Addr: conf.RedisPath,
		DB:   conf.RedisDbs["0"],
	}

	client := redis.NewClient(redisOptions)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "Cannot connect to redis db")
	}

	redisHandler := RedisHanlder{
		Client:    client,
		Db:        conf.RedisDbs["0"],
		Connected: true,
		RedisServices: RedisServices{
			Cache: services.NewCacheServices(client),
		},
	}

	return &redisHandler, nil
}

type ICacheServices interface {
	Set(string, string) error
	Get(string) (string, error)
	Delete(string) error
}
