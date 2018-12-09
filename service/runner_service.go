package service

import (
	"fmt"
	"net"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
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

	RegisterRunnerServiceServer(grpcServer, &RunnerService{Manager: brain})

	log.LogInfo("starting taask-server runner service on :3687")
	if err := grpcServer.Serve(lis); err != nil {
		errChan <- err
	}
}

// RunnerService describes the service available to taask runners
type RunnerService struct {
	Manager *brain.Manager
}

// AuthRunner allows a runner to advertise itself and perform auth with the server
func (rs *RunnerService) AuthRunner(ctx context.Context, req *model.AuthRunnerRequest) (*model.AuthRunnerResponse, error) {
	defer log.LogTrace("AuthRunner")()

	return rs.Manager.AuthRunner(req)
}

// RegisterRunner allows a runner to connect (with a valid session) get a stream of tasks to execute
func (rs *RunnerService) RegisterRunner(req *model.RegisterRunnerRequest, stream RunnerService_RegisterRunnerServer) error {
	defer log.LogTrace(fmt.Sprintf("RegisterRunner kind %s", req.Kind))()

	tasksChan := make(chan *model.Task, 10)

	runner := &model.Runner{
		UUID:        model.NewRunnerUUID(),
		Kind:        req.Kind,
		Tags:        req.Tags,
		TaskChannel: tasksChan,
	}

	if err := rs.Manager.RegisterRunner(runner, req.ChallengeSignature); err != nil {
		log.LogError(errors.Wrap(err, "failed to RegisterRunner"))
		return err
	}

	defer rs.Manager.UnregisterRunner(runner)

	for {
		task := <-tasksChan

		if err := stream.Send(task); err != nil {
			log.LogError(errors.Wrap(err, "failed to stream.Send"))
			break
		}
	}

	return nil
}

// UpdateTask handles update task calls
func (rs *RunnerService) UpdateTask(ctx context.Context, req *model.TaskUpdate) (*Empty, error) {
	defer log.LogTrace(fmt.Sprintf("UpdateTask task %s", req.UUID))

	if err := rs.Manager.UpdateTask(req); err != nil {
		log.LogError(errors.Wrap(err, "failed to UpdateTask"))
		return nil, err
	}

	return &Empty{}, nil
}
