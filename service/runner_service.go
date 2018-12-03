package service

import (
	"log"
	"net"

	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/model"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

const runnerServicePort = ":3687"

// StartRunnerService starts the runner service
func StartRunnerService(brain *brain.Manager, errChan chan error) {
	lis, err := net.Listen("tcp", runnerServicePort)
	if err != nil {
		errChan <- err
		return
	}

	grpcServer := grpc.NewServer()

	RegisterRunnerServiceServer(grpcServer, RunnerService{Manager: brain})

	log.Println("starting taask-server runner service on :3687")
	if err := grpcServer.Serve(lis); err != nil {
		errChan <- err
	}
}

// RunnerService describes the service available to taask runners
type RunnerService struct {
	Manager *brain.Manager
}

// RegisterRunner allows a runner to advertise itself and perform auth with the server
func (rs RunnerService) RegisterRunner(ctx context.Context, req *model.RegisterRunnerRequest) (*model.RegisterRunnerResponse, error) {
	return &model.RegisterRunnerResponse{}, nil
}

// StreamTasks allows a runner to connect (with a valid session) get a stream of tasks to execute
func (rs RunnerService) StreamTasks(req *model.StreamTasksRequest, stream RunnerService_StreamTasksServer) error {
	return nil
}
