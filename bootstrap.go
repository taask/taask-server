package main

import (
	"net/http"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/storage"
)

// Bootstrap bootstraps the service
func Bootstrap() *brain.Manager {
	defer log.LogTrace("Bootstrap")()

	joinCode := auth.GenerateJoinCode()

	brain := brain.NewManager(joinCode, storage.NewMemory())

	go startMetricsServer(brain)

	return brain
}

func startMetricsServer(brain *brain.Manager) {
	http.Handle("/metrics", brain.MetricsHandler())

	log.LogInfo("starting metrics server on :3689")

	if err := http.ListenAndServe(":3689", nil); err != nil {
		log.LogError(errors.Wrap(err, "failed to startMetricsServer"))
	}
}
