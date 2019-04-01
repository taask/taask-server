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

	m.prepareTask(task)

	if err := m.addNewTask(task); err != nil {
		return "", errors.Wrap(err, "failed to addNewTask")
	}

	update, err := m.setTaskPartnerIfNeeded(task)
	if err != nil {
		return "", errors.Wrap(err, "failed to prepareTask")
	}

	// if a partner was set, and it's ours, let's set it to waiting
	if update != nil {
		if m.isOurTask(task) {
			update.Changes.Status = model.TaskStatusWaiting
		}
	} else {
		// if no partner was set, create an update to wating so that our listeners get notified
		update = task.BuildUpdate(model.TaskChanges{Status: model.TaskStatusWaiting})
	}

	updatedTask, err := m.UpdateTask(update)
	if err != nil {
		return "", errors.Wrap(err, "failed to UpdateTask")
	}

	// only schedule the task if we own it
	if m.isOurTask(updatedTask) {
		go m.scheduler.ScheduleTask(updatedTask)
	}

	return updatedTask.UUID, nil
}

// GetTask gets a task from storage
func (m *Manager) GetTask(uuid string) (*model.Task, error) {
	return m.storage.Get(uuid)
}

// UpdateTask applies a task update from a runner
func (m *Manager) UpdateTask(update *model.TaskUpdate) (*model.Task, error) {
	if update.Changes.RetrySeconds != 0 {
		return nil, errors.New("RetrySeconds is immutable")
	}

	updatedTask, err := m.updater.UpdateTask(update)
	if err != nil {
		return nil, errors.Wrap(err, "failed to updater.UpdateTask")
	}

	if m.partnerManager.HasHealthyPartner() {
		go m.partnerManager.AddTaskUpdateForSync(update)
	}

	return updatedTask, nil
}

// GetTaskUpdateListener gets an update listener for a task
func (m *Manager) GetTaskUpdateListener(uuid string) chan model.Task {
	return m.updater.GetListener(uuid)
}

func (m *Manager) isOurTask(task *model.Task) bool {
	return task.Meta.PartnerUUID == "" || task.Meta.PartnerUUID == m.partnerManager.UUID
}

func (m *Manager) addNewTask(task *model.Task) error {
	if err := m.storage.Add(*task); err != nil {
		return errors.Wrap(err, "failed to storage.Add")
	}

	if m.partnerManager.HasHealthyPartner() {
		go m.partnerManager.AddTaskForSync(task)
	}

	return nil
}

func (m *Manager) prepareTask(task *model.Task) {
	task.UUID = model.NewTaskUUID()
	task.Status = ""      // clear this in case it was set
	task.Meta.Version = 0 // set this to 0
	task.Meta.PartnerUUID = ""

	if task.Meta.TimeoutSeconds == 0 {
		task.Meta.TimeoutSeconds = 600 // 10m default; TODO: make this configurable
	}
}

func (m *Manager) setTaskPartnerIfNeeded(task *model.Task) (*model.TaskUpdate, error) {
	var changes *model.TaskChanges

	if partnerUUID, partnerPubKey := m.partnerManager.HealthyPartner(); partnerPubKey != nil {
		randomizer := rand.Intn(100)

		if randomizer < 50 {
			task.Meta.PartnerUUID = partnerUUID

			encTaskKey := task.GetEncTaskKey(m.keyService.PubKey().KID)
			if encTaskKey == nil {
				return nil, fmt.Errorf("failed to find task key encrypted with our node key (KID %s)", m.keyService.PubKey().KID)
			}

			partnerEncTaskKey, err := m.keyService.ReEncryptTaskKey(encTaskKey, partnerPubKey)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to ReEncryptTaskKey for partner %s", partnerUUID)
			}

			changes = &model.TaskChanges{PartnerUUID: partnerUUID, AddedEncTaskKeys: []*simplcrypto.Message{partnerEncTaskKey}}
		} else {
			changes = &model.TaskChanges{PartnerUUID: m.partnerManager.UUID}
		}

		fmt.Println(fmt.Sprintf("adding task with PartnerUUID %s, mine is %s", task.Meta.PartnerUUID, m.partnerManager.UUID))
	}

	if changes != nil {
		return task.BuildUpdate(*changes), nil
	}

	return nil, nil
}
