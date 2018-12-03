package schedule

import (
	"time"

	"github.com/taask/taask-server/model"
)

// RetryTaskWorker manager the backoff retries of a delenquient task
type RetryTaskWorker struct {
	task      *model.Task
	retryTime time.Time
}
