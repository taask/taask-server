// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package auth

import (
	fmt "fmt"
	simplcrypto "github.com/cohix/simplcrypto"
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

type Attempt struct {
	MemberUUID           string                          `protobuf:"bytes,1,opt,name=MemberUUID,proto3" json:"MemberUUID,omitempty"`
	GroupUUID            string                          `protobuf:"bytes,2,opt,name=GroupUUID,proto3" json:"GroupUUID,omitempty"`
	PubKey               *simplcrypto.SerializablePubKey `protobuf:"bytes,3,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	AuthHashSignature    *simplcrypto.Signature          `protobuf:"bytes,4,opt,name=AuthHashSignature,proto3" json:"AuthHashSignature,omitempty"`
	Timestamp            int64                           `protobuf:"varint,5,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *Attempt) Reset()         { *m = Attempt{} }
func (m *Attempt) String() string { return proto.CompactTextString(m) }
func (*Attempt) ProtoMessage()    {}
func (*Attempt) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

func (m *Attempt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Attempt.Unmarshal(m, b)
}
func (m *Attempt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Attempt.Marshal(b, m, deterministic)
}
func (m *Attempt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Attempt.Merge(m, src)
}
func (m *Attempt) XXX_Size() int {
	return xxx_messageInfo_Attempt.Size(m)
}
func (m *Attempt) XXX_DiscardUnknown() {
	xxx_messageInfo_Attempt.DiscardUnknown(m)
}

var xxx_messageInfo_Attempt proto.InternalMessageInfo

func (m *Attempt) GetMemberUUID() string {
	if m != nil {
		return m.MemberUUID
	}
	return ""
}

func (m *Attempt) GetGroupUUID() string {
	if m != nil {
		return m.GroupUUID
	}
	return ""
}

func (m *Attempt) GetPubKey() *simplcrypto.SerializablePubKey {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Attempt) GetAuthHashSignature() *simplcrypto.Signature {
	if m != nil {
		return m.AuthHashSignature
	}
	return nil
}

func (m *Attempt) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type AttemptResponse struct {
	EncChallenge         *simplcrypto.Message            `protobuf:"bytes,1,opt,name=EncChallenge,proto3" json:"EncChallenge,omitempty"`
	MasterPubKey         *simplcrypto.SerializablePubKey `protobuf:"bytes,2,opt,name=MasterPubKey,proto3" json:"MasterPubKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *AttemptResponse) Reset()         { *m = AttemptResponse{} }
func (m *AttemptResponse) String() string { return proto.CompactTextString(m) }
func (*AttemptResponse) ProtoMessage()    {}
func (*AttemptResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{1}
}

func (m *AttemptResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AttemptResponse.Unmarshal(m, b)
}
func (m *AttemptResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AttemptResponse.Marshal(b, m, deterministic)
}
func (m *AttemptResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttemptResponse.Merge(m, src)
}
func (m *AttemptResponse) XXX_Size() int {
	return xxx_messageInfo_AttemptResponse.Size(m)
}
func (m *AttemptResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AttemptResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AttemptResponse proto.InternalMessageInfo

func (m *AttemptResponse) GetEncChallenge() *simplcrypto.Message {
	if m != nil {
		return m.EncChallenge
	}
	return nil
}

func (m *AttemptResponse) GetMasterPubKey() *simplcrypto.SerializablePubKey {
	if m != nil {
		return m.MasterPubKey
	}
	return nil
}

type Session struct {
	MemberUUID           string                 `protobuf:"bytes,1,opt,name=MemberUUID,proto3" json:"MemberUUID,omitempty"`
	GroupUUID            string                 `protobuf:"bytes,2,opt,name=GroupUUID,proto3" json:"GroupUUID,omitempty"`
	SessionChallengeSig  *simplcrypto.Signature `protobuf:"bytes,3,opt,name=SessionChallengeSig,proto3" json:"SessionChallengeSig,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()    {}
func (*Session) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{2}
}

func (m *Session) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Session.Unmarshal(m, b)
}
func (m *Session) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Session.Marshal(b, m, deterministic)
}
func (m *Session) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Session.Merge(m, src)
}
func (m *Session) XXX_Size() int {
	return xxx_messageInfo_Session.Size(m)
}
func (m *Session) XXX_DiscardUnknown() {
	xxx_messageInfo_Session.DiscardUnknown(m)
}

