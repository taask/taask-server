package brain

import (
	"net/http"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/metrics"
	"github.com/taask/taask-server/schedule"
	"github.com/taask/taask-server/storage"
	"github.com/taask/taask-server/update"
)

// Manager is the facade for the subsystem managers (schedule, storage, update, auth)
type Manager struct {
	// Updater manages updates to tasks and coordinates listeners, storage, and metrics
	Updater *update.Manager

	// Scheduler consumes tasks and schedules them onto the compute plane
	scheduler *schedule.Manager
	storage   storage.Manager

	// runnerAuth manages the authentication of runners
	runnerAuth     auth.Manager
	RunnerJoinCode string

	// clientAuth manages the authentication of clients
	clientAuth auth.Manager

	// metrics manages observability
	metrics *metrics.Manager
}

// NewManager creates a new manager
func NewManager(storage storage.Manager, runnerAuth, clientAuth auth.Manager) *Manager {
	metrics, err := metrics.NewManager()
	if err != nil {
		log.LogError(errors.Wrap(err, "failed to metrics.NewManager"))
		return nil
	}

	updater := update.NewManager(storage, metrics)

	scheduler := schedule.NewManager(updater)
	go scheduler.Start()

	return &Manager{
		Updater: updater,

		scheduler: scheduler,
		storage:   storage,

		runnerAuth: runnerAuth,

		clientAuth: clientAuth,

		metrics: metrics,
	}
}

// MetricsHandler returns the http handler for metrics scraping
func (m *Manager) MetricsHandler() http.Handler {
	return m.metrics.Handler()
}
