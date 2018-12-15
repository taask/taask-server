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

	if task.Meta.ResultPubKey == nil {
		result.addProblem("missing ResultPubKey")
	} else {
		// TODO: validate pubkey
	}

	if task.Kind == "" {
		result.addProblem("Kind is empty")
	}

	if task.Body == nil || len(task.Body) == 0 {
		result.addProblem("Body is empty")
	}

	return result
}
