// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task.proto

package model

import (
	proto1 "crypto/proto"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Task struct {
	UUID                 string                     `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	Type                 string                     `protobuf:"bytes,2,opt,name=Type,proto3" json:"Type,omitempty"`
	PubKey               *proto1.SerializablePubKey `protobuf:"bytes,3,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	Body                 []byte                     `protobuf:"bytes,4,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{0}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func (m *Task) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Task) GetPubKey() *proto1.SerializablePubKey {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Task) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type TaskResult struct {
	UUID                 string          `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	EncResponse          *proto1.Message `protobuf:"bytes,2,opt,name=EncResponse,proto3" json:"EncResponse,omitempty"`
	EncResponseSymKey    *proto1.Message `protobuf:"bytes,3,opt,name=EncResponseSymKey,proto3" json:"EncResponseSymKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *TaskResult) Reset()         { *m = TaskResult{} }
func (m *TaskResult) String() string { return proto.CompactTextString(m) }
func (*TaskResult) ProtoMessage()    {}
func (*TaskResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{1}
}

func (m *TaskResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskResult.Unmarshal(m, b)
}
func (m *TaskResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskResult.Marshal(b, m, deterministic)
}
func (m *TaskResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskResult.Merge(m, src)
}
func (m *TaskResult) XXX_Size() int {
	return xxx_messageInfo_TaskResult.Size(m)
}
func (m *TaskResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskResult.DiscardUnknown(m)
}

var xxx_messageInfo_TaskResult proto.InternalMessageInfo

func (m *TaskResult) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func (m *TaskResult) GetEncResponse() *proto1.Message {
	if m != nil {
		return m.EncResponse
	}
	return nil
}

func (m *TaskResult) GetEncResponseSymKey() *proto1.Message {
	if m != nil {
		return m.EncResponseSymKey
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
	return fileDescriptor_ce5d8dd45b4a91ff, []int{2}
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
	UUID                 string   `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckTaskRequest) Reset()         { *m = CheckTaskRequest{} }
func (m *CheckTaskRequest) String() string { return proto.CompactTextString(m) }
func (*CheckTaskRequest) ProtoMessage()    {}
func (*CheckTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{3}
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

type CheckTaskResponse struct {
	Completed            bool        `protobuf:"varint,1,opt,name=completed,proto3" json:"completed,omitempty"`
	Progress             int32       `protobuf:"varint,2,opt,name=progress,proto3" json:"progress,omitempty"`
	Result               *TaskResult `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *CheckTaskResponse) Reset()         { *m = CheckTaskResponse{} }
func (m *CheckTaskResponse) String() string { return proto.CompactTextString(m) }
func (*CheckTaskResponse) ProtoMessage()    {}
func (*CheckTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{4}
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

func (m *CheckTaskResponse) GetCompleted() bool {
	if m != nil {
		return m.Completed
	}
	return false
}

func (m *CheckTaskResponse) GetProgress() int32 {
	if m != nil {
		return m.Progress
	}
	return 0
}

func (m *CheckTaskResponse) GetResult() *TaskResult {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*Task)(nil), "taask.server.model.Task")
	proto.RegisterType((*TaskResult)(nil), "taask.server.model.TaskResult")
	proto.RegisterType((*QueueTaskResponse)(nil), "taask.server.model.QueueTaskResponse")
	proto.RegisterType((*CheckTaskRequest)(nil), "taask.server.model.CheckTaskRequest")
	proto.RegisterType((*CheckTaskResponse)(nil), "taask.server.model.CheckTaskResponse")
}

func init() { proto.RegisterFile("task.proto", fileDescriptor_ce5d8dd45b4a91ff) }

var fileDescriptor_ce5d8dd45b4a91ff = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0xc1, 0x4e, 0x02, 0x31,
	0x10, 0x4d, 0x15, 0x10, 0x06, 0x0f, 0xd2, 0xd3, 0x86, 0x10, 0x43, 0xf6, 0x20, 0x9c, 0x96, 0x04,
	0x13, 0x8f, 0x9a, 0xa0, 0x1e, 0xd4, 0x98, 0x68, 0x85, 0x8b, 0xb7, 0xb2, 0x4c, 0x90, 0xec, 0x2e,
	0xad, 0xed, 0xae, 0x49, 0xbd, 0x1a, 0xff, 0xc6, 0x8f, 0x34, 0x6d, 0x37, 0x82, 0xb2, 0xf1, 0x36,
	0x6f, 0xe6, 0xcd, 0xcc, 0x7b, 0x93, 0x01, 0xc8, 0xb9, 0x4e, 0x22, 0xa9, 0x44, 0x2e, 0x28, 0xcd,
	0xb9, 0x05, 0x1a, 0xd5, 0x1b, 0xaa, 0x28, 0x13, 0x0b, 0x4c, 0xbb, 0xdd, 0x58, 0x19, 0x99, 0x8b,
	0x91, 0x63, 0x8c, 0x32, 0xd4, 0x9a, 0x2f, 0xd1, 0xf3, 0xff, 0xd4, 0x12, 0x34, 0x92, 0xaf, 0x94,
	0xaf, 0x85, 0x1f, 0x04, 0x6a, 0x53, 0xae, 0x13, 0x4a, 0xa1, 0x36, 0x9b, 0xdd, 0x5c, 0x05, 0xa4,
	0x4f, 0x86, 0x2d, 0xe6, 0x62, 0x9b, 0x9b, 0x1a, 0x89, 0xc1, 0x9e, 0xcf, 0xd9, 0x98, 0x5e, 0x40,
	0xe3, 0xa1, 0x98, 0xdf, 0xa1, 0x09, 0xf6, 0xfb, 0x64, 0xd8, 0x1e, 0x0f, 0xa2, 0x5f, 0x6a, 0xfc,
	0x2a, 0x8b, 0x56, 0x3c, 0x5d, 0xbd, 0xf3, 0x79, 0x8a, 0x9e, 0xce, 0xca, 0x36, 0x3b, 0x74, 0x22,
	0x16, 0x26, 0xa8, 0xf5, 0xc9, 0xf0, 0x90, 0xb9, 0x38, 0xfc, 0x22, 0x00, 0x56, 0x05, 0x43, 0x5d,
	0xa4, 0x79, 0xa5, 0x96, 0x73, 0x68, 0x5f, 0xaf, 0x63, 0x86, 0x5a, 0x8a, 0xb5, 0xf6, 0x92, 0xda,
	0xe3, 0x5e, 0xe5, 0xf2, 0x7b, 0xef, 0x9e, 0x6d, 0x37, 0xd0, 0x5b, 0xe8, 0x6c, 0xc1, 0x27, 0x93,
	0x6d, 0x2c, 0xfc, 0x3f, 0x65, 0xb7, 0x2d, 0x1c, 0x40, 0xe7, 0xb1, 0xc0, 0x02, 0x4b, 0xc9, 0x7e,
	0x41, 0x85, 0xe8, 0xf0, 0x04, 0x8e, 0x2e, 0x5f, 0x30, 0x4e, 0x3c, 0xf1, 0xb5, 0x40, 0x5d, 0x69,
	0x2e, 0xfc, 0x24, 0xd0, 0xd9, 0x22, 0x96, 0x13, 0x7b, 0xd0, 0x8a, 0x45, 0x26, 0x53, 0xcc, 0x71,
	0xe1, 0xe8, 0x4d, 0xb6, 0x49, 0xd0, 0x2e, 0x34, 0xa5, 0x12, 0x4b, 0x85, 0x5a, 0xbb, 0x6b, 0xd4,
	0xd9, 0x0f, 0xa6, 0x67, 0xd0, 0x50, 0xee, 0x94, 0xa5, 0xc3, 0xe3, 0x68, 0xf7, 0x65, 0xa2, 0xcd,
	0xc1, 0x59, 0xc9, 0x9e, 0x1c, 0x3c, 0xd7, 0x5d, 0x6d, 0xde, 0x70, 0xdf, 0x71, 0xfa, 0x1d, 0x00,
	0x00, 0xff, 0xff, 0xed, 0x2c, 0x82, 0xaf, 0x77, 0x02, 0x00, 0x00,
}