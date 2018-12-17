// Code generated by protoc-gen-go. DO NOT EDIT.
// source: runnerservice.proto

package service

import (
	fmt "fmt"
	simplcrypto "github.com/cohix/simplcrypto"
	proto "github.com/golang/protobuf/proto"
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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_055c2bb08edcaeb3, []int{0}
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

type AuthRunnerRequest struct {
	PubKey               *simplcrypto.SerializablePubKey `protobuf:"bytes,1,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	JoinCodeSignature    *simplcrypto.Signature          `protobuf:"bytes,2,opt,name=JoinCodeSignature,proto3" json:"JoinCodeSignature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *AuthRunnerRequest) Reset()         { *m = AuthRunnerRequest{} }
func (m *AuthRunnerRequest) String() string { return proto.CompactTextString(m) }
func (*AuthRunnerRequest) ProtoMessage()    {}
func (*AuthRunnerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_055c2bb08edcaeb3, []int{1}
}

func (m *AuthRunnerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthRunnerRequest.Unmarshal(m, b)
}
func (m *AuthRunnerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthRunnerRequest.Marshal(b, m, deterministic)
}
func (m *AuthRunnerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthRunnerRequest.Merge(m, src)
}
func (m *AuthRunnerRequest) XXX_Size() int {
	return xxx_messageInfo_AuthRunnerRequest.Size(m)
}
func (m *AuthRunnerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthRunnerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AuthRunnerRequest proto.InternalMessageInfo

func (m *AuthRunnerRequest) GetPubKey() *simplcrypto.SerializablePubKey {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *AuthRunnerRequest) GetJoinCodeSignature() *simplcrypto.Signature {
	if m != nil {
		return m.JoinCodeSignature
	}
	return nil
}

type AuthRunnerResponse struct {
	EncChallenge         *simplcrypto.Message `protobuf:"bytes,1,opt,name=EncChallenge,proto3" json:"EncChallenge,omitempty"`
	EncChallengeKey      *simplcrypto.Message `protobuf:"bytes,2,opt,name=EncChallengeKey,proto3" json:"EncChallengeKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *AuthRunnerResponse) Reset()         { *m = AuthRunnerResponse{} }
func (m *AuthRunnerResponse) String() string { return proto.CompactTextString(m) }
func (*AuthRunnerResponse) ProtoMessage()    {}
func (*AuthRunnerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_055c2bb08edcaeb3, []int{2}
}

func (m *AuthRunnerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthRunnerResponse.Unmarshal(m, b)
}
func (m *AuthRunnerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthRunnerResponse.Marshal(b, m, deterministic)
}
func (m *AuthRunnerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthRunnerResponse.Merge(m, src)
}
func (m *AuthRunnerResponse) XXX_Size() int {
	return xxx_messageInfo_AuthRunnerResponse.Size(m)
}
func (m *AuthRunnerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthRunnerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AuthRunnerResponse proto.InternalMessageInfo

func (m *AuthRunnerResponse) GetEncChallenge() *simplcrypto.Message {
	if m != nil {
		return m.EncChallenge
	}
	return nil
}

func (m *AuthRunnerResponse) GetEncChallengeKey() *simplcrypto.Message {
	if m != nil {
		return m.EncChallengeKey
	}
	return nil
}

type RegisterRunnerRequest struct {
	UUID                 string                 `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	Kind                 string                 `protobuf:"bytes,2,opt,name=Kind,proto3" json:"Kind,omitempty"`
	Tags                 []string               `protobuf:"bytes,3,rep,name=Tags,proto3" json:"Tags,omitempty"`
	ChallengeSignature   *simplcrypto.Signature `protobuf:"bytes,4,opt,name=ChallengeSignature,proto3" json:"ChallengeSignature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *RegisterRunnerRequest) Reset()         { *m = RegisterRunnerRequest{} }
func (m *RegisterRunnerRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRunnerRequest) ProtoMessage()    {}
func (*RegisterRunnerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_055c2bb08edcaeb3, []int{3}
}

func (m *RegisterRunnerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRunnerRequest.Unmarshal(m, b)
}
func (m *RegisterRunnerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRunnerRequest.Marshal(b, m, deterministic)
}
func (m *RegisterRunnerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRunnerRequest.Merge(m, src)
}
func (m *RegisterRunnerRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRunnerRequest.Size(m)
}
func (m *RegisterRunnerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRunnerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRunnerRequest proto.InternalMessageInfo

func (m *RegisterRunnerRequest) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func (m *RegisterRunnerRequest) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *RegisterRunnerRequest) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *RegisterRunnerRequest) GetChallengeSignature() *simplcrypto.Signature {
	if m != nil {
		return m.ChallengeSignature
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "taask.server.service.Empty")
	proto.RegisterType((*AuthRunnerRequest)(nil), "taask.server.service.AuthRunnerRequest")
	proto.RegisterType((*AuthRunnerResponse)(nil), "taask.server.service.AuthRunnerResponse")
	proto.RegisterType((*RegisterRunnerRequest)(nil), "taask.server.service.RegisterRunnerRequest")
}

func init() { proto.RegisterFile("runnerservice.proto", fileDescriptor_055c2bb08edcaeb3) }

var fileDescriptor_055c2bb08edcaeb3 = []byte{
	// 438 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xdf, 0x6b, 0xd5, 0x30,
	0x14, 0xc7, 0xe9, 0xdd, 0x9c, 0xec, 0xf8, 0x8b, 0x45, 0xa7, 0x97, 0x0a, 0x72, 0xb9, 0x20, 0x1b,
	0x88, 0xe9, 0x98, 0x2f, 0xfa, 0x24, 0x3a, 0x2f, 0x32, 0x87, 0x20, 0xb9, 0xbb, 0x0f, 0xfa, 0x22,
	0xb9, 0xed, 0xa1, 0x0d, 0xb7, 0x4d, 0x6a, 0x92, 0x8a, 0xf5, 0xaf, 0xf0, 0x41, 0x7c, 0xf2, 0x8f,
	0x95, 0x24, 0xd5, 0xdd, 0x6e, 0x75, 0xe2, 0x4b, 0x39, 0x3d, 0xf9, 0x9c, 0x93, 0x6f, 0xbe, 0xc9,
	0x81, 0xdb, 0xba, 0x91, 0x12, 0xb5, 0x41, 0xfd, 0x59, 0xa4, 0x48, 0x6b, 0xad, 0xac, 0x22, 0x77,
	0x2c, 0xe7, 0x66, 0x45, 0x5d, 0x12, 0x35, 0xed, 0xd6, 0xe2, 0xbb, 0x95, 0xca, 0xb0, 0x4c, 0x3c,
	0x92, 0x58, 0x07, 0xf8, 0x30, 0x3e, 0xc8, 0x85, 0x2d, 0x9a, 0x25, 0x4d, 0x55, 0x95, 0xa4, 0xaa,
	0x10, 0x5f, 0x12, 0x23, 0xaa, 0xba, 0x4c, 0x75, 0x5b, 0x5b, 0xd5, 0xe1, 0x2b, 0x6c, 0x6b, 0x2e,
	0xf4, 0x7f, 0x54, 0x54, 0x68, 0x0c, 0xcf, 0x3b, 0x45, 0xd3, 0xab, 0x70, 0x65, 0x56, 0xd5, 0xb6,
	0x9d, 0x7e, 0x8f, 0x60, 0xe7, 0x45, 0x63, 0x0b, 0xe6, 0x65, 0x33, 0xfc, 0xd4, 0xa0, 0xb1, 0xe4,
	0x29, 0x6c, 0xbd, 0x6b, 0x96, 0x27, 0xd8, 0x8e, 0xa3, 0x49, 0xb4, 0x7f, 0xed, 0x70, 0x42, 0x7d,
	0x43, 0x1a, 0x3a, 0xd2, 0x39, 0x6a, 0xc1, 0x4b, 0xf1, 0x95, 0x2f, 0x4b, 0x0c, 0x1c, 0xeb, 0x78,
	0x32, 0x83, 0x9d, 0x37, 0x4a, 0xc8, 0x23, 0x95, 0xe1, 0x5c, 0xe4, 0x92, 0xdb, 0x46, 0xe3, 0x78,
	0xe4, 0x9b, 0xdc, 0x3b, 0xd7, 0xe4, 0xf7, 0x32, 0xbb, 0x58, 0x31, 0xfd, 0x16, 0x01, 0x59, 0x97,
	0x65, 0x6a, 0x25, 0x0d, 0x92, 0x67, 0x70, 0x7d, 0x26, 0xd3, 0xa3, 0x82, 0x97, 0x25, 0xca, 0x1c,
	0x3b, 0x75, 0xbb, 0xfd, 0xc6, 0x6f, 0xc3, 0x49, 0x59, 0x0f, 0x25, 0xcf, 0xe1, 0xd6, 0xfa, 0xbf,
	0x3b, 0xdb, 0xe8, 0xb2, 0xea, 0xf3, 0xf4, 0xf4, 0x67, 0x04, 0xbb, 0x0c, 0x73, 0x61, 0x2c, 0xea,
	0xbe, 0x5b, 0x04, 0x36, 0x17, 0x8b, 0xe3, 0x57, 0x5e, 0xcd, 0x36, 0xf3, 0xb1, 0xcb, 0x9d, 0x08,
	0x99, 0xf9, 0x3d, 0xb6, 0x99, 0x8f, 0x5d, 0xee, 0x94, 0xe7, 0x66, 0xbc, 0x31, 0xd9, 0x70, 0x39,
	0x17, 0x93, 0xd7, 0x40, 0xfe, 0xec, 0x72, 0x66, 0xd8, 0xe6, 0xe5, 0x86, 0x0d, 0x94, 0x1c, 0xfe,
	0x18, 0xc1, 0x8d, 0x20, 0x6b, 0x1e, 0xde, 0x17, 0xf9, 0x08, 0x70, 0x66, 0x21, 0xd9, 0xa3, 0x43,
	0x8f, 0x90, 0x5e, 0xb8, 0xfb, 0x78, 0xff, 0xdf, 0x60, 0x77, 0x1b, 0xef, 0xe1, 0x66, 0xdf, 0x10,
	0xf2, 0x68, 0xb8, 0x76, 0xd0, 0xb6, 0x78, 0xdc, 0x87, 0xfd, 0x34, 0xd0, 0x53, 0x6e, 0x56, 0x07,
	0x11, 0x39, 0x06, 0x58, 0xd4, 0x19, 0xb7, 0xe8, 0xfe, 0xc9, 0x83, 0xbf, 0x91, 0x81, 0x89, 0xef,
	0x0f, 0x6f, 0xeb, 0x5f, 0xf8, 0xcb, 0xbd, 0x0f, 0x0f, 0xd7, 0xc6, 0xc3, 0x83, 0xe1, 0xfb, 0x38,
	0xe0, 0x49, 0x87, 0x2f, 0xb7, 0xfc, 0x68, 0x3c, 0xf9, 0x15, 0x00, 0x00, 0xff, 0xff, 0xce, 0x13,
	0xcb, 0x68, 0xc3, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RunnerServiceClient is the client API for RunnerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RunnerServiceClient interface {
	AuthRunner(ctx context.Context, in *AuthRunnerRequest, opts ...grpc.CallOption) (*AuthRunnerResponse, error)
	RegisterRunner(ctx context.Context, in *RegisterRunnerRequest, opts ...grpc.CallOption) (RunnerService_RegisterRunnerClient, error)
	UpdateTask(ctx context.Context, in *model.TaskUpdate, opts ...grpc.CallOption) (*Empty, error)
}

type runnerServiceClient struct {
	cc *grpc.ClientConn
}

func NewRunnerServiceClient(cc *grpc.ClientConn) RunnerServiceClient {
	return &runnerServiceClient{cc}
}

func (c *runnerServiceClient) AuthRunner(ctx context.Context, in *AuthRunnerRequest, opts ...grpc.CallOption) (*AuthRunnerResponse, error) {
	out := new(AuthRunnerResponse)
	err := c.cc.Invoke(ctx, "/taask.server.service.RunnerService/AuthRunner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerServiceClient) RegisterRunner(ctx context.Context, in *RegisterRunnerRequest, opts ...grpc.CallOption) (RunnerService_RegisterRunnerClient, error) {
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
	AuthRunner(context.Context, *AuthRunnerRequest) (*AuthRunnerResponse, error)
	RegisterRunner(*RegisterRunnerRequest, RunnerService_RegisterRunnerServer) error
	UpdateTask(context.Context, *model.TaskUpdate) (*Empty, error)
}

func RegisterRunnerServiceServer(s *grpc.Server, srv RunnerServiceServer) {
	s.RegisterService(&_RunnerService_serviceDesc, srv)
}

func _RunnerService_AuthRunner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthRunnerRequest)
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
		return srv.(RunnerServiceServer).AuthRunner(ctx, req.(*AuthRunnerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RunnerService_RegisterRunner_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RegisterRunnerRequest)
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
	Metadata: "runnerservice.proto",
}
