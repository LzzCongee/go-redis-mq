package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

var Client *redis.Client

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 连接本地Redis服务器，端口6379
		DB:   0,                // 使用默认的数据库0
	})
}
