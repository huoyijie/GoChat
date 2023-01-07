// protoc --go_out=paths=source_relative:. *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.4
// source: packet.proto

package lib

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PackKind int32

const (
	PackKind_PING   PackKind = 0
	PackKind_PONG   PackKind = 1
	PackKind_SIGNUP PackKind = 2
	PackKind_SIGNIN PackKind = 3
	PackKind_TOKEN  PackKind = 4
	PackKind_USERS  PackKind = 5
	PackKind_MSG    PackKind = 6
	PackKind_ERR    PackKind = 7
)

// Enum value maps for PackKind.
var (
	PackKind_name = map[int32]string{
		0: "PING",
		1: "PONG",
		2: "SIGNUP",
		3: "SIGNIN",
		4: "TOKEN",
		5: "USERS",
		6: "MSG",
		7: "ERR",
	}
	PackKind_value = map[string]int32{
		"PING":   0,
		"PONG":   1,
		"SIGNUP": 2,
		"SIGNIN": 3,
		"TOKEN":  4,
		"USERS":  5,
		"MSG":    6,
		"ERR":    7,
	}
)

func (x PackKind) Enum() *PackKind {
	p := new(PackKind)
	*p = x
	return p
}

func (x PackKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PackKind) Descriptor() protoreflect.EnumDescriptor {
	return file_packet_proto_enumTypes[0].Descriptor()
}

func (PackKind) Type() protoreflect.EnumType {
	return &file_packet_proto_enumTypes[0]
}

func (x PackKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PackKind.Descriptor instead.
func (PackKind) EnumDescriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{0}
}

type MsgKind int32

const (
	MsgKind_TEXT MsgKind = 0
)

// Enum value maps for MsgKind.
var (
	MsgKind_name = map[int32]string{
		0: "TEXT",
	}
	MsgKind_value = map[string]int32{
		"TEXT": 0,
	}
)

func (x MsgKind) Enum() *MsgKind {
	p := new(MsgKind)
	*p = x
	return p
}

func (x MsgKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgKind) Descriptor() protoreflect.EnumDescriptor {
	return file_packet_proto_enumTypes[1].Descriptor()
}

func (MsgKind) Type() protoreflect.EnumType {
	return &file_packet_proto_enumTypes[1]
}

func (x MsgKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgKind.Descriptor instead.
func (MsgKind) EnumDescriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{1}
}

type Packet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Kind PackKind `protobuf:"varint,2,opt,name=kind,proto3,enum=lib.PackKind" json:"kind,omitempty"`
	Data []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Packet) Reset() {
	*x = Packet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Packet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Packet) ProtoMessage() {}

func (x *Packet) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Packet.ProtoReflect.Descriptor instead.
func (*Packet) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{0}
}

func (x *Packet) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Packet) GetKind() PackKind {
	if x != nil {
		return x.Kind
	}
	return PackKind_PING
}

func (x *Packet) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Auth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Passhash []byte `protobuf:"bytes,2,opt,name=passhash,proto3" json:"passhash,omitempty"`
}

func (x *Auth) Reset() {
	*x = Auth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Auth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Auth) ProtoMessage() {}

func (x *Auth) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Auth.ProtoReflect.Descriptor instead.
func (*Auth) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{1}
}

func (x *Auth) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Auth) GetPasshash() []byte {
	if x != nil {
		return x.Passhash
	}
	return nil
}

type Signup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Auth *Auth `protobuf:"bytes,1,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *Signup) Reset() {
	*x = Signup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Signup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Signup) ProtoMessage() {}

func (x *Signup) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Signup.ProtoReflect.Descriptor instead.
func (*Signup) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{2}
}

func (x *Signup) GetAuth() *Auth {
	if x != nil {
		return x.Auth
	}
	return nil
}

type Signin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Auth *Auth `protobuf:"bytes,1,opt,name=auth,proto3" json:"auth,omitempty"`
}

func (x *Signin) Reset() {
	*x = Signin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Signin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Signin) ProtoMessage() {}

func (x *Signin) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Signin.ProtoReflect.Descriptor instead.
func (*Signin) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{3}
}

func (x *Signin) GetAuth() *Auth {
	if x != nil {
		return x.Auth
	}
	return nil
}

type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token []byte `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{4}
}

func (x *Token) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

type TokenRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Id       uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Token    []byte `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *TokenRes) Reset() {
	*x = TokenRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenRes) ProtoMessage() {}

func (x *TokenRes) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenRes.ProtoReflect.Descriptor instead.
func (*TokenRes) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{5}
}

func (x *TokenRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *TokenRes) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *TokenRes) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *TokenRes) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

type UsersRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code  int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Users []string `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *UsersRes) Reset() {
	*x = UsersRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsersRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsersRes) ProtoMessage() {}

func (x *UsersRes) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsersRes.ProtoReflect.Descriptor instead.
func (*UsersRes) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{6}
}

func (x *UsersRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *UsersRes) GetUsers() []string {
	if x != nil {
		return x.Users
	}
	return nil
}

