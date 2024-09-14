package utils

import (
    "github.com/go-redis/redis/v8"
    "os"
)

var redisClient *redis.Client

func ConnectRedis() *redis.Client {
    redisAddr := os.Getenv("REDIS_URL")
    redisClient = redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    return redisClient
}
