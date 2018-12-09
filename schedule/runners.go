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
func (m *Manager) UnregisterRunner(runnerKind, uuid string) error {
	defer log.LogTrace("UnregisterRunner")()

	pool, ok := m.runnerPools[runnerKind]
	if !ok {
		return fmt.Errorf("no runner pool of Kind %s found", runnerKind)
	}

	go pool.removeRunner(uuid)

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
	AssignedCount int
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
		AssignedCount: 0,
	}

	rp.tracker.PushFront(tracker)
}

func (rp *runnerPool) removeRunner(uuid string) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	delete(rp.runners, uuid)

	rp.removeTracker(uuid)
}

func (rp *runnerPool) nextRunner() (*model.Runner, error) {
	rp.poolLock.Lock()
	defer rp.poolLock.Unlock()

	trackerElement := rp.tracker.Front()
	if trackerElement == nil {
		return nil, ErrorNoRunnersRegistered
	}

	tracker := trackerElement.Value.(*runnerLoad)

	runner := rp.runners[tracker.UUID]

	tracker.AssignedCount++
	defer rp.rebalance(trackerElement, tracker.AssignedCount) // since the unlock was deferred first, the lock will be held until the rebalance is finished

	return runner, nil
}

func (rp *runnerPool) rebalance(elem *list.Element, newVal int) {
	if rp.tracker.Len() == 1 {
		return
	}

	for e := rp.tracker.Back(); e != nil; e = e.Next() {
		if e == elem {
			break
		}

		tracker := e.Value.(*runnerLoad)

		if newVal >= tracker.AssignedCount {
			rp.tracker.MoveBefore(elem, e)
			break
		}
	}
}

func (rp *runnerPool) removeTracker(uuid string) {
	for e := rp.tracker.Back(); e != nil; e = e.Next() {
		tracker := e.Value.(*runnerLoad)

		if tracker.UUID == uuid {
			rp.tracker.Remove(e)
			break
		}
	}
}
