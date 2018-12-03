package service

import (
	"log"
	"net"

	"github.com/taask/taask-server/brain"
	model "github.com/taask/taask-server/model"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

const taskServicePort = ":3688"

// StartTaskService starts the runner service
func StartTaskService(brain *brain.Manager, errChan chan error) {
	lis, err := net.Listen("tcp", taskServicePort)
	if err != nil {
		errChan <- err
		return
	}

	grpcServer := grpc.NewServer()

	RegisterTaskServiceServer(grpcServer, TaskService{Manager: brain})

	log.Println("starting taask-server task service on :3688")
	if err := grpcServer.Serve(lis); err != nil {
		errChan <- err
	}
}

// TaskService describes the service available to taask clients
type TaskService struct {
	Manager *brain.Manager
}

// Queue handles queuing up a task to be distributed to runners
func (ts TaskService) Queue(context.Context, *model.Task) (*model.QueueTaskResponse, error) {
	return &model.QueueTaskResponse{}, nil
}

// CheckTask handles returning the state of a queued or running task to a client
func (ts TaskService) CheckTask(*model.CheckTaskRequest, TaskService_CheckTaskServer) error {
	return nil
}
