package validator

import "github.com/taask/taask-server/model"

// ValidateTask validates a task
func ValidateTask(task *model.Task) *Result {
	result := resultOfType(typeTask)
	if task == nil {
		result.addProblem("task is nil")
		return result
	}

	if task.Meta == nil {
		result.addProblem("missing Meta")
		return result
	}

	if task.Kind == "" {
		result.addProblem("Kind is empty")
	}

	if task.EncBody == nil || len(task.EncBody.Data) == 0 {
		result.addProblem("EncBody is empty")
	}

	return result
}
