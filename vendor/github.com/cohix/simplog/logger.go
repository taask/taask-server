package simplog

import (
	"fmt"
	"time"

	flag "github.com/cohix/simplflag"
)

var includeDebug bool

func init() {
	_, includeDebug = flag.CheckFlag("debug")
}

type logMessage struct {
	Time     time.Time `json:"time"`
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	TraceID  string    `json:"trace_id"`
}

// LogError logs an error
func LogError(err error) {
	fmt.Printf("(E) %s\n", err.Error())
}

// LogWarn logs a warning
func LogWarn(msg string) {
	fmt.Printf("(W) %s\n", msg)
}

// LogInfo logs an information message
func LogInfo(msg string) {
	fmt.Printf("(I) %s\n", msg)
}

// LogDebug logs an information message
func LogDebug(msg string) {
	if includeDebug {
		fmt.Printf("(D) %s\n", msg)
	}
}

// LogTrace logs a function call and returns a function to be deferred marking the end of the function
func LogTrace(name string) func() {
	LogInfo(fmt.Sprintf("[trace] %s began", name))
	return func() {
		LogInfo(fmt.Sprintf("[trace] %s completed", name))
	}
}
