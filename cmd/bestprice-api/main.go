package main

import (
	"bestprice/bestprice-api/internal/api"
	"bestprice/bestprice-api/internal/config"
	"bestprice/bestprice-api/internal/data/cache"
	"bestprice/bestprice-api/internal/data/database"

	"log"
)

func main() {
	conf := config.NewConfig()

	db, err := database.NewDb(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Connection.Close()

	redisHandler, err := cache.NewRedisHandler(conf)
	if err != nil {
		log.Fatal(err)
	}
	redisHandler.Client.FlushDB()
	defer redisHandler.Client.Close()

	api := api.NewApi(db, redisHandler, conf)
	api.Run()
}
