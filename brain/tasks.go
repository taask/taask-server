package brain

import (
	"fmt"
	"math/rand"

	"github.com/cohix/simplcrypto"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/model/validator"
)

// ScheduleTask schedules and persists a task
func (m *Manager) ScheduleTask(task *model.Task) (string, error) {
	if result := validator.ValidateTask(task); !result.Ok() {
		return "", errors.Wrap(errors.New(result.String()), "failed to ValidateTask")
	}

	task.UUID = model.NewTaskUUID()
	task.Status = ""      // clear this in case it was set
	task.Meta.Version = 0 // set this to 0
	task.Meta.PartnerUUID = ""

	if task.Meta.TimeoutSeconds == 0 {
		task.Meta.TimeoutSeconds = 600 // 10m default; TODO: make this configurable
	}

	encTaskKeys := []*simplcrypto.Message{}

	if partnerUUID, partnerPubKey := m.partnerManager.HealthyPartner(); partnerPubKey != nil {
		randomizer := rand.Intn(100)

		if randomizer < 50 {
			task.Meta.PartnerUUID = partnerUUID

			encTaskKey := task.GetEncTaskKey(m.keyService.PubKey().KID)
			if encTaskKey == nil {
				return "", fmt.Errorf("failed to find task key encrypted with our node key (KID %s)", m.keyService.PubKey().KID)
			}

			partnerEncTaskKey, err := m.keyService.ReEncryptTaskKey(encTaskKey, partnerPubKey)
			if err != nil {
				return "", errors.Wrapf(err, "failed to ReEncryptTaskKey for partner %s", partnerUUID)
			}

			encTaskKeys = append(encTaskKeys, partnerEncTaskKey)
		} else {
			task.Meta.PartnerUUID = m.partnerManager.UUID
		}

		fmt.Println(fmt.Sprintf("adding task with PartnerUUID %s, mine is %s", task.Meta.PartnerUUID, m.partnerManager.UUID))
	}

	if err := m.storage.Add(*task); err != nil {
		return "", errors.Wrap(err, "failed to storage.Add")
	}

	// we do a manual update to waiting to ensure the metrics catch the new task
	update, err := task.Update(model.TaskUpdate{Status: model.TaskStatusWaiting, AddedEncTaskKeys: encTaskKeys, PartnerUUID: task.Meta.PartnerUUID})
	if err != nil {
		return "", errors.Wrap(err, "failed to task.Update")
	}

	m.UpdateTask(update)

	// only schedule the task if we own it
	if task.Meta.PartnerUUID == "" || task.Meta.PartnerUUID == m.partnerManager.UUID {
		go m.scheduler.ScheduleTask(task)
	}

	return task.UUID, nil
}

// GetTask gets a task from storage
func (m *Manager) GetTask(uuid string) (*model.Task, error) {
	return m.storage.Get(uuid)
}

// UpdateTask applies a task update from a runner
func (m *Manager) UpdateTask(update model.TaskUpdate) error {
	if update.RetrySeconds != 0 {
		return errors.New("RetrySeconds is immutable")
	}

	task := m.updater.UpdateTask(update)

	go m.partnerManager.AddTaskForUpdate(task)

	return nil
}

// GetTaskUpdateListener gets an update listener for a task
func (m *Manager) GetTaskUpdateListener(uuid string) chan model.Task {
	return m.updater.GetListener(uuid)
}
