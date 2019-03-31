package brain

import (
	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
)

// RegisterRunner registers a runner with the manager's scheduler
func (m *Manager) RegisterRunner(runner *model.Runner) error {
	m.scheduler.RegisterRunner(runner)

	return nil
}

// UnregisterRunner unregisters a runner
func (m *Manager) UnregisterRunner(runner *model.Runner) {
	if err := m.scheduler.UnregisterRunner(runner.Kind, runner.UUID); err != nil {
		log.LogError(errors.Wrap(err, "failed to UnregisterRunner"))
	}

	if err := m.runnerAuth.DeleteMemberAuth(runner.UUID); err != nil {
		log.LogError(errors.Wrap(err, "failed to DeleteRunnerKey"))
	}
}

// EncryptTaskKeyForRunner encrypts a task key for a runner
func (m *Manager) EncryptTaskKeyForRunner(runnerUUID string, encTaskKey *simplcrypto.Message) (*simplcrypto.Message, error) {
	runnerPubKey, err := m.runnerAuth.MemberPubkey(runnerUUID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to MemberPubkey")
	}

	encKey, err := m.keyService.ReEncryptTaskKey(encTaskKey, runnerPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReEncryptTaskKey")
	}

	return encKey, nil
}

// GetMasterRunnerPubKey returns the master runner pubkey
func (m *Manager) GetMasterRunnerPubKey() *simplcrypto.SerializablePubKey {
	return m.keyService.PubKey()
}