var xxx_messageInfo_Session proto.InternalMessageInfo

func (m *Session) GetMemberUUID() string {
	if m != nil {
		return m.MemberUUID
	}
	return ""
}

func (m *Session) GetGroupUUID() string {
	if m != nil {
		return m.GroupUUID
	}
	return ""
}

func (m *Session) GetSessionChallengeSig() *simplcrypto.Signature {
	if m != nil {
		return m.SessionChallengeSig
	}
	return nil
}

type MemberGroup struct {
	UUID                 string   `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	JoinCode             string   `protobuf:"bytes,3,opt,name=JoinCode,proto3" json:"JoinCode,omitempty"`
	AuthHash             []byte   `protobuf:"bytes,4,opt,name=AuthHash,proto3" json:"AuthHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MemberGroup) Reset()         { *m = MemberGroup{} }
func (m *MemberGroup) String() string { return proto.CompactTextString(m) }
func (*MemberGroup) ProtoMessage()    {}
func (*MemberGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{3}
}

func (m *MemberGroup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MemberGroup.Unmarshal(m, b)
}
func (m *MemberGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MemberGroup.Marshal(b, m, deterministic)
}
func (m *MemberGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MemberGroup.Merge(m, src)
}
func (m *MemberGroup) XXX_Size() int {
	return xxx_messageInfo_MemberGroup.Size(m)
}
func (m *MemberGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_MemberGroup.DiscardUnknown(m)
}

var xxx_messageInfo_MemberGroup proto.InternalMessageInfo

func (m *MemberGroup) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func (m *MemberGroup) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MemberGroup) GetJoinCode() string {
	if m != nil {
		return m.JoinCode
	}
	return ""
}

func (m *MemberGroup) GetAuthHash() []byte {
	if m != nil {
		return m.AuthHash
	}
	return nil
}

