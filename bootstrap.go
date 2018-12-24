package main

import (
	"io/ioutil"
	"net/http"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/storage"
)

const (
	joinCodeWritePath = "/taask/config/joincode"
)

// Bootstrap bootstraps the service
func Bootstrap() (*brain.Manager, error) {
	defer log.LogTrace("Bootstrap")()

	joinCode := auth.GenerateJoinCode()
	defer writeJoinCode(joinCode)

	runnerAuth, err := configureRunnerAuthManager(joinCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to newRunnerAuthManager")
	}

	brain := brain.NewManager(joinCode, storage.NewMemory(), runnerAuth)

	go startMetricsServer(brain)

	return brain, nil
}

func startMetricsServer(brain *brain.Manager) {
	http.Handle("/metrics", brain.MetricsHandler())

	log.LogInfo("starting metrics server on :3689")

	if err := http.ListenAndServe(":3689", nil); err != nil {
		log.LogError(errors.Wrap(err, "failed to startMetricsServer"))
	}
}

func configureRunnerAuthManager(joinCode string) (*auth.InternalAuthManager, error) {
	manager, err := auth.NewInternalAuthManager(joinCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewRunnerAuthManager")
	}

	defaultGroup := generateDefaultMemberGroup(joinCode)
	manager.AddGroup(defaultGroup)

	return manager, nil
}

func generateDefaultMemberGroup(joinCode string) *auth.MemberGroup {
	authHash := auth.GroupAuthHash(joinCode, "")

	group := &auth.MemberGroup{
		UUID:     auth.DefaultGroupUUID,
		Name:     "default",
		JoinCode: joinCode,
		AuthHash: authHash,
	}

	return group
}

func writeJoinCode(joinCode string) {
	if err := ioutil.WriteFile(joinCodeWritePath, []byte(joinCode), 0666); err != nil {
		log.LogError(errors.Wrap(err, "failed to WriteFile join code"))
	}
}
