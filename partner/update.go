package partner

import (
	"sync"

	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/model"
)

// Update represents a sync between two partners
type Update struct {
	Tasks    []model.Task
	Groups   []auth.MemberGroup
	Sessions []auth.MemberAuth

	lock *sync.Mutex
}

func (u *Update) addTask(task model.Task) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.Tasks = append(u.Tasks, task)
}

func (u *Update) addGroup(group auth.MemberGroup) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.Groups = append(u.Groups, group)
}

func (u *Update) addSession(session auth.MemberAuth) {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.Sessions = append(u.Sessions, session)
}
