package redis

import (
	"clover/pkg/log"
	"clover/setting"
	"context"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis(conf *setting.RedisConf) {
	rd := redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
	})

	_, err := rd.Ping(context.Background()).Result()
	if err != nil {
		log.WithCategory("redis").WithError(err).Error("InitRedis: ping failed")
		panic(err)
	}

	redisClient = rd
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	err := redisClient.Close()
	if err != nil {
		log.WithCategory("redis").WithError(err).Error("CloseRedis: failed")
	}
}
