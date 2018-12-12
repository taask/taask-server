package schedule

import (
	"time"

	"github.com/taask/taask-server/model"
)

// RetryTaskWorker manager the backoff retries of a delenquient task
type RetryTaskWorker struct {
	task *model.Task
}

// StartRetryWorker starts a retry worker
func (sm *Manager) StartRetryWorker(task *model.Task) {
	if task.RetrySeconds == 0 {
		task.RetrySeconds = 1
	} else if task.RetrySeconds < 32 {
		task.RetrySeconds *= 2
	}

	worker := &RetryTaskWorker{
		task: task,
	}

	sm.retrying[task.UUID] = worker

	sm.updater.UpdateTask(&model.TaskUpdate{UUID: task.UUID, Status: model.TaskStatusRetrying, RetrySeconds: task.RetrySeconds})

	go func() {
		<-time.After(time.Second * time.Duration(worker.task.RetrySeconds))

		sm.requeueTask(task)

		delete(sm.retrying, task.UUID)
	}()
}
