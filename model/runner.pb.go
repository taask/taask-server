// Code generated by protoc-gen-go. DO NOT EDIT.
// source: runner.proto

package model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	crypto "github.com/taask/taask-server/crypto"
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

type RegisterRunnerRequest struct {
	UUID                 string            `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	Type                 string            `protobuf:"bytes,2,opt,name=Type,proto3" json:"Type,omitempty"`
	PubKey               string            `protobuf:"bytes,3,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	JoinCodeSignature    *crypto.Signature `protobuf:"bytes,4,opt,name=JoinCodeSignature,proto3" json:"JoinCodeSignature,omitempty"`
	Tags                 []string          `protobuf:"bytes,5,rep,name=Tags,proto3" json:"Tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RegisterRunnerRequest) Reset()         { *m = RegisterRunnerRequest{} }
func (m *RegisterRunnerRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRunnerRequest) ProtoMessage()    {}
func (*RegisterRunnerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{0}
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

func (m *RegisterRunnerRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *RegisterRunnerRequest) GetPubKey() string {
	if m != nil {
		return m.PubKey
	}
	return ""
}

func (m *RegisterRunnerRequest) GetJoinCodeSignature() *crypto.Signature {
	if m != nil {
		return m.JoinCodeSignature
	}
	return nil
}

func (m *RegisterRunnerRequest) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

type RegisterRunnerResponse struct {
	EncChallenge         string          `protobuf:"bytes,1,opt,name=EncChallenge,proto3" json:"EncChallenge,omitempty"`
	EncGlobalKey         *crypto.Message `protobuf:"bytes,2,opt,name=EncGlobalKey,proto3" json:"EncGlobalKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *RegisterRunnerResponse) Reset()         { *m = RegisterRunnerResponse{} }
func (m *RegisterRunnerResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterRunnerResponse) ProtoMessage()    {}
func (*RegisterRunnerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{1}
}

func (m *RegisterRunnerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRunnerResponse.Unmarshal(m, b)
}
func (m *RegisterRunnerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRunnerResponse.Marshal(b, m, deterministic)
}
func (m *RegisterRunnerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRunnerResponse.Merge(m, src)
}
func (m *RegisterRunnerResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterRunnerResponse.Size(m)
}
func (m *RegisterRunnerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRunnerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRunnerResponse proto.InternalMessageInfo

func (m *RegisterRunnerResponse) GetEncChallenge() string {
	if m != nil {
		return m.EncChallenge
	}
	return ""
}

func (m *RegisterRunnerResponse) GetEncGlobalKey() *crypto.Message {
	if m != nil {
		return m.EncGlobalKey
	}
	return nil
}

type StreamTasksRequest struct {
	ChallengeSignature   *crypto.Signature `protobuf:"bytes,1,opt,name=ChallengeSignature,proto3" json:"ChallengeSignature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *StreamTasksRequest) Reset()         { *m = StreamTasksRequest{} }
func (m *StreamTasksRequest) String() string { return proto.CompactTextString(m) }
func (*StreamTasksRequest) ProtoMessage()    {}
func (*StreamTasksRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_48eceea7e2abc593, []int{2}
}

func (m *StreamTasksRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamTasksRequest.Unmarshal(m, b)
}
func (m *StreamTasksRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamTasksRequest.Marshal(b, m, deterministic)
}
func (m *StreamTasksRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamTasksRequest.Merge(m, src)
}
func (m *StreamTasksRequest) XXX_Size() int {
	return xxx_messageInfo_StreamTasksRequest.Size(m)
}
func (m *StreamTasksRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamTasksRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StreamTasksRequest proto.InternalMessageInfo

func (m *StreamTasksRequest) GetChallengeSignature() *crypto.Signature {
	if m != nil {
		return m.ChallengeSignature
	}
	return nil
}

func init() {
	proto.RegisterType((*RegisterRunnerRequest)(nil), "taask.server.model.RegisterRunnerRequest")
	proto.RegisterType((*RegisterRunnerResponse)(nil), "taask.server.model.RegisterRunnerResponse")
	proto.RegisterType((*StreamTasksRequest)(nil), "taask.server.model.StreamTasksRequest")
}

func init() { proto.RegisterFile("runner.proto", fileDescriptor_48eceea7e2abc593) }

var fileDescriptor_48eceea7e2abc593 = []byte{
	// 309 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0x5f, 0x4b, 0xfb, 0x30,
	0x14, 0xa5, 0xdb, 0x7e, 0x83, 0x5f, 0xdc, 0x8b, 0x01, 0x47, 0x19, 0x22, 0xa3, 0x22, 0xec, 0xc5,
	0x14, 0xf4, 0x0b, 0x88, 0x53, 0xc4, 0xbf, 0x48, 0xb6, 0xbd, 0xf8, 0x96, 0xae, 0x97, 0xac, 0xac,
	0x4d, 0x6a, 0x6e, 0x2a, 0xf4, 0xc5, 0xef, 0xe5, 0xb7, 0x93, 0x26, 0x73, 0x52, 0xdd, 0x83, 0x2f,
	0xe1, 0xe6, 0x70, 0xce, 0xbd, 0xe7, 0xdc, 0x4b, 0x06, 0xa6, 0x52, 0x0a, 0x0c, 0x2b, 0x8d, 0xb6,
	0x9a, 0x52, 0x2b, 0x04, 0xae, 0x19, 0x82, 0x79, 0x03, 0xc3, 0x0a, 0x9d, 0x42, 0x3e, 0x1a, 0x2d,
	0x4d, 0x5d, 0x5a, 0x1d, 0x3b, 0x46, 0x5c, 0x00, 0xa2, 0x90, 0xe0, 0xf9, 0xd1, 0x47, 0x40, 0x0e,
	0x38, 0xc8, 0x0c, 0x2d, 0x18, 0xee, 0x1a, 0x71, 0x78, 0xad, 0x00, 0x2d, 0xa5, 0xa4, 0xb7, 0x58,
	0xdc, 0x5e, 0x85, 0xc1, 0x38, 0x98, 0xfc, 0xe7, 0xae, 0x6e, 0xb0, 0x79, 0x5d, 0x42, 0xd8, 0xf1,
	0x58, 0x53, 0xd3, 0x21, 0xe9, 0x3f, 0x57, 0xc9, 0x3d, 0xd4, 0x61, 0xd7, 0xa1, 0x9b, 0x1f, 0x7d,
	0x20, 0xfb, 0x77, 0x3a, 0x53, 0x53, 0x9d, 0xc2, 0x2c, 0x93, 0x4a, 0xd8, 0xca, 0x40, 0xd8, 0x1b,
	0x07, 0x93, 0xbd, 0xb3, 0x23, 0xd6, 0x72, 0xe9, 0xed, 0xb1, 0x2d, 0x8b, 0xff, 0x16, 0xba, 0xc9,
	0x42, 0x62, 0xf8, 0x6f, 0xdc, 0x75, 0x93, 0x85, 0xc4, 0xe8, 0x9d, 0x0c, 0x7f, 0x5a, 0xc7, 0x52,
	0x2b, 0x04, 0x1a, 0x91, 0xc1, 0xb5, 0x5a, 0x4e, 0x57, 0x22, 0xcf, 0x41, 0x49, 0xd8, 0x64, 0x68,
	0x61, 0xf4, 0xc2, 0x71, 0x6e, 0x72, 0x9d, 0x88, 0xbc, 0x71, 0xdf, 0x71, 0xd6, 0x0e, 0x77, 0x5a,
	0x7b, 0xf4, 0x3b, 0xe3, 0x2d, 0x45, 0x94, 0x12, 0x3a, 0xb3, 0x06, 0x44, 0x31, 0x17, 0xb8, 0xc6,
	0xaf, 0xbd, 0x3d, 0x11, 0xba, 0x1d, 0xf2, 0x1d, 0x3c, 0xf8, 0x53, 0xf0, 0x1d, 0xca, 0xcb, 0x93,
	0x97, 0x63, 0x99, 0xd9, 0x55, 0x95, 0xb0, 0xa5, 0x2e, 0x62, 0xa7, 0xf7, 0xef, 0xa9, 0xef, 0x12,
	0xbb, 0x23, 0x27, 0x7d, 0x77, 0xcf, 0xf3, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x7a, 0x12,
	0x19, 0x0f, 0x02, 0x00, 0x00,
}
