// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package service

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	model "github.com/taask/taask-server/model"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Empty)(nil), "taask.server.service.Empty")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 284 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xc1, 0x4a, 0xf3, 0x40,
	0x14, 0x85, 0xc9, 0x0f, 0xfd, 0xc5, 0x2b, 0x75, 0x31, 0x88, 0x94, 0x08, 0x2e, 0xc4, 0x2a, 0x2e,
	0x9c, 0x94, 0xfa, 0x04, 0x2a, 0x82, 0x2e, 0x8d, 0x8a, 0x50, 0x57, 0x69, 0x7a, 0x69, 0x42, 0x4c,
	0x26, 0xce, 0xdc, 0x11, 0x7c, 0x0b, 0xdf, 0xc5, 0x17, 0x94, 0xb9, 0x93, 0x68, 0x0b, 0x53, 0xdc,
	0x84, 0x70, 0xce, 0x77, 0xcf, 0xc9, 0xcc, 0x0d, 0x0c, 0x0d, 0xea, 0xf7, 0x32, 0x47, 0xd9, 0x6a,
	0x45, 0x4a, 0xec, 0x51, 0x96, 0x99, 0x4a, 0x3a, 0x11, 0xb5, 0xec, 0xbc, 0x78, 0xbf, 0x56, 0x0b,
	0x7c, 0x4d, 0x18, 0x49, 0xc8, 0x01, 0xfc, 0x1a, 0x8f, 0x56, 0x75, 0x6d, 0x9b, 0x06, 0xb5, 0x77,
	0x8e, 0xb6, 0x60, 0x70, 0x53, 0xb7, 0xf4, 0x31, 0xfd, 0x8a, 0x60, 0xe7, 0x31, 0x33, 0xd5, 0x83,
	0x8f, 0x12, 0xb7, 0x30, 0xb8, 0xb7, 0x68, 0x51, 0x8c, 0xe4, 0x5a, 0x15, 0x27, 0x49, 0x47, 0xc6,
	0xe3, 0x90, 0xc3, 0x43, 0xce, 0x4e, 0xd1, 0xb4, 0xaa, 0x31, 0x28, 0x66, 0xb0, 0x7d, 0x5d, 0x60,
	0x5e, 0x39, 0x51, 0x1c, 0x87, 0x66, 0x7e, 0xec, 0x14, 0xdf, 0x2c, 0x1a, 0x0a, 0x27, 0xaf, 0x50,
	0x3e, 0x79, 0x12, 0x4d, 0x3f, 0xff, 0xc1, 0x30, 0xe5, 0xf3, 0xf4, 0xdf, 0xfd, 0x02, 0x70, 0x69,
	0xa9, 0xf0, 0xa2, 0x08, 0x06, 0xfd, 0xfa, 0x7d, 0xdf, 0xc9, 0x5f, 0x58, 0x77, 0x94, 0x67, 0xd8,
	0x4d, 0x71, 0x59, 0x1a, 0x42, 0xdd, 0x15, 0x9c, 0x85, 0x26, 0xd7, 0x99, 0xbe, 0x64, 0xe3, 0x45,
	0x4e, 0x22, 0x71, 0x07, 0xf0, 0xd4, 0x2e, 0x32, 0xe2, 0x9b, 0x13, 0x87, 0x9b, 0x48, 0xcf, 0xc4,
	0x07, 0x32, 0xb4, 0x7d, 0xc9, 0x8b, 0xbc, 0x3a, 0x9d, 0x8d, 0x97, 0x25, 0x15, 0x76, 0x2e, 0x73,
	0x55, 0x27, 0x0c, 0xfa, 0xe7, 0xb9, 0xc7, 0x93, 0x0e, 0x9f, 0xff, 0xe7, 0x3f, 0xe0, 0xe2, 0x3b,
	0x00, 0x00, 0xff, 0xff, 0x5d, 0xab, 0xa0, 0x58, 0x5a, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TaskServiceClient interface {
	Queue(ctx context.Context, in *model.Task, opts ...grpc.CallOption) (*model.QueueTaskResponse, error)
	CheckTask(ctx context.Context, in *model.CheckTaskRequest, opts ...grpc.CallOption) (TaskService_CheckTaskClient, error)
}

type taskServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskServiceClient(cc *grpc.ClientConn) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) Queue(ctx context.Context, in *model.Task, opts ...grpc.CallOption) (*model.QueueTaskResponse, error) {
	out := new(model.QueueTaskResponse)
	err := c.cc.Invoke(ctx, "/taask.server.service.TaskService/Queue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) CheckTask(ctx context.Context, in *model.CheckTaskRequest, opts ...grpc.CallOption) (TaskService_CheckTaskClient, error) {
	stream, err := c.cc.NewStream(ctx, &_TaskService_serviceDesc.Streams[0], "/taask.server.service.TaskService/CheckTask", opts...)
	if err != nil {
		return nil, err
	}
	x := &taskServiceCheckTaskClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TaskService_CheckTaskClient interface {
	Recv() (*model.CheckTaskResponse, error)
	grpc.ClientStream
}

type taskServiceCheckTaskClient struct {
	grpc.ClientStream
}

func (x *taskServiceCheckTaskClient) Recv() (*model.CheckTaskResponse, error) {
	m := new(model.CheckTaskResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TaskServiceServer is the server API for TaskService service.
type TaskServiceServer interface {
	Queue(context.Context, *model.Task) (*model.QueueTaskResponse, error)
	CheckTask(*model.CheckTaskRequest, TaskService_CheckTaskServer) error
}

func RegisterTaskServiceServer(s *grpc.Server, srv TaskServiceServer) {
	s.RegisterService(&_TaskService_serviceDesc, srv)
}

func _TaskService_Queue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.Task)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).Queue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/taask.server.service.TaskService/Queue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).Queue(ctx, req.(*model.Task))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_CheckTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(model.CheckTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskServiceServer).CheckTask(m, &taskServiceCheckTaskServer{stream})
}

