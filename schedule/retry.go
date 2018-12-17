package schedule

import (
	"fmt"
	"time"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
)

// retryTaskWorker manager the backoff retries of a delenquient task
type retryTaskWorker struct {
	task    *model.Task
	nowChan chan bool
}

// StartRetryWorker starts a retryWorker for a task
func (sm *Manager) StartRetryWorker(taskUUID string) {
	// sm.retryLock.Lock()
	// _, exists := sm.retrying[taskUUID]
	// sm.retryLock.Unlock()

	// if exists {
	// 	log.LogInfo(fmt.Sprintf("retry worker for task %s already exists, canceling", taskUUID))
	// 	return
	// }

	listener := sm.updater.GetListener(taskUUID)
	task := <-listener // just get the first update, which will be the task's state when the listener was created

	retrySeconds := task.Meta.RetrySeconds

	if retrySeconds == 0 {
		retrySeconds = 1
	} else if retrySeconds < 16 {
		retrySeconds *= 2
	}

	worker := &retryTaskWorker{
		task:    &task,
		nowChan: make(chan bool),
	}

	sm.retryLock.Lock()
	sm.retrying[taskUUID] = worker
	sm.retryLock.Unlock()

	update, err := task.Update(model.TaskUpdate{Status: model.TaskStatusRetrying, RetrySeconds: retrySeconds, RunnerUUID: ""})
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

		log.LogInfo(fmt.Sprintf("task %s retrying after %d seconds", taskUUID, retrySeconds))

		sm.ScheduleTask(&task)

		sm.retryLock.Lock()
		delete(sm.retrying, taskUUID)
		sm.retryLock.Unlock()
	}()
}

func (rw *retryTaskWorker) retryNow() {
	go func() {
		rw.nowChan <- true
	}()
}
