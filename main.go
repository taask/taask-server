package main

import (
	"log"

	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/schedule"
	"github.com/taask/taask-server/service"
	"github.com/taask/taask-server/storage"
)

func main() {
	errChan := make(chan error)

	scheduler := schedule.NewManager()
	go scheduler.Start()

	brain := &brain.Manager{
		Storage:   storage.NewMemory(),
		Scheduler: scheduler,
	}

	go service.StartRunnerService(brain, errChan)
	go service.StartTaskService(brain, errChan)

	for {
		if err := <-errChan; err != nil {
			log.Fatal(err)
		}
	}
}
