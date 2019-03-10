package partner

import (
	"fmt"

	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/model"
)

// AddTaskForUpdate adds a task to be synced to partners
func (m *Manager) AddTaskForUpdate(task model.Task) {
	defer m.partner.lockUnlockUpdate()()

	m.partner.Update.AddTask(task)

	log.LogInfo(fmt.Sprintf("added task %s for partner update", task.UUID))
}
