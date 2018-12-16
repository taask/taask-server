package schedule

import (
	"fmt"
	"time"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
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
		log.LogWarn(err.Error())

		update, err := task.Update(model.TaskUpdate{Status: model.TaskStatusFailed, RunnerUUID: ""})
		if err != nil {
			log.LogWarn(errors.Wrap(err, "startRetryWorker failed to task.Update").Error())
		}

		sm.updater.UpdateTask(update)
		sm.startRetryWorker(task)
	}
}

func (rm *runMonitor) start(updateChan chan model.Task) error {
	log.LogInfo(fmt.Sprintf("starting run monitor for task %s", rm.taskUUID))

	timeoutChan := make(chan time.Time, 1)

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
