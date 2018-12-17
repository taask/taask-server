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

	// Running tasks that are being tracked
	running map[string]*runMonitor
	// A mutex to ensure the run monitors are kept thread safe
	runningLock *sync.Mutex

	// Delinquient tasks that need to be retried
	retrying map[string]*retryTaskWorker
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
		running:      make(map[string]*runMonitor),
		runningLock:  &sync.Mutex{},
		retrying:     make(map[string]*retryTaskWorker),
		retryLock:    &sync.Mutex{},
		updater:      updater,
	}
}

// Start begins the scheduler
func (sm *Manager) Start() {
	defer log.LogTrace("schedule.Manager.Start()")()

	for {
		sm.queueNewTaskIfExists()

		nextTask := sm.nextQueued()
		if nextTask == nil {
			sm.queueNewTaskUntilExists()
			continue
		}

		sm.runningLock.Lock()
		_, exists := sm.running[nextTask.UUID]
		sm.runningLock.Unlock()

		if exists && !nextTask.IsRetrying() {
			log.LogWarn(fmt.Sprintf("attempted to schedule task %s in state %s that has already been scheduled", nextTask.UUID, nextTask.Status))
			continue
		}

		runnerPool, ok := sm.runnerPools[nextTask.Kind]
		if !ok {
			log.LogWarn(fmt.Sprintf("schedule task %s: no runners of Kind %s registered", nextTask.UUID, nextTask.Kind))
			sm.startRetryWorker(nextTask.UUID)
			continue
		}

		runner, err := runnerPool.assignTaskToNextRunner(nextTask)
		if err != nil {
			log.LogWarn(errors.Wrap(err, fmt.Sprintf("schedule task %s: no runners of Kind %s available", nextTask.UUID, nextTask.Kind)).Error())
			sm.startRetryWorker(nextTask.UUID)
			continue
		}

		go sm.startRunMonitor(nextTask.UUID, runnerPool)

		runner.TaskChannel <- nextTask
	}
}

// ScheduleTask schedules a task
func (sm *Manager) ScheduleTask(task *model.Task) {
	defer log.LogTrace(fmt.Sprintf("ScheduleTask %s", task.UUID))()

	sm.scheduleChan <- task
}

// TODO: determine if this should flush the channel or not
func (sm *Manager) queueNewTaskIfExists() {
	select {
	case task := <-sm.scheduleChan:
		sm.queueLock.Lock()
		defer sm.queueLock.Unlock()

		sm.queued.PushBack(task)
	default:
		return
	}
}

func (sm *Manager) queueNewTaskUntilExists() {
	task := <-sm.scheduleChan

	sm.queueLock.Lock()
	defer sm.queueLock.Unlock()

	sm.queued.PushBack(task)
}

func (sm *Manager) nextQueued() *model.Task {
	sm.queueLock.Lock()
	defer sm.queueLock.Unlock()

	if sm.queued.Len() > 0 {
		task := sm.queued.Remove(sm.queued.Front()).(*model.Task)
		return task
	}

	return nil
}

func (sm *Manager) forceRetry() {
	sm.retryLock.Lock()
	defer sm.retryLock.Unlock()

	// range over the map to get a "random" one
	for uuid, worker := range sm.retrying {
		log.LogInfo(fmt.Sprintf("releasing retry worker for task %s", uuid))

		go worker.retryNow()
		break
	}
}
