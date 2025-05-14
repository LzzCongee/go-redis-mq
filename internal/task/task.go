package task

import (
	"fmt"
	"time"
)

type Task struct {
	ID        string
	Payload   string
	CreatedAt time.Time
}

func New(payload string) *Task {
	return &Task{
		ID:        fmt.Sprintf("task-%d", time.Now().UnixNano()),
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}
