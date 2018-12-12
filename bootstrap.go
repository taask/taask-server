package main

import (
	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/storage"
)

// Bootstrap bootstraps the service
func Bootstrap() *brain.Manager {
	defer log.LogTrace("Bootstrap")()

	joinCode := auth.GenerateJoinCode()

	brain := brain.NewManager(joinCode, storage.NewMemory())

	return brain
}
