package schedule

import (
	"fmt"
	"time"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
)

// RetryTaskWorker manager the backoff retries of a delenquient task
type RetryTaskWorker struct {
	task    *model.Task
	nowChan chan bool
}

func (sm *Manager) startRetryWorker(task *model.Task) {
	retrySeconds := task.Meta.RetrySeconds

	if retrySeconds == 0 {
		retrySeconds = 1
	} else if retrySeconds < 16 {
		retrySeconds *= 2
	}

	worker := &RetryTaskWorker{
		task:    task,
		nowChan: make(chan bool),
	}

	sm.retryLock.Lock()
	sm.retrying[task.UUID] = worker
	sm.retryLock.Unlock()

	update, err := task.Update(model.TaskUpdate{Status: model.TaskStatusRetrying, RetrySeconds: retrySeconds})
	if err != nil {
		log.LogWarn(errors.Wrap(err, "startRetryWorker failed to task.Update").Error())
	}

	sm.updater.UpdateTask(update)

	go func() {
		select {
		case <-time.After(time.Second * time.Duration(retrySeconds)):
			// run the task
		case <-worker.nowChan:
			// ignore the timer and run the task
		}

		log.LogInfo(fmt.Sprintf("task %s retrying after %d seconds", task.UUID, retrySeconds))

		task.Meta.RetrySeconds = retrySeconds
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
