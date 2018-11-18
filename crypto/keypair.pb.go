// Code generated by protoc-gen-go. DO NOT EDIT.
// source: keypair.proto

package crypto

import (
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

type SerializablePubKey struct {
	N                    string   `protobuf:"bytes,1,opt,name=N,proto3" json:"N,omitempty"`
	E                    int64    `protobuf:"varint,2,opt,name=E,proto3" json:"E,omitempty"`
	KID                  string   `protobuf:"bytes,3,opt,name=KID,proto3" json:"KID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SerializablePubKey) Reset()         { *m = SerializablePubKey{} }
func (m *SerializablePubKey) String() string { return proto.CompactTextString(m) }
func (*SerializablePubKey) ProtoMessage()    {}
func (*SerializablePubKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_40535d21a2c7da09, []int{0}
}

func (m *SerializablePubKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SerializablePubKey.Unmarshal(m, b)
}
func (m *SerializablePubKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SerializablePubKey.Marshal(b, m, deterministic)
}
func (m *SerializablePubKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SerializablePubKey.Merge(m, src)
}
func (m *SerializablePubKey) XXX_Size() int {
	return xxx_messageInfo_SerializablePubKey.Size(m)
}
func (m *SerializablePubKey) XXX_DiscardUnknown() {
	xxx_messageInfo_SerializablePubKey.DiscardUnknown(m)
}

var xxx_messageInfo_SerializablePubKey proto.InternalMessageInfo

func (m *SerializablePubKey) GetN() string {
	if m != nil {
		return m.N
	}
	return ""
}

func (m *SerializablePubKey) GetE() int64 {
	if m != nil {
		return m.E
	}
	return 0
}

func (m *SerializablePubKey) GetKID() string {
	if m != nil {
		return m.KID
	}
	return ""
}

func init() {
	proto.RegisterType((*SerializablePubKey)(nil), "taask.server.crypto.serializablePubKey")
}

func init() { proto.RegisterFile("keypair.proto", fileDescriptor_40535d21a2c7da09) }

var fileDescriptor_40535d21a2c7da09 = []byte{
	// 132 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4e, 0xad, 0x2c,
	0x48, 0xcc, 0x2c, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x2e, 0x49, 0x4c, 0x2c, 0xce,
	0xd6, 0x2b, 0x4e, 0x2d, 0x2a, 0x4b, 0x2d, 0xd2, 0x4b, 0x2e, 0xaa, 0x2c, 0x28, 0xc9, 0x57, 0x72,
	0xe2, 0x12, 0x2a, 0x4e, 0x2d, 0xca, 0x4c, 0xcc, 0xc9, 0xac, 0x4a, 0x4c, 0xca, 0x49, 0x0d, 0x28,
	0x4d, 0xf2, 0x4e, 0xad, 0x14, 0xe2, 0xe1, 0x62, 0xf4, 0x93, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c,
	0x62, 0xf4, 0x03, 0xf1, 0x5c, 0x25, 0x98, 0x14, 0x18, 0x35, 0x98, 0x83, 0x18, 0x5d, 0x85, 0x04,
	0xb8, 0x98, 0xbd, 0x3d, 0x5d, 0x24, 0x98, 0xc1, 0xb2, 0x20, 0xa6, 0x13, 0x47, 0x14, 0x1b, 0xc4,
	0xb4, 0x24, 0x36, 0xb0, 0x4d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa9, 0x30, 0x5c, 0x16,
	0x7a, 0x00, 0x00, 0x00,
}
