package service

import (
	"net"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
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

// AuthClient handles authenticating a client
func (ts TaskService) AuthClient(ctx context.Context, attempt *auth.Attempt) (*auth.AttemptResponse, error) {
	defer log.LogTrace("AuthClient")()

	encRunnerChallenge, err := ts.Manager.AttemptClientAuth(attempt)
	if err != nil {
		return nil, err
	}

	resp := &auth.AttemptResponse{
		EncChallenge: encRunnerChallenge.EncSessionChallenge,
		MasterPubKey: ts.Manager.GetMasterRunnerPubKey(),
	}

	return resp, nil
}

// Queue handles queuing up a task to be distributed to runners
func (ts TaskService) Queue(ctx context.Context, req *QueueTaskRequest) (*QueueTaskResponse, error) {
	if err := ts.Manager.CheckClientAuth(req.Session); err != nil {
		return nil, errors.Wrap(err, "failed to CheckClientAuth")
	}

	uuid, err := ts.Manager.ScheduleTask(req.Task)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ScheduleTask")
	}

	resp := &QueueTaskResponse{
		UUID: uuid,
	}

	return resp, nil
}

// CheckTask handles returning the state of a queued or running task to a client
func (ts TaskService) CheckTask(req *CheckTaskRequest, stream TaskService_CheckTaskServer) error {
	if err := ts.Manager.CheckClientAuth(req.Session); err != nil {
		return errors.Wrap(err, "failed to CheckClientAuth")
	}

	task, err := ts.Manager.GetTask(req.UUID)
	if err != nil {
		return errors.Wrap(err, "failed to GetTask")
	}

	status := ""

	listener := ts.Manager.GetTaskUpdateListener(req.UUID)

	for {
		// only send the update if the status has changed
		if task.Status != status {
			status = task.Status

			// TODO: make this more comprehensive
			update := &model.TaskUpdate{
				UUID:   task.UUID,
				Status: task.Status,
			}

			if task.EncResult != nil {
				update.EncResult = task.EncResult
			}

			resp := &CheckTaskResponse{
				Status:     task.Status,
				EncTaskKey: task.GetEncTaskKey(task.Meta.ClientKeyKID),
				Result:     update,
			}

			if err := stream.Send(resp); err != nil {
				log.LogError(errors.Wrap(err, "failed to Send"))
				return err
			}
		}

		if task.IsFinished() {
			break
		}

		updatedTask := <-listener
		task = &updatedTask
	}

	return nil
}
