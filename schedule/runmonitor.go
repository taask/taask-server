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

func (sm *Manager) startRunMonitor(task *model.Task, runnerPool *runnerPool, updateChan chan model.Task) {
	monitor := &runMonitor{
		taskUUID:       task.UUID,
		timeoutSeconds: task.Meta.TimeoutSeconds,
		runnerPool:     runnerPool,
	}

	if err := monitor.start(updateChan); err != nil {
		sm.startRetryWorker(task)
	}
}

func (rm *runMonitor) start(updateChan chan model.Task) error {
	log.LogInfo(fmt.Sprintf("starting run monitor for task %s", rm.taskUUID))

	timeoutChan := make(chan time.Time)

	for {
		var task model.Task
		select {
		case task = <-updateChan:
			// continue
		case <-timeoutChan:
			return fmt.Errorf("task %s completion not reached before timeout", rm.taskUUID)
		}

		if task.IsFinished() {
			log.LogInfo(fmt.Sprintf("runner %s %s task %s, updating runner tracker", task.Meta.RunnerUUID, task.Status, task.UUID))

			rm.runnerPool.runnerCompletedTask(task.Meta.RunnerUUID, task.UUID)

			if task.Status != model.TaskStatusCompleted {
				return fmt.Errorf("task %s reported in state %s", task.UUID, task.Status)
			}

			return nil
		}

		if task.IsRunning() {
			go rm.startTimeout(timeoutChan)
		}
	}
}

func (rm *runMonitor) startTimeout(timeoutChan chan time.Time) {
	<-time.After(time.Second * time.Duration(rm.timeoutSeconds))
	timeoutChan <- time.Now()
}
