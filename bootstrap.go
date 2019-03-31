package main

import (
	"fmt"
	"net/http"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/config"
	"github.com/taask/taask-server/keyservice"
	"github.com/taask/taask-server/partner"
	"github.com/taask/taask-server/storage"
)

// Bootstrap bootstraps the service
func Bootstrap() (*brain.Manager, error) {
	serverConfig, err := config.ServerConfigFromDefaultDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ServerConfigFromDefaultDir")
	}

	defer log.LogTrace("Bootstrap")()

	// TODO: load node key from config?
	keyservice, err := keyservice.NewManager()
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewManager for keyservice")
	}

	runnerAuth, err := configureRunnerAuthManager(&serverConfig.RunnerAuth.MemberGroup, keyservice)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configureRunnerAuthManager")
	}

	clientAuth, err := configureClientAuthManager(&serverConfig.ClientAuth.MemberGroup, keyservice)
	if err != nil {
		return nil, errors.Wrap(err, "failed to configureClientAuthManager")
	}

	var partnerManager *partner.Manager

	if serverConfig.PartnerAuth != nil {
		partnerManager, err = configurePartnerManager(serverConfig.PartnerAuth, keyservice)
		if err != nil {
			return nil, errors.Wrap(err, "failed to configurePartnerManager")
		} else if partnerManager == nil {
			return nil, errors.Wrap(err, "failed to configurePartnerManager, it was nil")
		}

		go partnerManager.StartOutgoingManager()

	} else {
		log.LogInfo("partner manager not configured")
	}

	brain := brain.NewManager(keyservice, storage.NewMemory(), runnerAuth, clientAuth, partnerManager)

	if partnerManager != nil {
		partnerManager.SetApplyUpdateFunc(brain.PartnerUpdateFunc())
	}

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

func configureRunnerAuthManager(defaultGroup *auth.MemberGroup, keyservice keyservice.KeyService) (*auth.InternalAuthManager, error) {
	manager, err := auth.NewInternalAuthManager(keyservice)
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewInternalAuthManager")
	}

	if defaultGroup.Name != "default" {
		return nil, fmt.Errorf("runner auth config with group name %s not allowed", defaultGroup.Name)
	}

	if defaultGroup.UUID != auth.DefaultGroupUUID {
		return nil, fmt.Errorf("runner auth config with group uuid %s not allowed", defaultGroup.UUID)
	}

	if err := manager.AddGroup(defaultGroup); err != nil {
		return nil, errors.Wrap(err, "failed to AddGroup")
	}

	return manager, nil
}

func configureClientAuthManager(adminGroup *auth.MemberGroup, keyservice keyservice.KeyService) (*auth.InternalAuthManager, error) {
	manager, err := auth.NewInternalAuthManager(keyservice)
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewInternalAuthManager")
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

func configurePartnerManager(config *config.ClientAuthConfig, keyservice keyservice.KeyService) (*partner.Manager, error) {
	partnerManager, err := partner.NewManager(config, keyservice)
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewManager")
	}

	return partnerManager, nil
}
