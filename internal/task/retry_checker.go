package task

import (
	"fmt"
	"strconv"
)

const MaxRetry = 3

// 重试策略
func CanRetry(retryCount any) bool {
	count, err := strconv.Atoi(fmt.Sprint(retryCount))
	if err != nil {
		return false
	}
	return count < MaxRetry
}
