package metrics

import (
	"fmt"
	"net/http"

	"github.com/taask/taask-server/model"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsNamespace   = "taask"
	metricsUUIDLabel   = "uuid"
	metricsStatusLabel = "state"
	metricsKindLabel   = "kind"
	metricsRunnerLabel = "runner"
)

// Manager tracks metrics for export to prometheus
type Manager struct {
	metricMap map[string]*prometheus.GaugeVec
}

// NewManager returns a metrics manager
func NewManager() (*Manager, error) {
	active, err := createAndRegisterGuageVecWithName("active")
	if err != nil {
		return nil, errors.Wrap(err, "failed to createAndRegisterGuageVecWithName")
	}

	waiting, err := createAndRegisterGuageVecWithName("waiting")
	if err != nil {
		return nil, errors.Wrap(err, "failed to createAndRegisterGuageVecWithName")
	}

	retrying, err := createAndRegisterGuageVecWithName("retrying")
	if err != nil {
		return nil, errors.Wrap(err, "failed to createAndRegisterGuageVecWithName")
	}

	queued, err := createAndRegisterGuageVecWithName("queued")
	if err != nil {
		return nil, errors.Wrap(err, "failed to createAndRegisterGuageVecWithName")
	}

	running, err := createAndRegisterGuageVecWithName("running")
	if err != nil {
		return nil, errors.Wrap(err, "failed to createAndRegisterGuageVecWithName")
	}

	complete, err := createAndRegisterGuageVecWithName("complete")
	if err != nil {
		return nil, errors.Wrap(err, "failed to createAndRegisterGuageVecWithName")
	}

	failed, err := createAndRegisterGuageVecWithName("failed")
	if err != nil {
		return nil, errors.Wrap(err, "failed to createAndRegisterGuageVecWithName")
	}

	metrics := map[string]*prometheus.GaugeVec{}

	metrics["active"] = active
	metrics[model.TaskStatusWaiting] = waiting
	metrics[model.TaskStatusRetrying] = retrying
	metrics[model.TaskStatusQueued] = queued
	metrics[model.TaskStatusRunning] = running
	metrics[model.TaskStatusCompleted] = complete
	metrics[model.TaskStatusFailed] = failed

	manager := &Manager{
		metricMap: metrics,
	}

	return manager, nil
}

// Handler returns the http handler for metrics scraping
func (m *Manager) Handler() http.Handler {
	return promhttp.Handler()
}

// UpdateTask updates the current task metrics
func (m *Manager) UpdateTask(task model.Task, update *model.TaskUpdate) {
	if update.Changes.Status != "" && task.Status != update.Changes.Status {
		m.updateStatusMetrics(task, update)
	}

	// if update.Changes.RunnerUUID != "" && task.Meta.RunnerUUID != update.Changes.RunnerUUID {
	// 	log.LogInfo("updating runner metric")

	// 	m.metricMap[update.Changes.Status].With(prometheus.Labels{metricsRunnerLabel: update.Changes.RunnerUUID, metricsKindLabel: task.Kind}).Inc()
	// }

	// if update.Changes.RetrySeconds != 0 && task.RetrySeconds != update.Changes.RetrySeconds {
	// 	log.LogInfo(fmt.Sprintf("task %s set to retry in %d seconds", task.UUID, update.Changes.RetrySeconds))
	// 	task.RetrySeconds = update.Changes.RetrySeconds
	// }
}

func (m *Manager) updateStatusMetrics(task model.Task, update *model.TaskUpdate) {
	if task.Status != "" {
		old := m.metricMap[task.Status]
		old.With(prometheus.Labels{metricsRunnerLabel: task.Meta.RunnerUUID, metricsKindLabel: task.Kind}).Dec()
	} else if task.Status == "" && update.Changes.Status == model.TaskStatusWaiting {
		// we only want to increment active on the task's first update.Changes, which will always be nothing->waiting
		m.metricMap["active"].With(prometheus.Labels{metricsRunnerLabel: task.Meta.RunnerUUID, metricsKindLabel: task.Kind}).Inc()
	}

	new := m.metricMap[update.Changes.Status]
	if update.Changes.RunnerUUID != "" {
		// update.Changes the runner label if it's changed
		new.With(prometheus.Labels{metricsRunnerLabel: update.Changes.RunnerUUID, metricsKindLabel: task.Kind}).Inc()
	} else {
		new.With(prometheus.Labels{metricsRunnerLabel: task.Meta.RunnerUUID, metricsKindLabel: task.Kind}).Inc()
	}

	if update.Changes.Status == model.TaskStatusCompleted || update.Changes.Status == model.TaskStatusFailed {
		m.metricMap["active"].With(prometheus.Labels{metricsRunnerLabel: task.Meta.RunnerUUID, metricsKindLabel: task.Kind}).Dec()
	}
}

func createAndRegisterGuageVecWithName(name string) (*prometheus.GaugeVec, error) {
	vec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      fmt.Sprintf("tasks_%s", name),
		Help:      fmt.Sprintf("Tasks currently %s", name),
	}, []string{metricsRunnerLabel, metricsKindLabel})

	if err := prometheus.Register(vec); err != nil {
		return nil, errors.Wrap(err, "failed to Register")
	}

	return vec, nil
}
