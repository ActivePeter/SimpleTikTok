package redisConfig

import (
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/go-redis/redis"
	"log"
)

var RD *redis.Client

func Init(config *utils.ServerConfig) {
	RD = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPw,
		DB:       1,
	})

	if res, err := RD.Ping().Result(); err != nil {
		log.Default().Println("ping err:", err.Error())
		return
	} else {
		log.Default().Println(res)
	}
}
