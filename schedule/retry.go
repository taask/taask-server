package schedule

import (
	"time"

	"github.com/taask/taask-server/model"
)

// RetryTaskWorker manager the backoff retries of a delenquient task
type RetryTaskWorker struct {
	task    *model.Task
	nowChan chan bool
}

// StartRetryWorker starts a retry worker
func (sm *Manager) StartRetryWorker(task *model.Task) {
	if task.RetrySeconds == 0 {
		task.RetrySeconds = 1
	} else if task.RetrySeconds < 16 {
		task.RetrySeconds *= 2
	}

	worker := &RetryTaskWorker{
		task:    task,
		nowChan: make(chan bool),
	}

	sm.retryLock.Lock()
	sm.retrying[task.UUID] = worker
	sm.retryLock.Unlock()

	sm.updater.UpdateTask(&model.TaskUpdate{UUID: task.UUID, Status: model.TaskStatusRetrying, RetrySeconds: task.RetrySeconds})

	go func() {
		select {
		case <-time.After(time.Second * time.Duration(worker.task.RetrySeconds)):
			// run the task
		case <-worker.nowChan:
			// ignore the timer and run the task
		}

		sm.requeueTask(task)

		sm.retryLock.Lock()
		delete(sm.retrying, task.UUID)
		sm.retryLock.Unlock()
	}()
}

func (rw *RetryTaskWorker) retryNow() {
	go func() {
		rw.nowChan <- true
	}()
}
