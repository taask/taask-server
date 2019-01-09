package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/config"
	"github.com/taask/taask-server/storage"
)

var joinCodeWritePath = filepath.Join(config.DefaultConfigDir(), "joincode")

// Bootstrap bootstraps the service
func Bootstrap() (*brain.Manager, error) {
	serverConfig, err := config.ServerConfigFromDefaultDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ServerConfigFromDefaultDir")
	}

	defer log.LogTrace("Bootstrap")()

	joinCode := auth.GenerateJoinCode()
	defer writeJoinCode(joinCode)

	runnerAuth, err := configureRunnerAuthManager(joinCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to newRunnerAuthManager")
	}

	clientAuth, err := configureClientAuthManager(&serverConfig.ClientAuth.AdminGroup)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configureClientAuthManager")
	}

	brain := brain.NewManager(joinCode, storage.NewMemory(), runnerAuth, clientAuth)

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
	manager, err := auth.NewInternalAuthManager()
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

func configureClientAuthManager(adminGroup *auth.MemberGroup) (*auth.InternalAuthManager, error) {
	manager, err := auth.NewInternalAuthManager()
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewRunnerAuthManager")
	}

	if adminGroup.Name != "admin" {
		return nil, fmt.Errorf("client auth config with group name %s not allowed", adminGroup.Name)
	}

	if adminGroup.UUID != auth.AdminGroupUUID {
		return nil, fmt.Errorf("client auth config with group uuid %s not allowed", adminGroup.UUID)
	}

	if err := manager.AddGroup(adminGroup); err != nil {
		return nil, errors.Wrap(err, "failed to AddGroup")
	}

	return manager, nil
}

func writeJoinCode(joinCode string) {
	if err := ioutil.WriteFile(joinCodeWritePath, []byte(joinCode), 0666); err != nil {
		log.LogError(errors.Wrap(err, "failed to WriteFile join code"))
	}
}
