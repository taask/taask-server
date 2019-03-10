package main

import (
	"os"

	log "github.com/cohix/simplog"
	"github.com/spf13/cobra"
	"github.com/taask/taask-server/service"
)

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "taask-server",
		Short: "Taask Server is the control plane for Taask Core, managing and scheduling tasks in coordination with the compute plane",
		Long: `A distributed task execution platform
allowing developers to run intensive and long-running compute tasks
on any infrastructure. Taask is cloud-independent, fully open source,
secure and observable by defult, and runs with zero config in most scenarios`,
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
}

func run() {
	errChan := make(chan error)

	brain, err := Bootstrap()
	if err != nil {
		log.LogError(err)
		os.Exit(1)
	}

	go service.StartRunnerService(brain, errChan)
	go service.StartTaskService(brain, errChan)
	go service.StartPartnerService(brain, errChan)

	// runtime.SetMutexProfileFraction(2)
	// runtime.SetBlockProfileRate(2)

	// log.LogInfo("Starting provile server on :6060")
	// go func() {
	// 	fmt.Println(http.ListenAndServe(":6060", nil))
	// }()

	log.LogInfo("starting taask-server")

	for {
		if err := <-errChan; err != nil {
			log.LogError(err)
			os.Exit(1)
		}
	}
}
