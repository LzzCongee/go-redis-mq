package redisclient

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client redis.UniversalClient

func InitRedis() {
	// 加载 .env 文件（仅第一次初始化时调用）
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, falling back to system environment")
	}

	// 从环境变量读取模式： "single" | "cluster" | "sentinel"
	mode := os.Getenv("REDIS_MODE")
	addrs := strings.Split(os.Getenv("REDIS_ADDRS"), ",") // e.g. "10.0.0.1:6379,10.0.0.2:6379"
	password := os.Getenv("REDIS_PASSWORD")

	switch mode {
	case "cluster":
		// Cluster 模式
		Client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: password,
		})
	case "sentinel":
		// Sentinel 模式
		Client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    os.Getenv("REDIS_MASTER_NAME"), // e.g. "mymaster"
			SentinelAddrs: addrs,
			Password:      password,
		})
	default:
		// 单机模式
		Client = redis.NewClient(&redis.Options{
			Addr:     addrs[0],
			Password: password,
		})
	}
}
