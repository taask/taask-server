package main

import (
	"fmt"
	"os"

	log "github.com/cohix/simplog"
	"github.com/taask/taask-server/service"
)

func main() {
	errChan := make(chan error)

	brain := Bootstrap()

	go service.StartRunnerService(brain, errChan)
	go service.StartTaskService(brain, errChan)

	log.LogInfo("starting taask-server")
	log.LogInfo(fmt.Sprintf("runner join code: %s", brain.JoinCode()))

	for {
		if err := <-errChan; err != nil {
			log.LogError(err)
			os.Exit(1)
		}
	}
}
