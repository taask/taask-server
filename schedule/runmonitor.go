package schedule

import (
	"fmt"
	"time"

	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/model"
)

type runMonitor struct {
	taskUUID       string
	timeoutSeconds int32
	runnerPool     *runnerPool
}

func (sm *Manager) startRunMonitor(taskUUID string, runnerPool *runnerPool) {
	updateChan := sm.updater.GetListener(taskUUID)

	monitor := &runMonitor{
		taskUUID:   taskUUID,
		runnerPool: runnerPool,
	}

	sm.runningLock.Lock()
	sm.running[taskUUID] = monitor
	sm.runningLock.Unlock()

	if task, err := monitor.start(updateChan); err != nil {
		log.LogWarn(err.Error())

		sm.startRetryWorker(task.UUID)
	}

	sm.runningLock.Lock()
	delete(sm.running, taskUUID)
	sm.runningLock.Unlock()
}

func (rm *runMonitor) start(updateChan chan model.Task) (*model.Task, error) {
	log.LogInfo(fmt.Sprintf("starting run monitor for task %s", rm.taskUUID))

	timeoutChan := make(chan time.Time, 1)

	var task model.Task

	for {
		select {
		case task = <-updateChan:
			// continue
		case <-timeoutChan:
			if !task.IsRetrying() && !task.IsFinished() {
				// if the timeout is thrown but the task is already retrying, don't throw an error
				return &task, fmt.Errorf("task %s not completed before timeout", rm.taskUUID)
			}

			continue
		}

		if task.IsFinished() {
			log.LogInfo(fmt.Sprintf("runner %s %s task %s, updating runner tracker", task.Meta.RunnerUUID, task.Status, task.UUID))

			rm.runnerPool.runnerCompletedTask(task.Meta.RunnerUUID, task.UUID)

			return nil, nil
		}

		if task.IsRunning() {
			rm.timeoutSeconds = task.Meta.TimeoutSeconds
			go rm.startTimeout(timeoutChan)
		}

		if task.IsRetrying() {
			log.LogInfo(fmt.Sprintf("task %s began retrying, updating runner %s tracker", task.UUID, task.Meta.RunnerUUID))

			if task.Meta.RunnerUUID != "" {
				rm.runnerPool.runnerCompletedTask(task.Meta.RunnerUUID, task.UUID)
			}
		}
	}
}

func (rm *runMonitor) startTimeout(timeoutChan chan time.Time) {
	<-time.After(time.Second * time.Duration(rm.timeoutSeconds))
	timeoutChan <- time.Now()
}
