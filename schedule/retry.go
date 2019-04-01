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

func (sm *Manager) startRetryWorker(taskUUID string) {
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

	update := task.BuildUpdate(model.TaskChanges{Status: model.TaskStatusRetrying, RetrySeconds: retrySeconds, RunnerUUID: ""})

	_, err := sm.updater.UpdateTask(update)
	if err != nil {
		log.LogError(errors.Wrapf(err, "retryWorker failed to UpdateTask for task %s", task.UUID))
	}

	go func() {
		for {
			select {
			case task = <-listener:
				if task.IsRunning() || task.IsFinished() {
					log.LogInfo(fmt.Sprintf("retry monitor for task %s detected task with staus %s, canceling", task.UUID, task.Status))

					sm.retryLock.Lock()
					delete(sm.retrying, taskUUID)
					sm.retryLock.Unlock()

					break
				}

				continue
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

			break
		}
	}()
}

func (rw *retryTaskWorker) retryNow() {
	go func() {
		rw.nowChan <- true
	}()
}