type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind MsgKind `protobuf:"varint,1,opt,name=kind,proto3,enum=lib.MsgKind" json:"kind,omitempty"`
	From string  `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To   string  `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Data []byte  `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{7}
}

func (x *Msg) GetKind() MsgKind {
	if x != nil {
		return x.Kind
	}
	return MsgKind_TEXT
}

func (x *Msg) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Msg) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Msg) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type ErrRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *ErrRes) Reset() {
	*x = ErrRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packet_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrRes) ProtoMessage() {}

func (x *ErrRes) ProtoReflect() protoreflect.Message {
	mi := &file_packet_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrRes.ProtoReflect.Descriptor instead.
func (*ErrRes) Descriptor() ([]byte, []int) {
	return file_packet_proto_rawDescGZIP(), []int{8}
}

func (x *ErrRes) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

var File_packet_proto protoreflect.FileDescriptor

var file_packet_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x6c, 0x69, 0x62, 0x22, 0x4f, 0x0a, 0x06, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x21, 0x0a,
	0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x6c, 0x69,
	0x62, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x3e, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x68, 0x61, 0x73, 0x68, 0x22, 0x27, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x1d,
	0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x6c,
	0x69, 0x62, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x27, 0x0a,
	0x06, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x12, 0x1d, 0x0a, 0x04, 0x61, 0x75, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x6c, 0x69, 0x62, 0x2e, 0x41, 0x75, 0x74, 0x68,
	0x52, 0x04, 0x61, 0x75, 0x74, 0x68, 0x22, 0x1d, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x60, 0x0a, 0x08, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x34, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x5f, 0x0a,
	0x03, 0x4d, 0x73, 0x67, 0x12, 0x20, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x6c, 0x69, 0x62, 0x2e, 0x4d, 0x73, 0x67, 0x4b, 0x69, 0x6e, 0x64,
	0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x1c,
	0x0a, 0x06, 0x45, 0x72, 0x72, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x2a, 0x5e, 0x0a, 0x08,
	0x50, 0x61, 0x63, 0x6b, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x49, 0x4e, 0x47,
	0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x4f, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06,
	0x53, 0x49, 0x47, 0x4e, 0x55, 0x50, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x49, 0x47, 0x4e,
	0x49, 0x4e, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x10, 0x04, 0x12,
	0x09, 0x0a, 0x05, 0x55, 0x53, 0x45, 0x52, 0x53, 0x10, 0x05, 0x12, 0x07, 0x0a, 0x03, 0x4d, 0x53,
	0x47, 0x10, 0x06, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x52, 0x52, 0x10, 0x07, 0x2a, 0x13, 0x0a, 0x07,
	0x4d, 0x73, 0x67, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x45, 0x58, 0x54, 0x10,
	0x00, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x68, 0x75, 0x6f, 0x79, 0x69, 0x6a, 0x69, 0x65, 0x2f, 0x47, 0x6f, 0x43, 0x68, 0x61, 0x74, 0x2f,
	0x6c, 0x69, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_packet_proto_rawDescOnce sync.Once
	file_packet_proto_rawDescData = file_packet_proto_rawDesc
)

func file_packet_proto_rawDescGZIP() []byte {
	file_packet_proto_rawDescOnce.Do(func() {
		file_packet_proto_rawDescData = protoimpl.X.CompressGZIP(file_packet_proto_rawDescData)
	})
	return file_packet_proto_rawDescData
}

var file_packet_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_packet_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_packet_proto_goTypes = []interface{}{
	(PackKind)(0),    // 0: lib.PackKind
	(MsgKind)(0),     // 1: lib.MsgKind
	(*Packet)(nil),   // 2: lib.Packet
	(*Auth)(nil),     // 3: lib.Auth
	(*Signup)(nil),   // 4: lib.Signup
	(*Signin)(nil),   // 5: lib.Signin
	(*Token)(nil),    // 6: lib.Token
	(*TokenRes)(nil), // 7: lib.TokenRes
	(*UsersRes)(nil), // 8: lib.UsersRes
	(*Msg)(nil),      // 9: lib.Msg
	(*ErrRes)(nil),   // 10: lib.ErrRes
}
var file_packet_proto_depIdxs = []int32{
	0, // 0: lib.Packet.kind:type_name -> lib.PackKind
	3, // 1: lib.Signup.auth:type_name -> lib.Auth
	3, // 2: lib.Signin.auth:type_name -> lib.Auth
	1, // 3: lib.Msg.kind:type_name -> lib.MsgKind
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_packet_proto_init() }
func file_packet_proto_init() {
	if File_packet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_packet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Packet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Auth); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Signup); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Signin); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Token); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsersRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_packet_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_packet_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_packet_proto_goTypes,
		DependencyIndexes: file_packet_proto_depIdxs,
		EnumInfos:         file_packet_proto_enumTypes,
		MessageInfos:      file_packet_proto_msgTypes,
	}.Build()
	File_packet_proto = out.File
	file_packet_proto_rawDesc = nil
	file_packet_proto_goTypes = nil
	file_packet_proto_depIdxs = nil
}
