package main

import (
	"fmt"
	"time"

	"go-redis-mq/internal/redisclient"
	"go-redis-mq/internal/task"

	"github.com/redis/go-redis/v9"
)

const (
	retryStream = "retry_stream" // 重试任务流
	mainStream  = "task_stream"  // 主任务流
)

func main() {
	redisclient.InitRedis()
	rdb := redisclient.Client
	fmt.Println("🚀 Retry worker started...")

	for {
		streams, err := rdb.XRead(redisclient.Ctx, &redis.XReadArgs{
			Streams: []string{retryStream, "0"},
			Count:   1,               // 每次最多读取1条消息
			Block:   3 * time.Second, // 如果没有消息，最多阻塞3秒
		}).Result()

		if err != nil && err != redis.Nil {
			fmt.Println("❌ Error reading retry stream:", err)
			continue
		}
		if len(streams) == 0 {
			// 只有当流中有数据时才会继续执行
			continue
		}

		for _, msg := range streams[0].Messages {
			fmt.Println("🔁 Retrying task:", msg.Values["task_id"])

			retryCount := msg.Values["retry_count"]
			if !task.CanRetry(retryCount) {
				// 检查重试次数是否超限
				fmt.Println("⚠️  Retry limit exceeded for task:", msg.Values["task_id"])
				continue
			}

			// Optional: exponential backoff or delay 添加延迟（可选的退避策略）
			time.Sleep(1 * time.Second)

			// Re-publish to main stream 将任务重新发布到主流
			_, err := rdb.XAdd(redisclient.Ctx, &redis.XAddArgs{
				Stream: mainStream,
				Values: map[string]interface{}{
					"task_id":     msg.Values["task_id"],
					"payload":     msg.Values["payload"],
					"created_at":  msg.Values["created_at"],
					"retry_count": retryCount,
					"from_retry":  true,
				},
			}).Result()

			if err != nil {
				fmt.Println("❌ Retry publish failed:", err)
			} else {
				fmt.Println("✅ Retry sent back to stream:", msg.Values["task_id"])
			}

			// 删除 retry_stream 中的任务（避免重复处理）
			rdb.XDel(redisclient.Ctx, retryStream, msg.ID)
		}
	}
}
