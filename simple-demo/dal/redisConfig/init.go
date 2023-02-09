package redisConfig

import (
	"github.com/go-redis/redis"
	"log"
)

var RD *redis.Client

func Init() {
	RD = redis.NewClient(&redis.Options{
		Addr:     "43.143.166.162:6379",
		Password: "root",
		DB:       1,
	})

	if res, err := RD.Ping().Result(); err != nil {
		log.Default().Println("ping err:", err.Error())
		return
	} else {
		log.Default().Println(res)
	}
}
