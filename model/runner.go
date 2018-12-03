package model

// Runner describes a runner available for tasks
type Runner struct {
	UUID        string
	Type        string
	Tags        []string
	TaskChannel (chan *Task)
}
