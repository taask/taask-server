package brain

import (
	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/model/validator"
	"github.com/taask/taask-server/schedule"
	"github.com/taask/taask-server/storage"
)

// Manager is the facade for the subsystem managers (schedule, storage, update, auth)
type Manager struct {
	scheduler  *schedule.Manager
	storage    storage.Manager
	runnerAuth *auth.RunnerAuthManager
}

// NewManager creates a new manager
func NewManager(scheduler *schedule.Manager, storage storage.Manager, runnerAuth *auth.RunnerAuthManager) *Manager {
	return &Manager{
		scheduler:  scheduler,
		storage:    storage,
		runnerAuth: runnerAuth,
	}
}

// AuthRunner allows a runner to auth
func (m *Manager) AuthRunner(authReq *model.AuthRunnerRequest) (*model.AuthRunnerResponse, error) {
	return m.runnerAuth.AttemptAuth(authReq.PubKey, authReq.JoinCodeSignature)
}

// RegisterRunner registers a runner with the manager's scheduler
func (m *Manager) RegisterRunner(runner *model.Runner, challengeSignature *simplcrypto.Signature) error {
	if err := m.runnerAuth.CheckRunnerChallenge(challengeSignature); err != nil {
		return errors.Wrap(err, "failed to CheckRunnerChallenge")
	}

	m.scheduler.RegisterRunner(runner)

	return nil
}

// UnregisterRunner unregisters a runner
func (m *Manager) UnregisterRunner(runner *model.Runner) {
	// TODO: reassign currently running tasks
	if err := m.scheduler.UnregisterRunner(runner.Kind, runner.UUID); err != nil {
		log.LogError(errors.Wrap(err, "failed to UnregisterRunner"))
	}
}

// ScheduleTask schedules and persists a task
func (m *Manager) ScheduleTask(task *model.Task) (string, error) {
	if result := validator.ValidateTask(task); !result.Ok() {
		return "", errors.Wrap(errors.New(result.String()), "failed to ValidateTask")
	}

	task.UUID = model.NewTaskUUID()
	task.Status = model.TaskStatusWaiting

	if err := m.storage.Add(task); err != nil {
		return "", errors.Wrap(err, "failed to storage.Add")
	}

	go func(storage storage.Manager) {
		m.scheduler.ScheduleTask(task)

		task.Status = model.TaskStatusQueued
		if err := storage.Update(task); err != nil {
			log.LogError(errors.Wrap(err, "failed to storage.Update"))
		}
	}(m.storage)

	return task.UUID, nil
}

// UpdateTask updates a task
func (m *Manager) UpdateTask(update *model.TaskUpdate) error {
	task, err := m.storage.Get(update.UUID)
	if err != nil {
		return errors.Wrap(err, "failed to storage.Get")
	}

	task.Status = update.Status
	if update.EncResult != nil {
		task.EncResult = update.EncResult
		task.EncResultSymKey = update.EncResultSymKey
	}

	return m.storage.Update(task)
}

// JoinCode returns the runner join code
func (m *Manager) JoinCode() string {
	return m.runnerAuth.JoinCode
}
