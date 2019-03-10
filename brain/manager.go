package brain

import (
	"net/http"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/metrics"
	"github.com/taask/taask-server/partner"
	"github.com/taask/taask-server/schedule"
	"github.com/taask/taask-server/storage"
	"github.com/taask/taask-server/update"
)

// Manager is the facade for the subsystem managers (schedule, storage, update, auth)
type Manager struct {
	// Updater manages updates to tasks and coordinates listeners, storage, and metrics
	updater *update.Manager

	// Scheduler consumes tasks and schedules them onto the compute plane
	scheduler *schedule.Manager

	// storage manages persistence of tasks, etc
	storage storage.Manager

	// runnerAuth manages the authentication of runners
	runnerAuth auth.Manager

	// clientAuth manages the authentication of clients
	clientAuth auth.Manager

	// partnerManager manages sync with partners
	partnerManager *partner.Manager

	// metrics manages observability
	metrics *metrics.Manager
}

// NewManager creates a new manager
func NewManager(storage storage.Manager, runnerAuth, clientAuth auth.Manager, partnerManager *partner.Manager) *Manager {
	metrics, err := metrics.NewManager()
	if err != nil {
		log.LogError(errors.Wrap(err, "failed to metrics.NewManager"))
		return nil
	}

	updater := update.NewManager(storage, metrics)

	if partnerManager == nil {
		partnerManager = &partner.Manager{UUID: ""}
	}

	scheduler := schedule.NewManager(updater)
	go scheduler.Start()

	return &Manager{
		updater: updater,

		scheduler: scheduler,

		storage: storage,

		runnerAuth: runnerAuth,

		clientAuth: clientAuth,

		partnerManager: partnerManager,

		metrics: metrics,
	}
}

// MetricsHandler returns the http handler for metrics scraping
func (m *Manager) MetricsHandler() http.Handler {
	return m.metrics.Handler()
}
