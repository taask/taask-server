package brain

import (
	"fmt"

	"github.com/cohix/simplcrypto"
	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/partner"
	"github.com/taask/taask-server/update"
)

// GetMasterPartnerPubKey returns the master runner pubkey
func (m *Manager) GetMasterPartnerPubKey() *simplcrypto.SerializablePubKey {
	return m.keyService.PubKey()
}

// RunPartnerManagerWithServer runs the partner manager with a stream server
func (m *Manager) RunPartnerManagerWithServer(server partner.PartnerService_StreamUpdatesServer) error {
	return m.partnerManager.RunWithServer(server)
}

// PartnerUpdateFunc gets the update func for the partner manager
func (m *Manager) PartnerUpdateFunc() func(update.PartnerUpdate) {
	return func(update update.PartnerUpdate) {
		if err := m.updateTasksFromPartner(update.Tasks); err != nil {
			log.LogError(errors.Wrap(err, "PartnerUpdateFunc failed to updateTasksFromPartner"))
		}
	}
}

func (m *Manager) updateTasksFromPartner(tasks []model.Task) error {
	errs := []error{}

	for i, t := range tasks {
		canSchedule := true

		if err := m.storage.Add(t); err != nil {
			log.LogInfo(fmt.Sprintf("task %s from partner already exists, will update", t.UUID))

			if err := m.storage.Update(t); err != nil {
				errs = append(errs, errors.Wrap(err, "failed to update task"))
				canSchedule = false
			}
		}

		// if we own it, schedule it
		if canSchedule && t.Meta.PartnerUUID == m.partnerManager.UUID {
			if t.IsNotStarted() {
				m.scheduler.ScheduleTask(&tasks[i])
			} else {
				log.LogWarn("received task from partner owned by self, but is started")
			}
		}
	}

	var combinedErr error
	if len(errs) > 0 {
		combinedErrString := ""

		for _, err := range errs {
			combinedErrString += err.Error() + "\n"
		}

		combinedErr = errors.New(combinedErrString)
	}

	if combinedErr != nil {
		return errors.Wrap(combinedErr, "encountered errors updating tasks")
	}

	return nil
}
