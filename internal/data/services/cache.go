package services

import (
	"bestprice/bestprice-api/internal/helper"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type CacheServices struct {
	RedisClient *redis.Client
	Expiration  time.Duration
}

func NewCacheServices(redisClient *redis.Client) *CacheServices {
	return &CacheServices{
		RedisClient: redisClient,
		Expiration:  1 * time.Hour,
	}
}

func (o *CacheServices) Set(key string, value string) error {
	var err error

	if !helper.ValidRedisKeyValuePair(key, value) {
		return errors.New("Redis Key or Value are not valid")
	}

	err = o.RedisClient.Set(key, value, o.Expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (o *CacheServices) Get(key string) (string, error) {
	var output string
	mockedValue := "v"

	if !helper.ValidRedisKeyValuePair(key, mockedValue) {
		return output, errors.New("Redis Key or Value are not valid")
	}

	value, err := o.RedisClient.Get(key).Result()
	if err != nil {
		if err.Error() != "redis: nil" {
			log.Println("Error Redis:", err)
		}
		return output, err
	}
	output = value

	return output, nil
}

func (o *CacheServices) Delete(key string) error {
	mockedValue := "v"

	if !helper.ValidRedisKeyValuePair(key, mockedValue) {
		return errors.New("Redis Key or Value are not valid")
	}

	_, err := o.RedisClient.Unlink(key).Result()
	if err != nil {
		return err
	}

	return nil
}
