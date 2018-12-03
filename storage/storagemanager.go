package storage

import "github.com/taask/taask-server/model"

// Manager describes a storage implementation for tasks
type Manager interface {
	Add(*model.Task) error
	Update(*model.Task) error
	Get(string) (*model.Task, error)
	Delete(string) error
}
