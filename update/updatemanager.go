package update

import (
	"fmt"
	"sync"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/metrics"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/storage"
)

// Manager manages task updates
type Manager struct {
	storage   storage.Manager
	listeners map[string]*taskListener
	lock      *sync.Mutex
	metrics   *metrics.Manager
}

type taskListener struct {
	listenerChans [](chan model.Task)
}

// NewManager creates an update manager
func NewManager(storage storage.Manager, metrics *metrics.Manager) *Manager {
	return &Manager{
		storage:   storage,
		listeners: make(map[string]*taskListener),
		lock:      &sync.Mutex{},
		metrics:   metrics,
	}
}

// UpdateTask updates a task in storage and notifies listeners of the new status
func (m *Manager) UpdateTask(update *model.TaskUpdate) {
	if update.UUID == "" {
		log.LogError(errors.New("attempted to update task without providing UUID"))
		return
	}

	task, err := m.storage.Get(update.UUID)
	if err != nil {
		log.LogError(errors.Wrap(err, "failed to storage.Get"))
		return
	}

	go m.metrics.UpdateTask(*task, update)

	// if update is nil, then we just wanted to update metrics
	if update != nil {
		if update.Status != "" && task.Status != update.Status {
			log.LogInfo(fmt.Sprintf("task %s status updated (%s -> %s)", task.UUID, task.Status, update.Status))
			task.Status = update.Status
		}

		if update.RunnerUUID != "" && task.RunnerUUID != update.RunnerUUID {
			log.LogInfo(fmt.Sprintf("task %s assigned to runner %s", task.UUID, update.RunnerUUID))
			task.RunnerUUID = update.RunnerUUID
		}

		if update.RetrySeconds != 0 && task.RetrySeconds != update.RetrySeconds {
			log.LogInfo(fmt.Sprintf("task %s set to retry in %d seconds", task.UUID, update.RetrySeconds))
			task.RetrySeconds = update.RetrySeconds
		}

		// TODO: determine if we should block updating this more than once
		if update.EncResult != nil {
			task.EncResult = update.EncResult
			task.EncResultSymKey = update.EncResultSymKey
		}

		if err := m.storage.Update(*task); err != nil {
			log.LogError(errors.Wrap(err, "failed to m.storage.Update"))
		}

		m.updateListeners(task)
	}
}

// GetListener gets a channel to listen to task updates
func (m *Manager) GetListener(taskUUID string) chan model.Task {
	m.lock.Lock()
	defer m.lock.Unlock()

	var listener *taskListener

	if existing, ok := m.listeners[taskUUID]; ok {
		listener = existing
	} else {
		newListener := &taskListener{
			listenerChans: [](chan model.Task){},
		}

		m.listeners[taskUUID] = newListener

		listener = newListener
	}

	// allow 64 updates to buffer
	listenerChan := make(chan model.Task, 64)

	listener.listenerChans = append(listener.listenerChans, listenerChan)

	return listenerChan
}

func (m *Manager) updateListeners(task *model.Task) {
	m.lock.Lock()
	defer m.lock.Unlock()

	listener, ok := m.listeners[task.UUID]
	if !ok {
		return
	}

	for i := range listener.listenerChans {
		taskCopy := *task

		go func(listenerChan chan model.Task) {
			listenerChan <- taskCopy
		}(listener.listenerChans[i])
	}

	if task.Status == model.TaskStatusCompleted || task.Status == model.TaskStatusFailed {
		log.LogInfo(fmt.Sprintf("task %s completed, removing all update listeners", task.UUID))
		delete(m.listeners, task.UUID)
	}
}
