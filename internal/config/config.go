package config

import "bestprice/bestprice-api/internal/helper"

type Config struct {
	DbPath           string
	RedisPath        string
	RedisDbs         map[string]int
	DbMaxConnections int
	ApiAddress       string
	JwtKey           []byte
}

func NewConfig() *Config {
	return &Config{
		DbPath:           helper.GetEnv("BESTPRICE_MYSQL_PATH", "user:password@tcp(localhost:33066)/bestprice?parseTime=true"),
		RedisPath:        helper.GetEnv("BESTPRICE_REDIS_PATH", "localhost:63799"),
		RedisDbs:         map[string]int{"0": 0},
		DbMaxConnections: 50,
		ApiAddress:       helper.GetEnv("ADDRESS", ":8080"),
		JwtKey:           []byte("bestprice"),
	}
}
