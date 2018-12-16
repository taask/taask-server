package brain

import (
	"net/http"

	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/metrics"
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

	metrics *metrics.Manager
}

// NewManager creates a new manager
func NewManager(joinCode string, storage storage.Manager) *Manager {
	metrics, err := metrics.NewManager()
	if err != nil {
		log.LogError(errors.Wrap(err, "failed to metrics.NewManager"))
		return nil
	}

	updater := update.NewManager(storage, metrics)

	scheduler := schedule.NewManager(updater)
	go scheduler.Start()

	runnerAuth := auth.NewRunnerAuthManager(joinCode)

	return &Manager{
		scheduler:  scheduler,
		storage:    storage,
		runnerAuth: runnerAuth,
		Updater:    updater,
		metrics:    metrics,
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

// ScheduleTaskRetry schedules a task to be retried
func (m *Manager) ScheduleTaskRetry(task *model.Task) {
	go m.scheduler.ScheduleTask(task)
}

// GetTask gets a task from storage
func (m *Manager) GetTask(uuid string) (*model.Task, error) {
	return m.storage.Get(uuid)
}

// UpdateTask applies a task update from a runner
func (m *Manager) UpdateTask(update *model.TaskUpdate) error {
	if update.RunnerUUID != "" {
		return errors.New("RunnerUUID is immutable")
	}

	if update.RetrySeconds != 0 {
		return errors.New("RetrySeconds is immutable")
	}

	m.Updater.UpdateTask(update)

	return nil
}

// JoinCode returns the runner join code
func (m *Manager) JoinCode() string {
	return m.runnerAuth.JoinCode
}

// MetricsHandler returns the http handler for metrics scraping
func (m *Manager) MetricsHandler() http.Handler {
	return m.metrics.Handler()
}
