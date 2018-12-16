package service

import (
	"fmt"
	"net"
	"time"

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

	tasksChan := make(chan *model.Task, 128)

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
	go startRunnerHeartbeat(tasksChan)

	log.LogInfo(fmt.Sprintf("runner %s ready to receive tasks", runner.UUID))

	for {
		task := <-tasksChan

		if task.UUID != "" {
			log.LogInfo(fmt.Sprintf("runner %s handling task %s", runner.UUID, task.UUID))
		} else {
			log.LogInfo(fmt.Sprintf("sending runner %s heartbeat", runner.UUID))
		}

		if err := stream.Send(task); err != nil {
			log.LogError(errors.Wrap(err, "failed to stream.Send"))

			if task.UUID != "" {
				rs.Manager.ScheduleTaskRetry(task)
			}

			break
		}
	}

	return nil
}

// UpdateTask handles update task calls
func (rs *RunnerService) UpdateTask(ctx context.Context, req *model.TaskUpdate) (*Empty, error) {
	defer log.LogTrace(fmt.Sprintf("UpdateTask task %s", req.UUID))()

	rs.Manager.Updater.UpdateTask(req)

	return &Empty{}, nil
}

// TODO: find a way to terminate this
func startRunnerHeartbeat(taskChan chan *model.Task) {
	for {
		<-time.After(time.Second * time.Duration(10))

		taskChan <- &model.Task{}
	}
}
