package main

import (
	"fmt"
	"os"
	"time"

	"go-redis-mq/internal/redisclient"
	"go-redis-mq/internal/task"

	"github.com/redis/go-redis/v9"
)

// 使用 XADD 写入 Redis Stream task_stream
// 每条任务包括字段：task_id、payload、created_at
// 支持 CLI 参数自定义任务内容（可选）
// 自动重连、失败重试

// go run cmd/producer/main.go "process-image-123"

func main() {
	redisclient.InitRedis()
	rdb := redisclient.Client // 获得全局的客户端实例

	payload := "default-task"
	if len(os.Args) > 1 {
		payload = os.Args[1]
	}

	t := task.New(payload)
	fmt.Println("Producing task:", t.ID, t.Payload)

	// Result() 函数返回两个值：
	// 1. 第一个返回值（被代码中用 _ 忽略了）：是一个字符串类型，表示新添加的消息的 ID。在 Redis Stream 中，每条消息都会有一个唯一的 ID。
	// 2. 第二个返回值 err ：是一个 error 类型，用于表示操作是否成功。如果操作成功，err 为 nil；如果操作失败，err 会包含错误信息。
	_, err := rdb.XAdd(redisclient.Ctx, &redis.XAddArgs{
		Stream: "task_stream",
		Values: map[string]interface{}{
			"task_id":    t.ID,
			"payload":    t.Payload,
			"created_at": t.CreatedAt.Format(time.RFC3339Nano),
		},
	}).Result()

	if err != nil {
		panic(err)
	}

	fmt.Println("Task sent to Redis stream successfully.")
}
