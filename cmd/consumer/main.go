package main

import (
	"fmt"
	"strings"
	"time"

	"go-redis-mq/internal/redisclient"
	"go-redis-mq/internal/retry"

	"github.com/redis/go-redis/v9"
)

const (
	stream      = "task_stream"
	group       = "task_group"
	consumer    = "worker-1" // 可修改为 worker-2, worker-N
	pollTimeout = 5 * time.Second
)

// 确保消费者组存在
func ensureGroupExists() {
	rdb := redisclient.Client
	err := rdb.XGroupCreateMkStream(redisclient.Ctx, stream, group, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		// 只有真正的错误（比如连接失败、权限问题等）才会导致程序终止，而不是因为 BUSYGROUP 错误。
		// 对于 BUSYGROUP 错误，我们只是简单地忽略它，因为它表示消费者组已经存在，不需要重复创建。
		panic(err)
	}
	fmt.Println("Consumer group ensured.")
}

func main() {
	redisclient.InitRedis()
	ensureGroupExists()
	rdb := redisclient.Client

	for { // 持续监听主任务流
		streams, err := rdb.XReadGroup(redisclient.Ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Block:    pollTimeout,
			Count:    1,
		}).Result()

		if err != nil && err != redis.Nil {
			fmt.Println("Read error:", err)
			continue
		}

		if len(streams) == 0 {
			continue
		}

		for _, msg := range streams[0].Messages {
			fmt.Println("🟢 Received task:", msg.Values["task_id"], msg.Values["payload"])

			// 假设任务处理函数
			err := handleTask(msg.Values["payload"].(string))
			if err != nil {
				retry.HandleFailedTask(msg, err)
			} else {
				rdb.XAck(redisclient.Ctx, stream, group, msg.ID)
				rdb.Incr(redisclient.Ctx, "task_success_count")
				fmt.Println("✅ ACK task:", msg.ID)
			}
		}
	}
}

// 模拟任务处理逻辑
func handleTask(payload string) error {
	fmt.Println("Processing:", payload)
	time.Sleep(1 * time.Second)

	if payload == "fail" {
		return fmt.Errorf("simulated failure")
	}
	return nil
}
