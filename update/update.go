package update

import (
	"sync"

	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/model"
)

// PartnerUpdate represents a sync between two partners
type PartnerUpdate struct {
	Tasks    []model.Task
	Groups   []auth.MemberGroup
	Sessions []auth.MemberAuth

	// we need to lock the update itself to protect against multiple goroutines
	lock *sync.Mutex
}

// NewPartnerUpdate returns an empty partner update
func NewPartnerUpdate() *PartnerUpdate {
	update := &PartnerUpdate{
		Tasks:    []model.Task{},
		Groups:   []auth.MemberGroup{},
		Sessions: []auth.MemberAuth{},
		lock:     &sync.Mutex{},
	}

	return update
}

// AddTask adds a task to be updated
func (u *PartnerUpdate) AddTask(task model.Task) {
	defer u.lockUnlock()()

	for i, t := range u.Tasks {
		if t.UUID == task.UUID {
			u.Tasks[i] = task
			return
		}
	}

	u.Tasks = append(u.Tasks, task)
}

// AddGroup adds a member group to be updated
func (u *PartnerUpdate) AddGroup(group auth.MemberGroup) {
	defer u.lockUnlock()()

	u.Groups = append(u.Groups, group)
}

// AddSession adds a member session to be synced
func (u *PartnerUpdate) AddSession(session auth.MemberAuth) {
	defer u.lockUnlock()()

	u.Sessions = append(u.Sessions, session)
}

func (u *PartnerUpdate) lockUnlock() func() {
	u.lock.Lock()

	return func() {
		u.lock.Unlock()
	}
}
