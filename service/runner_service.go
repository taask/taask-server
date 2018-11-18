package service

import (
	proto1 "model/proto"

	context "golang.org/x/net/context"
)

// RunnerService describes the service available to taask runners
type RunnerService struct{}

// RegisterRunner allows a runner to advertise itself and perform auth with the server
func (rs *RunnerService) RegisterRunner(context.Context, *proto1.RegisterRunnerRequest) (*proto1.RegisterRunnerResponse, error) {

}

// StreamTasks allows a runner to connect (with a valid session) get a stream of tasks to execute
func (rs *RunnerService) StreamTasks(*proto1.StreamTasksRequest, TaaskService_StreamTasksServer) error {

}
