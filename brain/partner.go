package brain

import (
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
		if err := m.addTasksFromPartner(update.Tasks); err != nil {
			log.LogError(errors.Wrap(err, "PartnerUpdateFunc failed to updateTasksFromPartner"))
		}
	}
}

func (m *Manager) addTasksFromPartner(tasks []model.Task) error {
	errs := []error{}

	for _, t := range tasks {
		if err := m.storage.Add(t); err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to Add to storage for task %s", t.UUID))
		}

		// only schedule the task if we own it
		if m.isOurTask(&t) {
			update := t.BuildUpdate(model.TaskChanges{Status: model.TaskStatusWaiting})

			updatedTask, err := m.updater.UpdateTask(update)
			if err != nil {
				errs = append(errs, errors.Wrap(err, "failed to updater.UpdateTask"))
			}

			go m.scheduler.ScheduleTask(updatedTask)
		}
	}

	if err := combinedErr(errs); err != nil {
		return errors.Wrap(err, "encountered errors adding tasks")
	}

	return nil
}

func (m *Manager) updateTasksFromPartner(updates []model.TaskUpdate) error {
	errs := []error{}

	for _, u := range updates {
		_, err := m.UpdateTask(&u)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to UpdateTask"))
		}
	}

	if err := combinedErr(errs); err != nil {
		return errors.Wrap(err, "encountered errors updating tasks")
	}

	return nil
}

func combinedErr(errs []error) error {
	var combinedErr error

	if len(errs) > 0 {
		combinedErrString := ""

		for _, err := range errs {
			combinedErrString += err.Error() + "\n"
		}

		combinedErr = errors.New(combinedErrString)
	}

	return combinedErr
}
