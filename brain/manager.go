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
	"github.com/taask/taask-server/update"
)

// Manager is the facade for the subsystem managers (schedule, storage, update, auth)
type Manager struct {
	scheduler  *schedule.Manager
	storage    storage.Manager
	runnerAuth *auth.RunnerAuthManager
	Updater    *update.Manager
}

// NewManager creates a new manager
func NewManager(joinCode string, storage storage.Manager) *Manager {
	updater := update.NewManager(storage)

	scheduler := schedule.NewManager(updater)
	go scheduler.Start()

	runnerAuth := auth.NewRunnerAuthManager(joinCode)

	return &Manager{
		scheduler:  scheduler,
		storage:    storage,
		runnerAuth: runnerAuth,
		Updater:    updater,
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

	if err := m.storage.Add(*task); err != nil {
		return "", errors.Wrap(err, "failed to storage.Add")
	}

	go func() {
		m.scheduler.ScheduleTask(task)
	}()

	return task.UUID, nil
}

// ScheduleTaskRetry schedules a task to be retried
func (m *Manager) ScheduleTaskRetry(task *model.Task) {
	m.scheduler.StartRetryWorker(task)
}

// GetTask gets a task from storage
func (m *Manager) GetTask(uuid string) (*model.Task, error) {
	return m.storage.Get(uuid)
}

// JoinCode returns the runner join code
func (m *Manager) JoinCode() string {
	return m.runnerAuth.JoinCode
}
