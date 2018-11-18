package main

import (
	"log"

	"github.com/taask/taask-server/service"
)

func main() {
	errChan := make(chan error)

	go service.StartRunnerService(errChan)
	go service.StartTaskService(errChan)

	for {
		if err := <-errChan; err != nil {
			log.Fatal(err)
		}
	}
}
