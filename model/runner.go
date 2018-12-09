package model

// Runner describes a runner available for tasks
type Runner struct {
	UUID        string
	Kind        string
	Tags        []string
	TaskChannel (chan *Task)
}
