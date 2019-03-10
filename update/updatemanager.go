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
func (m *Manager) UpdateTask(update model.TaskUpdate) *model.Task {
	if update.UUID == "" {
		log.LogError(errors.New("attempted to update task without providing UUID"))
		return nil
	}

	task, err := m.storage.Get(update.UUID)
	if err != nil {
		log.LogError(errors.Wrap(err, "failed to storage.Get"))
		return nil
	}

	go m.metrics.UpdateTask(*task, update)

	if err := task.ApplyUpdate(update, true); err != nil {
		log.LogWarn(errors.Wrap(err, "update.Manager failed to ApplyUpdate").Error())
		return nil
	}

	if err := m.storage.Update(*task); err != nil {
		log.LogError(errors.Wrap(err, "failed to m.storage.Update"))
		return nil
	}

	m.updateListeners(task)

	return task
}

// GetListener gets a channel to listen to task updates, immediately updates the listener with the current state
func (m *Manager) GetListener(taskUUID string) chan model.Task {
	m.lock.Lock()
	defer m.lock.Unlock()

	task, err := m.storage.Get(taskUUID)
	if err != nil {
		log.LogWarn(errors.Wrap(err, "GetListener failed to storage.Get, listener will not get current state").Error())
	}

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

	if task != nil {
		listenerChan <- *task
	}

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

	if task.Status == model.TaskStatusCompleted {
		log.LogInfo(fmt.Sprintf("task %s completed, removing all update listeners", task.UUID))
		delete(m.listeners, task.UUID)
	} else if task.Status == model.TaskStatusFailed {
		log.LogInfo(fmt.Sprintf("task %s failed, removing all update listeners", task.UUID))
		delete(m.listeners, task.UUID)
	}
}
