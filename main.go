package main

import (
	"fmt"
	"os"

	_ "net/http/pprof"

	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/service"
)

func main() {
	errChan := make(chan error)

	brain := Bootstrap()

	go service.StartRunnerService(brain, errChan)
	go service.StartTaskService(brain, errChan)

	// runtime.SetMutexProfileFraction(2)
	// runtime.SetBlockProfileRate(2)

	// log.LogInfo("Starting provile server on :6060")
	// go func() {
	// 	fmt.Println(http.ListenAndServe(":6060", nil))
	// }()

	log.LogInfo("starting taask-server")
	log.LogInfo(fmt.Sprintf("runner join code: %s", brain.JoinCode()))

	for {
		if err := <-errChan; err != nil {
			log.LogError(err)
			os.Exit(1)
		}
	}
}
