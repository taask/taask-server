package main

import (
	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/schedule"
	"github.com/taask/taask-server/storage"
)

// Bootstrap bootstraps the service
func Bootstrap() *brain.Manager {
	defer log.LogTrace("Bootstrap")()

	scheduler := schedule.NewManager()
	go scheduler.Start()

	joinCode := auth.GenerateJoinCode()
	runnerAuthManager := auth.NewRunnerAuthManager(joinCode)

	brain := brain.NewManager(scheduler, storage.NewMemory(), runnerAuthManager)

	return brain
}
