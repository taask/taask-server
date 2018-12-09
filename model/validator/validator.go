package validator

import (
	"errors"
	"fmt"
)

const (
	typeTask = "com.taask.task"
)

// Result is the result of a validation
type Result struct {
	Type   string
	Errors []error
}

func resultOfType(resType string) *Result {
	return &Result{
		Type:   resType,
		Errors: make([]error, 0),
	}
}

func (vr *Result) addProblem(prob string) {
	vr.Errors = append(vr.Errors, errors.New(prob))
}

// String returns a string representation of a result
func (vr *Result) String() string {
	errString := fmt.Sprintf("validator found %d problems:")

	for i, e := range vr.Errors {
		errString += fmt.Sprintf("\n\t%d: %s", i+1, e.Error())
	}

	return errString
}

// Ok returns the result's status
func (vr *Result) Ok() bool {
	return len(vr.Errors) == 0
}
