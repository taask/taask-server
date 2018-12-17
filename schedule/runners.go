package schedule

import (
	"container/list"
	"fmt"
	"sync"

	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/model"
)

// RegisterRunner registers the existence of a runner
func (sm *Manager) RegisterRunner(runner *model.Runner) {
	defer log.LogTrace("RegisterRunner")()

	pool, ok := sm.runnerPools[runner.Kind]
	if !ok {
		pool = newRunnerPool()
		sm.runnerPools[runner.Kind] = pool
	}

	go pool.addRunner(runner)
}

// UnregisterRunner removes a runner from eligibility
// this contains some lock-contentious calls, but we're ok if it's slow if it means that we never lose a task
func (sm *Manager) UnregisterRunner(runnerKind, uuid string) error {
	defer log.LogTrace("UnregisterRunner")()

	pool, ok := sm.runnerPools[runnerKind]
	if !ok {
		return fmt.Errorf("no runner pool of Kind %s found", runnerKind)
	}

	deadTasks := pool.removeRunner(uuid)
	if len(deadTasks) > 0 {
		log.LogWarn(fmt.Sprintf("re-queuing %d tasks from unregistered runner %s", len(deadTasks), uuid))

		for i := range deadTasks {
			log.LogInfo(fmt.Sprintf("starting retry worker for dead task %s", deadTasks[i]))
			sm.startRetryWorker(deadTasks[i])
		}
	}

	return nil
}

type runnerPool struct {
	// runners is a map from a UUID to the runner with that UUID
	runners map[string]*model.Runner

	// runnerTracker is a doubly linked list of runners in priority of lowest current load
	runnerTracker *list.List

	// poolLock is a lock for the pool, lock to be applied on all operations that modify `runners` or `runnerTracker`
	poolLock *sync.Mutex
}

type runnerLoad struct {
	// UUID is the uuid of the runner
	UUID string

	// AssignedTasks are the tasks queued and running on a runner
	AssignedTasks map[string]string
}

func (rl *runnerLoad) AssignedCount() int {
	return len(rl.AssignedTasks)
}

func newRunnerPool() *runnerPool {
	rp := &runnerPool{
		runners:       make(map[string]*model.Runner),
		runnerTracker: list.New(),
		poolLock:      &sync.Mutex{},
	}

	return rp
}

func (rp *runnerPool) addRunner(runner *model.Runner) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	rp.runners[runner.UUID] = runner

	runnerTracker := &runnerLoad{
		UUID:          runner.UUID,
		AssignedTasks: make(map[string]string),
	}

	rp.runnerTracker.PushFront(runnerTracker)
}

// returns the tasks that were assigned to the runner when it was removed
func (rp *runnerPool) removeRunner(uuid string) []string {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	deadTasks := rp.removeTracker(uuid)

	delete(rp.runners, uuid)

	return deadTasks
}

func (rp *runnerPool) assignTaskToNextRunner(task *model.Task) (*model.Runner, error) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	elem := rp.runnerTracker.Front()
	if elem == nil {
		return nil, ErrorNoRunnersRegistered
	}

	runnerTracker := elem.Value.(*runnerLoad)

	if runnerTracker.AssignedCount() >= 10 {
		return nil, ErrorCapacityReached
	}

	runner := rp.runners[runnerTracker.UUID]

	runnerTracker.AssignedTasks[task.UUID] = task.UUID

	defer rp.rebalance(elem, runnerTracker.AssignedCount()) // since the unlock was deferred first, the lock will be held until the rebalance is finished

	return runner, nil
}

func (rp *runnerPool) print() {
	fmt.Println("====printing runner pool====")
	for e := rp.runnerTracker.Front(); e != nil; e = e.Next() {
		runnerTracker := e.Value.(*runnerLoad)

		fmt.Printf("runner %s, assigned:\n", runnerTracker.UUID)
		for _, uuid := range runnerTracker.AssignedTasks {
			fmt.Printf("\task %s\n", uuid)
		}
	}
	fmt.Println("====done runner pool====")
}

func (rp *runnerPool) runnerCompletedTask(runnerUUID, taskUUID string) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	for e := rp.runnerTracker.Back(); e != nil; e = e.Prev() {
		runnerTracker := e.Value.(*runnerLoad)

		if runnerTracker.UUID == runnerUUID {
			delete(runnerTracker.AssignedTasks, taskUUID)
			rp.rebalance(e, runnerTracker.AssignedCount())
			return
		}
	}
}

func (rp *runnerPool) rebalance(elem *list.Element, newVal int) {
	if rp.runnerTracker.Len() == 1 {
		return
	}

	for e := rp.runnerTracker.Back(); e != nil; e = e.Prev() {
		if e == elem {
			return
		}

		runnerTracker := e.Value.(*runnerLoad)

		if newVal >= runnerTracker.AssignedCount() {
			rp.runnerTracker.MoveAfter(elem, e)
			return
		}
	}

	if rp.runnerTracker.Front() != elem {
		rp.runnerTracker.MoveBefore(elem, rp.runnerTracker.Front())
	}
}

// removeTracker returns the uuids of the tasks that were assigned to the runner at the time it was removed
func (rp *runnerPool) removeTracker(uuid string) []string {
	tasks := []string{}

	for e := rp.runnerTracker.Back(); e != nil; e = e.Prev() {
		runnerTracker := e.Value.(*runnerLoad)

		if runnerTracker.UUID == uuid {
			for uuid := range runnerTracker.AssignedTasks {
				tasks = append(tasks, runnerTracker.AssignedTasks[uuid])
			}

			rp.runnerTracker.Remove(e)
			break
		}
	}

	return tasks
}
