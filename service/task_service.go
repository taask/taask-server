package service

import (
	"net"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
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

	log.LogInfo("starting taask-server task service on :3688")
	if err := grpcServer.Serve(lis); err != nil {
		errChan <- err
	}
}

// TaskService describes the service available to taask clients
type TaskService struct {
	Manager *brain.Manager
}

// Queue handles queuing up a task to be distributed to runners
func (ts TaskService) Queue(ctx context.Context, task *model.Task) (*model.QueueTaskResponse, error) {
	uuid, err := ts.Manager.ScheduleTask(task)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ScheduleTask")
	}

	resp := &model.QueueTaskResponse{
		UUID: uuid,
	}

	return resp, nil
}

// CheckTask handles returning the state of a queued or running task to a client
func (ts TaskService) CheckTask(*model.CheckTaskRequest, TaskService_CheckTaskServer) error {
	return nil
}
