// Code generated by protoc-gen-go. DO NOT EDIT.
// source: taskservice.proto

package service

import (
	fmt "fmt"
	simplcrypto "github.com/cohix/simplcrypto"
	proto "github.com/golang/protobuf/proto"
	auth "github.com/taask/taask-server/auth"
	model "github.com/taask/taask-server/model"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
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

type QueueTaskRequest struct {
	Task                 *model.Task   `protobuf:"bytes,1,opt,name=Task,proto3" json:"Task,omitempty"`
	Session              *auth.Session `protobuf:"bytes,2,opt,name=Session,proto3" json:"Session,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *QueueTaskRequest) Reset()         { *m = QueueTaskRequest{} }
func (m *QueueTaskRequest) String() string { return proto.CompactTextString(m) }
func (*QueueTaskRequest) ProtoMessage()    {}
func (*QueueTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f9f4b593773f8c1, []int{0}
}

func (m *QueueTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueueTaskRequest.Unmarshal(m, b)
}
func (m *QueueTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueueTaskRequest.Marshal(b, m, deterministic)
}
func (m *QueueTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueueTaskRequest.Merge(m, src)
}
func (m *QueueTaskRequest) XXX_Size() int {
	return xxx_messageInfo_QueueTaskRequest.Size(m)
}
func (m *QueueTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueueTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueueTaskRequest proto.InternalMessageInfo

func (m *QueueTaskRequest) GetTask() *model.Task {
	if m != nil {
		return m.Task
	}
	return nil
}

func (m *QueueTaskRequest) GetSession() *auth.Session {
	if m != nil {
		return m.Session
	}
	return nil
}

type QueueTaskResponse struct {
	UUID                 string   `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueueTaskResponse) Reset()         { *m = QueueTaskResponse{} }
func (m *QueueTaskResponse) String() string { return proto.CompactTextString(m) }
func (*QueueTaskResponse) ProtoMessage()    {}
func (*QueueTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f9f4b593773f8c1, []int{1}
}

func (m *QueueTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueueTaskResponse.Unmarshal(m, b)
}
func (m *QueueTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueueTaskResponse.Marshal(b, m, deterministic)
}
func (m *QueueTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueueTaskResponse.Merge(m, src)
}
func (m *QueueTaskResponse) XXX_Size() int {
	return xxx_messageInfo_QueueTaskResponse.Size(m)
}
func (m *QueueTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueueTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueueTaskResponse proto.InternalMessageInfo

func (m *QueueTaskResponse) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

type CheckTaskRequest struct {
	UUID                 string        `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	Session              *auth.Session `protobuf:"bytes,2,opt,name=Session,proto3" json:"Session,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *CheckTaskRequest) Reset()         { *m = CheckTaskRequest{} }
func (m *CheckTaskRequest) String() string { return proto.CompactTextString(m) }
func (*CheckTaskRequest) ProtoMessage()    {}
func (*CheckTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f9f4b593773f8c1, []int{2}
}

func (m *CheckTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckTaskRequest.Unmarshal(m, b)
}
func (m *CheckTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckTaskRequest.Marshal(b, m, deterministic)
}
func (m *CheckTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckTaskRequest.Merge(m, src)
}
func (m *CheckTaskRequest) XXX_Size() int {
	return xxx_messageInfo_CheckTaskRequest.Size(m)
}
func (m *CheckTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CheckTaskRequest proto.InternalMessageInfo

func (m *CheckTaskRequest) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func (m *CheckTaskRequest) GetSession() *auth.Session {
	if m != nil {
		return m.Session
	}
	return nil
}

type CheckTaskResponse struct {
	Status               string               `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Progress             int32                `protobuf:"varint,2,opt,name=progress,proto3" json:"progress,omitempty"`
	EncTaskKey           *simplcrypto.Message `protobuf:"bytes,3,opt,name=EncTaskKey,proto3" json:"EncTaskKey,omitempty"`
	Result               *model.TaskUpdate    `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CheckTaskResponse) Reset()         { *m = CheckTaskResponse{} }
func (m *CheckTaskResponse) String() string { return proto.CompactTextString(m) }
func (*CheckTaskResponse) ProtoMessage()    {}
func (*CheckTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f9f4b593773f8c1, []int{3}
}

func (m *CheckTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckTaskResponse.Unmarshal(m, b)
}
func (m *CheckTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckTaskResponse.Marshal(b, m, deterministic)
}
func (m *CheckTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckTaskResponse.Merge(m, src)
}
func (m *CheckTaskResponse) XXX_Size() int {
	return xxx_messageInfo_CheckTaskResponse.Size(m)
}
func (m *CheckTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CheckTaskResponse proto.InternalMessageInfo

func (m *CheckTaskResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *CheckTaskResponse) GetProgress() int32 {
	if m != nil {
		return m.Progress
	}
	return 0
}

func (m *CheckTaskResponse) GetEncTaskKey() *simplcrypto.Message {
	if m != nil {
		return m.EncTaskKey
	}
	return nil
}

func (m *CheckTaskResponse) GetResult() *model.TaskUpdate {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*QueueTaskRequest)(nil), "taask.server.service.QueueTaskRequest")
	proto.RegisterType((*QueueTaskResponse)(nil), "taask.server.service.QueueTaskResponse")
	proto.RegisterType((*CheckTaskRequest)(nil), "taask.server.service.CheckTaskRequest")
	proto.RegisterType((*CheckTaskResponse)(nil), "taask.server.service.CheckTaskResponse")
}

func init() { proto.RegisterFile("taskservice.proto", fileDescriptor_2f9f4b593773f8c1) }

var fileDescriptor_2f9f4b593773f8c1 = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xdf, 0x8b, 0xd3, 0x40,
	0x10, 0xc7, 0xc9, 0xd9, 0xab, 0xde, 0x9c, 0x0f, 0xd7, 0xc5, 0x3b, 0x42, 0x1e, 0x44, 0x0a, 0x9a,
	0x7b, 0xd0, 0x4d, 0xa9, 0x3f, 0xde, 0xb5, 0xfa, 0x20, 0xd2, 0x07, 0x53, 0x0b, 0x22, 0xa2, 0x6c,
	0xd3, 0xa1, 0x09, 0xf9, 0xb1, 0x71, 0x67, 0xb7, 0xd8, 0xff, 0xc8, 0xbf, 0xc0, 0xbf, 0x4f, 0xb2,
	0x59, 0x4b, 0x5a, 0x5b, 0x0a, 0xbe, 0x84, 0xdd, 0xc9, 0xe7, 0x3b, 0xf3, 0x65, 0x66, 0x16, 0x06,
	0x5a, 0x50, 0x4e, 0xa8, 0xd6, 0x59, 0x82, 0xbc, 0x56, 0x52, 0x4b, 0xf6, 0x40, 0x0b, 0x41, 0x39,
	0x6f, 0x82, 0xa8, 0xb8, 0xfb, 0x17, 0xdc, 0x4f, 0x64, 0x59, 0xca, 0xaa, 0x65, 0x82, 0x6b, 0x61,
	0x74, 0x1a, 0xd9, 0x73, 0xd4, 0x1c, 0x5d, 0xf8, 0xa6, 0x94, 0x4b, 0x2c, 0x5c, 0xbc, 0xc9, 0xec,
	0xe2, 0xa3, 0x55, 0xa6, 0x53, 0xb3, 0xe0, 0x89, 0x2c, 0xa3, 0x44, 0xa6, 0xd9, 0xcf, 0x88, 0xb2,
	0xb2, 0x2e, 0x12, 0xb5, 0xa9, 0xb5, 0x74, 0x78, 0x89, 0x44, 0x62, 0xe5, 0x4c, 0x0c, 0xd7, 0x70,
	0xf5, 0xd1, 0xa0, 0xc1, 0x4f, 0x82, 0xf2, 0x18, 0x7f, 0x18, 0x24, 0xcd, 0x9e, 0x42, 0xaf, 0xb9,
	0xfa, 0xde, 0x23, 0xef, 0xf6, 0x72, 0xec, 0xf3, 0x1d, 0x9f, 0xb6, 0x32, 0xb7, 0xb8, 0xa5, 0xd8,
	0x0b, 0xb8, 0x3b, 0x43, 0xa2, 0x4c, 0x56, 0xfe, 0x99, 0x15, 0x04, 0xbb, 0x02, 0x6b, 0xdb, 0x11,
	0xf1, 0x5f, 0x74, 0x18, 0xc2, 0xa0, 0x53, 0x97, 0x6a, 0x59, 0x11, 0x32, 0x06, 0xbd, 0xf9, 0xfc,
	0xfd, 0x5b, 0x5b, 0xf8, 0x22, 0xb6, 0xe7, 0xe1, 0x57, 0xb8, 0x9a, 0xa4, 0x98, 0xe4, 0x5d, 0x83,
	0x07, 0xb8, 0xff, 0xb4, 0xf1, 0xdb, 0x83, 0x41, 0x27, 0xbd, 0xf3, 0x71, 0x03, 0x7d, 0xd2, 0x42,
	0x1b, 0x72, 0x15, 0xdc, 0x8d, 0x05, 0x70, 0xaf, 0x56, 0x72, 0xa5, 0x90, 0xc8, 0x16, 0x39, 0x8f,
	0xb7, 0x77, 0xf6, 0x12, 0xe0, 0x5d, 0x95, 0x34, 0x69, 0x3e, 0xe0, 0xc6, 0xbf, 0x63, 0x2d, 0x5c,
	0x73, 0xdb, 0x7e, 0xde, 0xf6, 0x9f, 0x4f, 0xdb, 0xce, 0xc7, 0x1d, 0x90, 0xbd, 0x82, 0xbe, 0x42,
	0x32, 0x85, 0xf6, 0x7b, 0x56, 0xf2, 0xf0, 0x58, 0xb7, 0xe7, 0xf5, 0x52, 0x68, 0x8c, 0x1d, 0x3d,
	0xfe, 0x75, 0x06, 0x97, 0x4d, 0x78, 0xd6, 0xae, 0x0d, 0xfb, 0x0e, 0xf0, 0xda, 0xe8, 0x74, 0x52,
	0x64, 0x58, 0x69, 0x16, 0xf2, 0x43, 0xbb, 0xc5, 0x1b, 0x62, 0x8a, 0xe5, 0x02, 0x95, 0xeb, 0x64,
	0x70, 0x7b, 0x1a, 0x74, 0x3d, 0xf9, 0x0c, 0xe7, 0x76, 0x60, 0xec, 0xc9, 0x61, 0xc9, 0xfe, 0x16,
	0x05, 0xe1, 0x49, 0xce, 0x65, 0xfe, 0x06, 0x17, 0xdb, 0x11, 0x1c, 0xcb, 0xbe, 0xbf, 0x02, 0xc7,
	0xb2, 0xff, 0x33, 0xcb, 0x91, 0xf7, 0x26, 0xfc, 0xf2, 0xb8, 0xf3, 0x2c, 0xac, 0xac, 0xfd, 0x3e,
	0x6b, 0xc5, 0x91, 0x13, 0x2f, 0xfa, 0xf6, 0x49, 0x3c, 0xff, 0x13, 0x00, 0x00, 0xff, 0xff, 0x34,
	0x86, 0x64, 0x3f, 0xac, 0x03, 0x00, 0x00,
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
	AuthClient(ctx context.Context, in *AuthMemberRequest, opts ...grpc.CallOption) (*AuthMemberResponse, error)
	Queue(ctx context.Context, in *QueueTaskRequest, opts ...grpc.CallOption) (*QueueTaskResponse, error)
	CheckTask(ctx context.Context, in *CheckTaskRequest, opts ...grpc.CallOption) (TaskService_CheckTaskClient, error)
}

type taskServiceClient struct {
	cc *grpc.ClientConn
}

func NewTaskServiceClient(cc *grpc.ClientConn) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) AuthClient(ctx context.Context, in *AuthMemberRequest, opts ...grpc.CallOption) (*AuthMemberResponse, error) {
	out := new(AuthMemberResponse)
	err := c.cc.Invoke(ctx, "/taask.server.service.TaskService/AuthClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) Queue(ctx context.Context, in *QueueTaskRequest, opts ...grpc.CallOption) (*QueueTaskResponse, error) {
	out := new(QueueTaskResponse)
	err := c.cc.Invoke(ctx, "/taask.server.service.TaskService/Queue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) CheckTask(ctx context.Context, in *CheckTaskRequest, opts ...grpc.CallOption) (TaskService_CheckTaskClient, error) {
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
	Recv() (*CheckTaskResponse, error)
	grpc.ClientStream
}

type taskServiceCheckTaskClient struct {
	grpc.ClientStream
}

func (x *taskServiceCheckTaskClient) Recv() (*CheckTaskResponse, error) {
	m := new(CheckTaskResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TaskServiceServer is the server API for TaskService service.
type TaskServiceServer interface {
	AuthClient(context.Context, *AuthMemberRequest) (*AuthMemberResponse, error)
	Queue(context.Context, *QueueTaskRequest) (*QueueTaskResponse, error)
	CheckTask(*CheckTaskRequest, TaskService_CheckTaskServer) error
}

func RegisterTaskServiceServer(s *grpc.Server, srv TaskServiceServer) {
	s.RegisterService(&_TaskService_serviceDesc, srv)
}

func _TaskService_AuthClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).AuthClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/taask.server.service.TaskService/AuthClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).AuthClient(ctx, req.(*AuthMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_Queue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueTaskRequest)
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
		return srv.(TaskServiceServer).Queue(ctx, req.(*QueueTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_CheckTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CheckTaskRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskServiceServer).CheckTask(m, &taskServiceCheckTaskServer{stream})
}

type TaskService_CheckTaskServer interface {
	Send(*CheckTaskResponse) error
	grpc.ServerStream
}

type taskServiceCheckTaskServer struct {
	grpc.ServerStream
}

func (x *taskServiceCheckTaskServer) Send(m *CheckTaskResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _TaskService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "taask.server.service.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthClient",
			Handler:    _TaskService_AuthClient_Handler,
		},
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
	Metadata: "taskservice.proto",
}