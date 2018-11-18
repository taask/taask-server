package service

import (
	proto1 "model/proto"

	context "golang.org/x/net/context"
)

// TaskService describes the service available to taask clients
type TaskService struct{}

// Queue handles queuing up a task to be distributed to runners
func (ts *TaaskService) Queue(context.Context, *proto1.Task) (*proto1.QueueTaskResponse, error) {

}

// CheckTask handles returning the state of a queued or running task to a client
func (ts *TaaskService) CheckTask(*proto1.CheckTaskRequest, TaaskService_CheckTaskServer) error {

}
