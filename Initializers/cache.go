package Initializers

import "github.com/redis/go-redis/v9"

var RDb *redis.Client

func InitCache() {
	RDb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
