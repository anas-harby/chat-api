package redis

import (
	"log"
	"github.com/go-redis/redis"
	"github.com/bsm/redislock"
	"github.com/anas-harby/go-chat-creation-api/configs"
)

var redisClient *redis.Client
var redisLocker *redislock.Client

func GetRedis() (*redis.Client, error) {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     configs.RedisAddress,
			Password: "",
			DB:       0,
		})

		err := redisClient.Ping().Err()
		if err != nil {
			return nil, err
		}
		redisLocker = redislock.New(redisClient)
	}
	return redisClient, nil
}

func GetLocker() (*redislock.Client) {
	if redisClient == nil {
		log.Fatalln("Redis client is not initialized yet")
	}
	return redisLocker
}