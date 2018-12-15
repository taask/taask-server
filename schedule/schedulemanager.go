package schedule

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	log "github.com/cohix/simplog"

	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/update"
)

// ErrorNoRunnersRegistered is returned when a task is scheduled to a Kind that has no runners
var ErrorNoRunnersRegistered = errors.New("no runners registered")

// ErrorCapacityReached is returned when a runner is at capacity
var ErrorCapacityReached = errors.New("runner capacity reached")

// Manager manages the scheduling of tasks to runners
type Manager struct {
	// A map of runner Kinds to pools of runners
	runnerPools map[string]*runnerPool

	// scheduleChan is used to schedule tasks
	scheduleChan chan *model.Task

	// Tasks waiting to be assigned to a runner
	queued *list.List

	// A mutex to ensure the queue is kept thread safe
	queueLock *sync.Mutex

	// Delinquient tasks that need to be retried
	retrying map[string]*RetryTaskWorker

	// A mutex to ensure the retry workers are kept thread safe
	retryLock *sync.Mutex

	// updater allows the scheduler to report on updates to the system
	updater *update.Manager
}

// NewManager creates a new ScheduleManager
func NewManager(updater *update.Manager) *Manager {
	return &Manager{
		runnerPools:  make(map[string]*runnerPool),
		scheduleChan: make(chan *model.Task, 256),
		queued:       list.New(),
		queueLock:    &sync.Mutex{},
		retrying:     make(map[string]*RetryTaskWorker),
		retryLock:    &sync.Mutex{},
		updater:      updater,
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
			m.startRetryWorker(nextTask)
			continue
		}

		runner, err := runnerPool.assignTaskToNextRunner(nextTask)
		if err != nil {
			log.LogWarn(errors.Wrap(err, fmt.Sprintf("schedule task %s: no runners of Kind %s available", nextTask.UUID, nextTask.Kind)).Error())
			m.startRetryWorker(nextTask)
			continue
		}

		m.updater.UpdateTask(&model.TaskUpdate{UUID: nextTask.UUID, Status: model.TaskStatusQueued, RunnerUUID: runner.UUID})

		listener := m.updater.GetListener(nextTask.UUID)
		go m.startRunMonitor(nextTask, runnerPool, listener)

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
		m.queueLock.Lock()
		defer m.queueLock.Unlock()

		m.queued.PushBack(task)
	default:
		return
	}
}

func (m *Manager) queueNewTaskUntilExists() {
	task := <-m.scheduleChan

	m.queueLock.Lock()
	defer m.queueLock.Unlock()

	m.queued.PushBack(task)
}

func (m *Manager) nextQueued() *model.Task {
	m.queueLock.Lock()
	defer m.queueLock.Unlock()

	if m.queued.Len() > 0 {
		task := m.queued.Remove(m.queued.Front()).(*model.Task)
		return task
	}

	return nil
}

func (m *Manager) requeueTask(task *model.Task) {
	m.queueLock.Lock()
	defer m.queueLock.Unlock()

	log.LogInfo(fmt.Sprintf("requeing task %s", task.UUID))

	if m.queued.Len() == 0 {
		m.ScheduleTask(task) // if there's nothing, then the run loop will be waiting for a new scheduled task before continuing
		return
	}

	e := m.queued.Front()
	for i := 0; i < m.queued.Len()/3 && e.Next() != nil; i++ {
		e = e.Next()
	}

	m.queued.InsertAfter(task, e)
}

func (m *Manager) forceRetry() {
	m.retryLock.Lock()
	defer m.retryLock.Unlock()

	// range over the map to get a "random" one
	for uuid, worker := range m.retrying {
		log.LogInfo(fmt.Sprintf("releasing retry worker for task %s", uuid))

		go worker.retryNow()
		break
	}
}
