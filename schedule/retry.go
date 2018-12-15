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

func (sm *Manager) startRetryWorker(task *model.Task) {
	if task.Meta.RetrySeconds == 0 {
		task.Meta.RetrySeconds = 1
	} else if task.Meta.RetrySeconds < 16 {
		task.Meta.RetrySeconds *= 2
	}

	worker := &RetryTaskWorker{
		task:    task,
		nowChan: make(chan bool),
	}

	sm.retryLock.Lock()
	sm.retrying[task.UUID] = worker
	sm.retryLock.Unlock()

	sm.updater.UpdateTask(&model.TaskUpdate{UUID: task.UUID, Status: model.TaskStatusRetrying, RetrySeconds: task.Meta.RetrySeconds})

	go func() {
		select {
		case <-time.After(time.Second * time.Duration(worker.task.Meta.RetrySeconds)):
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
