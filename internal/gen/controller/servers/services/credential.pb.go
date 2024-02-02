// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: controller/servers/services/v1/credential.proto

package services

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

type Credential struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Credential:
	//
	//	*Credential_UsernamePassword
	//	*Credential_SshPrivateKey
	//	*Credential_SshCertificate
	Credential isCredential_Credential `protobuf_oneof:"credential"`
}

func (x *Credential) Reset() {
	*x = Credential{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_servers_services_v1_credential_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Credential) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Credential) ProtoMessage() {}

func (x *Credential) ProtoReflect() protoreflect.Message {
	mi := &file_controller_servers_services_v1_credential_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Credential.ProtoReflect.Descriptor instead.
func (*Credential) Descriptor() ([]byte, []int) {
	return file_controller_servers_services_v1_credential_proto_rawDescGZIP(), []int{0}
}

func (m *Credential) GetCredential() isCredential_Credential {
	if m != nil {
		return m.Credential
	}
	return nil
}

func (x *Credential) GetUsernamePassword() *UsernamePassword {
	if x, ok := x.GetCredential().(*Credential_UsernamePassword); ok {
		return x.UsernamePassword
	}
	return nil
}

func (x *Credential) GetSshPrivateKey() *SshPrivateKey {
	if x, ok := x.GetCredential().(*Credential_SshPrivateKey); ok {
		return x.SshPrivateKey
	}
	return nil
}

func (x *Credential) GetSshCertificate() *SshCertificate {
	if x, ok := x.GetCredential().(*Credential_SshCertificate); ok {
		return x.SshCertificate
	}
	return nil
}

type isCredential_Credential interface {
	isCredential_Credential()
}

type Credential_UsernamePassword struct {
	UsernamePassword *UsernamePassword `protobuf:"bytes,2,opt,name=username_password,json=usernamePassword,proto3,oneof"`
}

type Credential_SshPrivateKey struct {
	SshPrivateKey *SshPrivateKey `protobuf:"bytes,3,opt,name=ssh_private_key,json=sshPrivateKey,proto3,oneof"`
}

type Credential_SshCertificate struct {
	SshCertificate *SshCertificate `protobuf:"bytes,4,opt,name=ssh_certificate,json=sshCertificate,proto3,oneof"`
}

func (*Credential_UsernamePassword) isCredential_Credential() {}

func (*Credential_SshPrivateKey) isCredential_Credential() {}

func (*Credential_SshCertificate) isCredential_Credential() {}

// UsernamePassword is a credential containing a username and a password.
type UsernamePassword struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The username of the credential
	Username string `protobuf:"bytes,10,opt,name=username,proto3" json:"username,omitempty"` // @gotags: `class:"public"`
	// The password of the credential
	Password string `protobuf:"bytes,20,opt,name=password,proto3" json:"password,omitempty"` // @gotags: `class:"secret"`
}

func (x *UsernamePassword) Reset() {
	*x = UsernamePassword{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_servers_services_v1_credential_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsernamePassword) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsernamePassword) ProtoMessage() {}

func (x *UsernamePassword) ProtoReflect() protoreflect.Message {
	mi := &file_controller_servers_services_v1_credential_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsernamePassword.ProtoReflect.Descriptor instead.
func (*UsernamePassword) Descriptor() ([]byte, []int) {
	return file_controller_servers_services_v1_credential_proto_rawDescGZIP(), []int{1}
}

func (x *UsernamePassword) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UsernamePassword) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// SshPrivateKey is a credential containing a username a private key and an optional
// private key passphrase.
type SshPrivateKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The username of the credential
	Username string `protobuf:"bytes,10,opt,name=username,proto3" json:"username,omitempty"` // @gotags: `class:"public"`
	// The private key of the credential
	PrivateKey string `protobuf:"bytes,20,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"` // @gotags: `class:"secret"`
	// The optional passphrase of the private_key
	PrivateKeyPassphrase string `protobuf:"bytes,30,opt,name=private_key_passphrase,json=privateKeyPassphrase,proto3" json:"private_key_passphrase,omitempty"` // @gotags: `class:"secret"`
}

func (x *SshPrivateKey) Reset() {
	*x = SshPrivateKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_servers_services_v1_credential_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SshPrivateKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SshPrivateKey) ProtoMessage() {}

func (x *SshPrivateKey) ProtoReflect() protoreflect.Message {
	mi := &file_controller_servers_services_v1_credential_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SshPrivateKey.ProtoReflect.Descriptor instead.
func (*SshPrivateKey) Descriptor() ([]byte, []int) {
	return file_controller_servers_services_v1_credential_proto_rawDescGZIP(), []int{2}
}

func (x *SshPrivateKey) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SshPrivateKey) GetPrivateKey() string {
	if x != nil {
		return x.PrivateKey
	}
	return ""
}

func (x *SshPrivateKey) GetPrivateKeyPassphrase() string {
	if x != nil {
		return x.PrivateKeyPassphrase
	}
	return ""
}

// SshCertificate is a credential containing a username, private key, and
// client certificate.
type SshCertificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The username of the credential
	Username string `protobuf:"bytes,10,opt,name=username,proto3" json:"username,omitempty"` // @gotags: `class:"public"`
	// The private key of the credential
	PrivateKey string `protobuf:"bytes,20,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"` // @gotags: `class:"secret"`
	// The client certificate signed by a CA to establish trust of the private key.
	Certificate string `protobuf:"bytes,30,opt,name=certificate,proto3" json:"certificate,omitempty"` // @gotags: `class:"public"`
}

