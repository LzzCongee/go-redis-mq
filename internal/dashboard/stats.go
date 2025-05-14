package dashboard

import (
	"go-redis-mq/internal/redisclient"
)

// Redis 任务统计逻辑
// 获取Redis中任务队列的长度，以及消费组中待处理的消息数量
func GetStats() (map[string]interface{}, error) {
	ctx := redisclient.Ctx
	rdb := redisclient.Client

	taskLen, _ := rdb.XLen(ctx, "task_stream").Result()
	retryLen, _ := rdb.XLen(ctx, "retry_stream").Result()

	// 使用Redis的XPending命令来获取消费组中待处理的消息信息
	// task_stream ：要查询的流名称
	// task_group ：消费组的名称
	// 返回的pending对象包含了该消费组中所有待处理消息的统计信息，如待处理消息的总数（Count）
	pending, _ := rdb.XPending(ctx, "task_stream", "task_group").Result()

	successCount, _ := rdb.Get(ctx, "task_success_count").Int64()

	return map[string]interface{}{
		"task_stream_len":    taskLen,
		"retry_stream_len":   retryLen,
		"pending_tasks":      pending.Count,
		"success_task_count": successCount,
	}, nil
}
