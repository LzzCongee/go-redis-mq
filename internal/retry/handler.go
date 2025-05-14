package retry

import (
	"context"
	"fmt"

	"go-redis-mq/internal/redisclient"

	"github.com/redis/go-redis/v9"
)

func HandleFailedTask(msg redis.XMessage, err error) {
	fmt.Println("❌ Task failed:", msg.ID, "reason:", err)

	retryCount := 0
	if val, ok := msg.Values["retry_count"]; ok {
		// Sscanf函数用于从字符串中解析格式化的输入。它类似于scanf函数，但是它从字符串读取而不是从标准输入。
		// 第一个参数是Sprint函数返回的字符串，即val的字符串表示。
		// 第二个参数"%d"是一个格式字符串，指定了期望解析的数据类型。%d表示期望解析一个整数。
		// 第三个参数&retryCount是一个指向整数的指针，用于存储解析后的整数值。retryCount应该是一个已经在其他地方声明过的整数变量。
		fmt.Sscanf(fmt.Sprint(val), "%d", &retryCount)
	}
	retryCount++

	if retryCount > 3 {
		fmt.Println("⚠️  Discarding task after 3 retries:", msg.ID)
		return
	}

	// 写入 retry_stream
	_, e := redisclient.Client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "retry_stream",
		Values: map[string]interface{}{
			"task_id":     msg.Values["task_id"],
			"payload":     msg.Values["payload"],
			"created_at":  msg.Values["created_at"],
			"retry_count": retryCount,
			"last_error":  err.Error(),
		},
	}).Result()

	if e != nil {
		fmt.Println("Failed to send to retry_stream:", e)
	}
}
