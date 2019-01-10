package brain

import (
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/model/validator"
)

// ScheduleTask schedules and persists a task
func (m *Manager) ScheduleTask(task *model.Task) (string, error) {
	if result := validator.ValidateTask(task); !result.Ok() {
		return "", errors.Wrap(errors.New(result.String()), "failed to ValidateTask")
	}

	task.UUID = model.NewTaskUUID()
	task.Status = ""      // clear this in case it was set
	task.Meta.Version = 0 // set this to 0
	if task.Meta.TimeoutSeconds == 0 {
		task.Meta.TimeoutSeconds = 600 // 10m default
	}

	if err := m.storage.Add(*task); err != nil {
		return "", errors.Wrap(err, "failed to storage.Add")
	}

	// we do a manual update to waiting to ensure the metrics catch the new task
	update, err := task.Update(model.TaskUpdate{Status: model.TaskStatusWaiting})
	if err != nil {
		return "", errors.Wrap(err, "failed to task.Update")
	}

	m.Updater.UpdateTask(update)

	go m.scheduler.ScheduleTask(task)

	return task.UUID, nil
}

// GetTask gets a task from storage
func (m *Manager) GetTask(uuid string) (*model.Task, error) {
	return m.storage.Get(uuid)
}

// UpdateTask applies a task update from a runner
func (m *Manager) UpdateTask(update model.TaskUpdate) error {
	if update.RunnerUUID != "" {
		return errors.New("RunnerUUID is immutable")
	}

	if update.RetrySeconds != 0 {
		return errors.New("RetrySeconds is immutable")
	}

	m.Updater.UpdateTask(update)

	return nil
}