func init() {
	proto.RegisterType((*Attempt)(nil), "taask.server.auth.Attempt")
	proto.RegisterType((*AttemptResponse)(nil), "taask.server.auth.AttemptResponse")
	proto.RegisterType((*Session)(nil), "taask.server.auth.Session")
	proto.RegisterType((*MemberGroup)(nil), "taask.server.auth.MemberGroup")
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874) }

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 402 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x53, 0xdf, 0xeb, 0xd3, 0x30,
	0x10, 0xa7, 0xdb, 0xdc, 0x5c, 0x36, 0x90, 0x45, 0xc4, 0x32, 0x44, 0x4a, 0xf1, 0x61, 0x2f, 0xb6,
	0x32, 0x5f, 0xf4, 0x71, 0x6e, 0x43, 0xa7, 0x4c, 0x24, 0x75, 0x2f, 0xbe, 0xa5, 0xf5, 0x68, 0xc3,
	0x9a, 0x26, 0x24, 0xa9, 0x38, 0xff, 0x0c, 0xff, 0x48, 0xff, 0x0e, 0x69, 0xd2, 0xfd, 0xf2, 0x07,
	0x0c, 0xbe, 0x2f, 0x21, 0xb9, 0xbb, 0xcf, 0x7d, 0x3e, 0x77, 0x97, 0x43, 0x88, 0xd6, 0xa6, 0x88,
	0xa4, 0x12, 0x46, 0xe0, 0x89, 0xa1, 0x54, 0xef, 0x23, 0x0d, 0xea, 0x1b, 0xa8, 0xa8, 0x71, 0x4c,
	0x5f, 0xe4, 0xcc, 0x14, 0x75, 0x1a, 0x65, 0x82, 0xc7, 0x99, 0x28, 0xd8, 0xf7, 0x58, 0x33, 0x2e,
	0xcb, 0x4c, 0x1d, 0xa4, 0x11, 0xb1, 0x85, 0xc5, 0x1c, 0xb4, 0xa6, 0x39, 0xb8, 0x24, 0x37, 0x21,
	0xf6, 0x70, 0x90, 0x94, 0x29, 0x87, 0x08, 0x7f, 0x79, 0x68, 0xb0, 0x30, 0x06, 0xb8, 0x34, 0xf8,
	0x29, 0x42, 0x5b, 0xe0, 0x29, 0xa8, 0xdd, 0x6e, 0xb3, 0xf2, 0xbd, 0xc0, 0x9b, 0x0d, 0xc9, 0x85,
	0x05, 0x3f, 0x41, 0xc3, 0xb7, 0x4a, 0xd4, 0xd2, 0xba, 0x3b, 0xd6, 0x7d, 0x36, 0xe0, 0x57, 0xa8,
	0xff, 0xa9, 0x4e, 0x3f, 0xc0, 0xc1, 0xef, 0x06, 0xde, 0x6c, 0x34, 0x0f, 0x22, 0xcb, 0x1d, 0x39,
	0xf2, 0x28, 0x01, 0xc5, 0x68, 0xc9, 0x7e, 0xd0, 0xb4, 0x04, 0x17, 0x47, 0xda, 0x78, 0xbc, 0x46,
	0x93, 0x45, 0x6d, 0x8a, 0x77, 0x54, 0x17, 0x09, 0xcb, 0x2b, 0x6a, 0x6a, 0x05, 0x7e, 0xcf, 0x26,
	0x79, 0xfc, 0x47, 0x92, 0xa3, 0x9b, 0xfc, 0x8d, 0x68, 0xe4, 0x7d, 0x66, 0x1c, 0xb4, 0xa1, 0x5c,
	0xfa, 0xf7, 0x02, 0x6f, 0xd6, 0x25, 0x67, 0x43, 0xf8, 0xd3, 0x43, 0x0f, 0xda, 0x42, 0x09, 0x68,
	0x29, 0x2a, 0x0d, 0xf8, 0x35, 0x1a, 0xaf, 0xab, 0x6c, 0x59, 0xd0, 0xb2, 0x84, 0x2a, 0x07, 0x5b,
	0xf2, 0x68, 0xfe, 0xe8, 0x9a, 0x73, 0xeb, 0x3a, 0x4c, 0xae, 0x42, 0xf1, 0x0a, 0x8d, 0xb7, 0x54,
	0x1b, 0x50, 0x6d, 0xcd, 0x9d, 0x1b, 0x6b, 0xbe, 0x42, 0x35, 0xa2, 0x06, 0x09, 0x68, 0xcd, 0x44,
	0x75, 0xc7, 0xee, 0x6f, 0xd0, 0xc3, 0x36, 0xd1, 0x49, 0x63, 0xc2, 0xf2, 0x76, 0x14, 0xff, 0xed,
	0xe2, 0xbf, 0x30, 0x21, 0x47, 0x23, 0x47, 0x6b, 0xb3, 0x63, 0x8c, 0x7a, 0x17, 0x8a, 0xec, 0xbd,
	0xb1, 0x7d, 0xa4, 0x1c, 0x5a, 0x19, 0xf6, 0x8e, 0xa7, 0xe8, 0xfe, 0x7b, 0xc1, 0xaa, 0xa5, 0xf8,
	0x0a, 0x96, 0x76, 0x48, 0x4e, 0xef, 0xc6, 0x77, 0x9c, 0x97, 0x1d, 0xec, 0x98, 0x9c, 0xde, 0x6f,
	0x9e, 0x7d, 0x09, 0x2f, 0x7e, 0xad, 0xdd, 0x02, 0x77, 0x3e, 0x77, 0xbb, 0x10, 0x37, 0xbb, 0x90,
	0xf6, 0xed, 0x77, 0x7d, 0xf9, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x0f, 0xfa, 0x06, 0x33, 0x03,
	0x00, 0x00,
}
