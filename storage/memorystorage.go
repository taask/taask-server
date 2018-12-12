package storage

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
)

// Memory is the default, in-memory storage option for persisting tasks
type Memory struct {
	tasks map[string]*model.Task
	lock  *sync.Mutex
}

// NewMemory creates a new in-memory store
func NewMemory() *Memory {
	return &Memory{
		tasks: make(map[string]*model.Task),
		lock:  &sync.Mutex{},
	}
}

// Add adds a task to storage
func (sm *Memory) Add(task model.Task) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	if _, ok := sm.tasks[task.UUID]; ok {
		return errors.Wrap(fmt.Errorf("task with uuid %s already exists", task.UUID), "failed to add task to map")
	}

	sm.tasks[task.UUID] = &task

	return nil
}

// Update updates a task in storage (task is matched on UUID and replaced)
func (sm *Memory) Update(task model.Task) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	if _, ok := sm.tasks[task.UUID]; !ok {
		return errors.Wrap(fmt.Errorf("task with uuid %s does not exist", task.UUID), "failed to update task")
	}

	sm.tasks[task.UUID] = &task

	return nil
}

// Get retreives a task from storage by UUID
func (sm *Memory) Get(uuid string) (*model.Task, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	task, ok := sm.tasks[uuid]
	if !ok {
		return nil, errors.Wrap(fmt.Errorf("task with uuid %s does not exist", uuid), "failed to get task")
	}

	taskCopy := *task
	return &taskCopy, nil
}

// Delete deletes a task by UUID
func (sm *Memory) Delete(uuid string) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	if _, ok := sm.tasks[uuid]; !ok {
		return errors.Wrap(fmt.Errorf("task with uuid %s does not exist", uuid), "failed to delete task")
	}

	delete(sm.tasks, uuid)

	return nil
}