type TaskService_CheckTaskServer interface {
	Send(*model.CheckTaskResponse) error
	grpc.ServerStream
}

type taskServiceCheckTaskServer struct {
	grpc.ServerStream
}

func (x *taskServiceCheckTaskServer) Send(m *model.CheckTaskResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _TaskService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "taask.server.service.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Queue",
			Handler:    _TaskService_Queue_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CheckTask",
			Handler:       _TaskService_CheckTask_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}

// RunnerServiceClient is the client API for RunnerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RunnerServiceClient interface {
	AuthRunner(ctx context.Context, in *model.AuthRunnerRequest, opts ...grpc.CallOption) (*model.AuthRunnerResponse, error)
	RegisterRunner(ctx context.Context, in *model.RegisterRunnerRequest, opts ...grpc.CallOption) (RunnerService_RegisterRunnerClient, error)
	UpdateTask(ctx context.Context, in *model.TaskUpdate, opts ...grpc.CallOption) (*Empty, error)
}

type runnerServiceClient struct {
	cc *grpc.ClientConn
}

func NewRunnerServiceClient(cc *grpc.ClientConn) RunnerServiceClient {
	return &runnerServiceClient{cc}
}

func (c *runnerServiceClient) AuthRunner(ctx context.Context, in *model.AuthRunnerRequest, opts ...grpc.CallOption) (*model.AuthRunnerResponse, error) {
	out := new(model.AuthRunnerResponse)
	err := c.cc.Invoke(ctx, "/taask.server.service.RunnerService/AuthRunner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerServiceClient) RegisterRunner(ctx context.Context, in *model.RegisterRunnerRequest, opts ...grpc.CallOption) (RunnerService_RegisterRunnerClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RunnerService_serviceDesc.Streams[0], "/taask.server.service.RunnerService/RegisterRunner", opts...)
	if err != nil {
		return nil, err
	}
	x := &runnerServiceRegisterRunnerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RunnerService_RegisterRunnerClient interface {
	Recv() (*model.Task, error)
	grpc.ClientStream
}

type runnerServiceRegisterRunnerClient struct {
	grpc.ClientStream
}

func (x *runnerServiceRegisterRunnerClient) Recv() (*model.Task, error) {
	m := new(model.Task)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *runnerServiceClient) UpdateTask(ctx context.Context, in *model.TaskUpdate, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/taask.server.service.RunnerService/UpdateTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RunnerServiceServer is the server API for RunnerService service.
type RunnerServiceServer interface {
	AuthRunner(context.Context, *model.AuthRunnerRequest) (*model.AuthRunnerResponse, error)
	RegisterRunner(*model.RegisterRunnerRequest, RunnerService_RegisterRunnerServer) error
	UpdateTask(context.Context, *model.TaskUpdate) (*Empty, error)
}

func RegisterRunnerServiceServer(s *grpc.Server, srv RunnerServiceServer) {
	s.RegisterService(&_RunnerService_serviceDesc, srv)
}

func _RunnerService_AuthRunner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.AuthRunnerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServiceServer).AuthRunner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/taask.server.service.RunnerService/AuthRunner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServiceServer).AuthRunner(ctx, req.(*model.AuthRunnerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RunnerService_RegisterRunner_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(model.RegisterRunnerRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RunnerServiceServer).RegisterRunner(m, &runnerServiceRegisterRunnerServer{stream})
}

type RunnerService_RegisterRunnerServer interface {
	Send(*model.Task) error
	grpc.ServerStream
}

type runnerServiceRegisterRunnerServer struct {
	grpc.ServerStream
}

func (x *runnerServiceRegisterRunnerServer) Send(m *model.Task) error {
	return x.ServerStream.SendMsg(m)
}

func _RunnerService_UpdateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.TaskUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServiceServer).UpdateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/taask.server.service.RunnerService/UpdateTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServiceServer).UpdateTask(ctx, req.(*model.TaskUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

var _RunnerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "taask.server.service.RunnerService",
	HandlerType: (*RunnerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthRunner",
			Handler:    _RunnerService_AuthRunner_Handler,
		},
		{
			MethodName: "UpdateTask",
			Handler:    _RunnerService_UpdateTask_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RegisterRunner",
			Handler:       _RunnerService_RegisterRunner_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}
