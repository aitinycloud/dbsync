package rediscli

import (
	"fmt"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func RedisInit(addrs string, password string) *redis.Client {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     addrs,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	pong, err := redisdb.Ping().Result()
	if err != nil || pong != "PONG" {
		panic(fmt.Sprintln("Redis ping,Is redis config error ? , error : ", err))
	}
	RedisClient = redisdb
	return redisdb
}
