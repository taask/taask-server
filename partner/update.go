package partner

import (
	"fmt"

	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/model"
)

// AddTaskForSync adds a task to be synced to partners
func (m *Manager) AddTaskForSync(task *model.Task) {
	if task == nil {
		// TODO: handle health checks better
		return
	}

	// lock the update at the manager level, AddTask locks the update internally
	defer m.partner.lockUnlockUpdate()()

	m.partner.Update.AddTask(*task)

	log.LogInfo(fmt.Sprintf("added task %s for partner sync", task.UUID))
}

// AddTaskUpdateForSync adds a task to be synced to partners
func (m *Manager) AddTaskUpdateForSync(update *model.TaskUpdate) {
	// lock the update at the manager level, AddTaskUpdate locks the update internally
	defer m.partner.lockUnlockUpdate()()

	m.partner.Update.AddTaskUpdate(update)

	log.LogInfo(fmt.Sprintf("added task %s update for partner sync", update.UUID))
}
