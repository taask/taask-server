package schedule

import (
	"container/list"
	"errors"
	"fmt"

	log "github.com/cohix/simplog"

	"github.com/taask/taask-server/model"
)

// ErrorNoRunnersRegistered is returned when a task is scheduled to a Kind that has no runners
var ErrorNoRunnersRegistered = errors.New("no runners registered")

// Manager manages the scheduling of tasks to runners
type Manager struct {
	// A map of runner Kinds to pools of runners
	runnerPools map[string]*runnerPool

	// scheduleChan is used to schedule tasks
	scheduleChan chan *model.Task

	// Tasks waiting to be assigned to a runner
	queued *list.List

	// Delinquient tasks that need to be retried
	retrying []*RetryTaskWorker
}

// NewManager creates a new ScheduleManager
func NewManager() *Manager {
	return &Manager{
		runnerPools:  make(map[string]*runnerPool),
		scheduleChan: make(chan *model.Task, 256),
		queued:       list.New(),
		retrying:     []*RetryTaskWorker{},
	}
}

// Start begins the scheduler
func (m *Manager) Start() {
	defer log.LogTrace("schedule.Manager.Start()")()

	for {
		m.queueNewTaskIfExists()

		nextTask := m.nextQueued()
		if nextTask == nil {
			m.queueNewTaskUntilExists()
			continue
		}

		runnerPool, ok := m.runnerPools[nextTask.Kind]
		if !ok {
			log.LogWarn(fmt.Sprintf("schedule task %s: no runners of Kind %s registered", nextTask.UUID, nextTask.Kind))
			m.requeueTask(nextTask)
			continue
		}

		runner, err := runnerPool.nextRunner()
		if err != nil {
			log.LogWarn(fmt.Sprintf("schedule task %s: no runners of Kind %s registered", nextTask.UUID, nextTask.Kind))
			m.requeueTask(nextTask)
			continue
		} else {
			log.LogInfo(fmt.Sprintf("scheduling task %s to runner %s", nextTask.UUID, runner.UUID))
		}

		runner.TaskChannel <- nextTask
	}
}

// ScheduleTask schedules a task
func (m *Manager) ScheduleTask(task *model.Task) {
	defer log.LogTrace(fmt.Sprintf("ScheduleTask %s", task.UUID))()

	m.scheduleChan <- task
}

// TODO: determine if this should flush the channel or not
func (m *Manager) queueNewTaskIfExists() {
	select {
	case task := <-m.scheduleChan:
		m.queued.PushBack(task)
	default:
		return
	}
}

func (m *Manager) queueNewTaskUntilExists() {
	task := <-m.scheduleChan
	m.queued.PushBack(task)
}

func (m *Manager) nextQueued() *model.Task {
	if m.queued.Len() > 0 {
		task := m.queued.Remove(m.queued.Front()).(*model.Task)
		return task
	}

	return nil
}

func (m *Manager) requeueTask(task *model.Task) {
	e := m.queued.Front()
	for i := 0; i < m.queued.Len()/3 && e.Next() != nil; i++ {
		e = e.Next()
	}

	m.queued.InsertAfter(task, e)
}
