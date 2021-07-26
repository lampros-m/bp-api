package cache

import (
	"bestprice/bestprice-api/internal/data/services"
)

func NewRedisHandlerMock() *RedisHanlder {
	redisHandler := RedisHanlder{
		RedisServices: RedisServices{
			Cache: services.NewCacheServicesMock(),
		},
	}

	return &redisHandler
}
