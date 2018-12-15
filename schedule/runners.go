package schedule

import (
	"container/list"
	"fmt"
	"sync"

	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/model"
)

// RegisterRunner registers the existence of a runner
func (m *Manager) RegisterRunner(runner *model.Runner) {
	defer log.LogTrace("RegisterRunner")()

	pool, ok := m.runnerPools[runner.Kind]
	if !ok {
		pool = newRunnerPool()
		m.runnerPools[runner.Kind] = pool
	}

	go pool.addRunner(runner)
}

// UnregisterRunner removes a runner from eligibility
// this contains some lock-contentious calls, but we're ok if it's slow if it means that we never lose a task
func (m *Manager) UnregisterRunner(runnerKind, uuid string) error {
	defer log.LogTrace("UnregisterRunner")()

	pool, ok := m.runnerPools[runnerKind]
	if !ok {
		return fmt.Errorf("no runner pool of Kind %s found", runnerKind)
	}

	deadTasks := pool.removeRunner(uuid)
	if len(deadTasks) > 0 {
		log.LogWarn(fmt.Sprintf("re-queuing %d tasks from unregistered runner %s", len(deadTasks), uuid))

		for i := range deadTasks {
			m.updater.UpdateTask(&model.TaskUpdate{UUID: deadTasks[i].UUID, Status: model.TaskStatusFailed, RetrySeconds: deadTasks[i].Meta.RetrySeconds})
			m.startRetryWorker(deadTasks[i])
		}
	}

	return nil
}

type runnerPool struct {
	// runners is a map from a UUID to the runner with that UUID
	runners map[string]*model.Runner

	// tracker is a doubly linked list of runners in priority of lowest current load
	tracker *list.List

	// poolLock is a lock for the pool, lock to be applied on all operations that modify `runners` or `tracker`
	poolLock *sync.Mutex
}

type runnerLoad struct {
	// UUID is the uuid of the runner
	UUID string

	// AssignedCount is the number of tasks queued and running on a runner
	AssignedTasks map[string]*model.Task
}

func (rl *runnerLoad) AssignedCount() int {
	return len(rl.AssignedTasks)
}

func newRunnerPool() *runnerPool {
	rp := &runnerPool{
		runners:  make(map[string]*model.Runner),
		tracker:  list.New(),
		poolLock: &sync.Mutex{},
	}

	return rp
}

func (rp *runnerPool) addRunner(runner *model.Runner) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	rp.runners[runner.UUID] = runner

	tracker := &runnerLoad{
		UUID:          runner.UUID,
		AssignedTasks: make(map[string]*model.Task),
	}

	rp.tracker.PushFront(tracker)
}

func (rp *runnerPool) removeRunner(uuid string) []*model.Task {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	delete(rp.runners, uuid)

	return rp.removeTracker(uuid)
}

func (rp *runnerPool) assignTaskToNextRunner(task *model.Task) (*model.Runner, error) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	trackerElement := rp.tracker.Front()
	if trackerElement == nil {
		return nil, ErrorNoRunnersRegistered
	}

	tracker := trackerElement.Value.(*runnerLoad)

	if tracker.AssignedCount() >= 10 {
		return nil, ErrorCapacityReached
	}

	runner := rp.runners[tracker.UUID]

	tracker.AssignedTasks[task.UUID] = task

	defer rp.rebalance(trackerElement, tracker.AssignedCount()) // since the unlock was deferred first, the lock will be held until the rebalance is finished

	return runner, nil
}

func (rp *runnerPool) runnerCompletedTask(runnerUUID, taskUUID string) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	for e := rp.tracker.Back(); e != nil; e = e.Prev() {
		tracker := e.Value.(*runnerLoad)

		if tracker.UUID == runnerUUID {
			delete(tracker.AssignedTasks, taskUUID)
			rp.rebalance(e, tracker.AssignedCount())
			return
		}
	}
}

func (rp *runnerPool) rebalance(elem *list.Element, newVal int) {
	if rp.tracker.Len() == 1 {
		return
	}

	for e := rp.tracker.Back(); e != nil; e = e.Prev() {
		if e == elem {
			return
		}

		tracker := e.Value.(*runnerLoad)

		if newVal >= tracker.AssignedCount() {
			rp.tracker.MoveAfter(elem, e)
			return
		}
	}

	if rp.tracker.Front() != elem {
		rp.tracker.MoveBefore(elem, rp.tracker.Front())
	}
}

// removeTracker returns the tasks that were assigned to the runner at the time it was removed
func (rp *runnerPool) removeTracker(uuid string) []*model.Task {
	tasks := []*model.Task{}

	for e := rp.tracker.Back(); e != nil; e = e.Prev() {
		tracker := e.Value.(*runnerLoad)

		for uuid := range tracker.AssignedTasks {
			tasks = append(tasks, tracker.AssignedTasks[uuid])
		}

		if tracker.UUID == uuid {
			rp.tracker.Remove(e)
			break
		}
	}

	return tasks
}
