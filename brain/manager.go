package brain

import (
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/schedule"
	"github.com/taask/taask-server/storage"
)

// Manager is the facade for the subsystem managers (schedule, storage, update)
type Manager struct {
	Scheduler *schedule.Manager
	Storage   storage.Manager
}

// RegisterRunner registers a runner with the manager's scheduler
func (m *Manager) RegisterRunner(runner *model.Runner) {
	m.Scheduler.RegisterRunner(runner)
}

// ScheduleTask schedules and persists a task
func (m *Manager) ScheduleTask(task *model.Task) error {
	if err := m.Storage.Add(task); err != nil {
		return errors.Wrap(err, "failed to storage.Add")
	}

	m.Scheduler.ScheduleTask(task)

	return nil
}
