package schedule

import (
	"container/list"
	"sync"

	"github.com/taask/taask-server/model"
)

// RegisterRunner registers the existence of a runner
func (m *Manager) RegisterRunner(runner *model.Runner) {
	pool, ok := m.runnerPools[runner.Type]
	if !ok {
		pool = newRunnerPool()
		m.runnerPools[runner.Type] = pool
	}

	go pool.addRunner(runner)
}

// UnregisterRunner removes a runner from eligibility
func (m *Manager) UnregisterRunner(runnerType, uuid string) error {
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
	defer rp.rebalance(trackerElement, tracker.AssignedCount)

	return runner, nil
}

func (rp *runnerPool) rebalance(elem *list.Element, newVal int) {
	for e := rp.tracker.Back(); e != nil; e = e.Next() {
		tracker := e.Value.(*runnerLoad)

		if newVal >= tracker.AssignedCount {
			rp.tracker.MoveBefore(elem, e)
			break
		}
	}
}