func (x *SshCertificate) Reset() {
	*x = SshCertificate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_servers_services_v1_credential_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SshCertificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SshCertificate) ProtoMessage() {}

func (x *SshCertificate) ProtoReflect() protoreflect.Message {
	mi := &file_controller_servers_services_v1_credential_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SshCertificate.ProtoReflect.Descriptor instead.
func (*SshCertificate) Descriptor() ([]byte, []int) {
	return file_controller_servers_services_v1_credential_proto_rawDescGZIP(), []int{3}
}

func (x *SshCertificate) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *SshCertificate) GetPrivateKey() string {
	if x != nil {
		return x.PrivateKey
	}
	return ""
}

func (x *SshCertificate) GetCertificate() string {
	if x != nil {
		return x.Certificate
	}
	return ""
}

var File_controller_servers_services_v1_credential_proto protoreflect.FileDescriptor

var file_controller_servers_services_v1_credential_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x73, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76,
	0x31, 0x22, 0xc4, 0x02, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x12, 0x5f, 0x0a, 0x11, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x48, 0x00, 0x52,
	0x10, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x12, 0x57, 0x0a, 0x0f, 0x73, 0x73, 0x68, 0x5f, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x73, 0x68, 0x50,
	0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x48, 0x00, 0x52, 0x0d, 0x73, 0x73, 0x68,
	0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x59, 0x0a, 0x0f, 0x73, 0x73,
	0x68, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x73, 0x68, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x0e, 0x73, 0x73, 0x68, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x61, 0x6c, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x52, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x4a, 0x0a, 0x10, 0x55, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x22, 0x82, 0x01, 0x0a, 0x0d, 0x53, 0x73, 0x68, 0x50, 0x72, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x4b, 0x65, 0x79, 0x12, 0x34, 0x0a, 0x16, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x6b,
	0x65, 0x79, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x70, 0x68, 0x72, 0x61, 0x73, 0x65, 0x18, 0x1e, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x14, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x50,
	0x61, 0x73, 0x73, 0x70, 0x68, 0x72, 0x61, 0x73, 0x65, 0x22, 0x6f, 0x0a, 0x0e, 0x53, 0x73, 0x68,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x69, 0x76, 0x61,
	0x74, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x42, 0x51, 0x5a, 0x4f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f,
	0x72, 0x70, 0x2f, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x61, 0x72, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x6c, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x3b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_controller_servers_services_v1_credential_proto_rawDescOnce sync.Once
	file_controller_servers_services_v1_credential_proto_rawDescData = file_controller_servers_services_v1_credential_proto_rawDesc
)

func file_controller_servers_services_v1_credential_proto_rawDescGZIP() []byte {
	file_controller_servers_services_v1_credential_proto_rawDescOnce.Do(func() {
		file_controller_servers_services_v1_credential_proto_rawDescData = protoimpl.X.CompressGZIP(file_controller_servers_services_v1_credential_proto_rawDescData)
	})
	return file_controller_servers_services_v1_credential_proto_rawDescData
}

var file_controller_servers_services_v1_credential_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_controller_servers_services_v1_credential_proto_goTypes = []interface{}{
	(*Credential)(nil),       // 0: controller.servers.services.v1.Credential
	(*UsernamePassword)(nil), // 1: controller.servers.services.v1.UsernamePassword
	(*SshPrivateKey)(nil),    // 2: controller.servers.services.v1.SshPrivateKey
	(*SshCertificate)(nil),   // 3: controller.servers.services.v1.SshCertificate
}
var file_controller_servers_services_v1_credential_proto_depIdxs = []int32{
	1, // 0: controller.servers.services.v1.Credential.username_password:type_name -> controller.servers.services.v1.UsernamePassword
	2, // 1: controller.servers.services.v1.Credential.ssh_private_key:type_name -> controller.servers.services.v1.SshPrivateKey
	3, // 2: controller.servers.services.v1.Credential.ssh_certificate:type_name -> controller.servers.services.v1.SshCertificate
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_controller_servers_services_v1_credential_proto_init() }
func file_controller_servers_services_v1_credential_proto_init() {
	if File_controller_servers_services_v1_credential_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_controller_servers_services_v1_credential_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Credential); i {
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
		file_controller_servers_services_v1_credential_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsernamePassword); i {
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
		file_controller_servers_services_v1_credential_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SshPrivateKey); i {
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
		file_controller_servers_services_v1_credential_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SshCertificate); i {
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
	file_controller_servers_services_v1_credential_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Credential_UsernamePassword)(nil),
		(*Credential_SshPrivateKey)(nil),
		(*Credential_SshCertificate)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_controller_servers_services_v1_credential_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_controller_servers_services_v1_credential_proto_goTypes,
		DependencyIndexes: file_controller_servers_services_v1_credential_proto_depIdxs,
		MessageInfos:      file_controller_servers_services_v1_credential_proto_msgTypes,
	}.Build()
	File_controller_servers_services_v1_credential_proto = out.File
	file_controller_servers_services_v1_credential_proto_rawDesc = nil
	file_controller_servers_services_v1_credential_proto_goTypes = nil
	file_controller_servers_services_v1_credential_proto_depIdxs = nil
}
