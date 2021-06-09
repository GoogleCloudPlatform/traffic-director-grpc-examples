// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0
// source: proto/grpc/examples/wallet/account/account.proto

package account

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

type MembershipType int32

const (
	MembershipType_UNKNOWN_MEMBERSHIP_TYPE MembershipType = 0
	MembershipType_NORMAL                  MembershipType = 1
	MembershipType_PREMIUM                 MembershipType = 2
)

// Enum value maps for MembershipType.
var (
	MembershipType_name = map[int32]string{
		0: "UNKNOWN_MEMBERSHIP_TYPE",
		1: "NORMAL",
		2: "PREMIUM",
	}
	MembershipType_value = map[string]int32{
		"UNKNOWN_MEMBERSHIP_TYPE": 0,
		"NORMAL":                  1,
		"PREMIUM":                 2,
	}
)

func (x MembershipType) Enum() *MembershipType {
	p := new(MembershipType)
	*p = x
	return p
}

func (x MembershipType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MembershipType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_grpc_examples_wallet_account_account_proto_enumTypes[0].Descriptor()
}

func (MembershipType) Type() protoreflect.EnumType {
	return &file_proto_grpc_examples_wallet_account_account_proto_enumTypes[0]
}

func (x MembershipType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MembershipType.Descriptor instead.
func (MembershipType) EnumDescriptor() ([]byte, []int) {
	return file_proto_grpc_examples_wallet_account_account_proto_rawDescGZIP(), []int{0}
}

type GetUserInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *GetUserInfoRequest) Reset() {
	*x = GetUserInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_grpc_examples_wallet_account_account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoRequest) ProtoMessage() {}

func (x *GetUserInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_grpc_examples_wallet_account_account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoRequest.ProtoReflect.Descriptor instead.
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) {
	return file_proto_grpc_examples_wallet_account_account_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserInfoRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type GetUserInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Membership MembershipType `protobuf:"varint,2,opt,name=membership,proto3,enum=grpc.examples.wallet.account.MembershipType" json:"membership,omitempty"`
}

func (x *GetUserInfoResponse) Reset() {
	*x = GetUserInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_grpc_examples_wallet_account_account_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInfoResponse) ProtoMessage() {}

func (x *GetUserInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_grpc_examples_wallet_account_account_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInfoResponse.ProtoReflect.Descriptor instead.
func (*GetUserInfoResponse) Descriptor() ([]byte, []int) {
	return file_proto_grpc_examples_wallet_account_account_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserInfoResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetUserInfoResponse) GetMembership() MembershipType {
	if x != nil {
		return x.Membership
	}
	return MembershipType_UNKNOWN_MEMBERSHIP_TYPE
}

var File_proto_grpc_examples_wallet_account_account_proto protoreflect.FileDescriptor

var file_proto_grpc_examples_wallet_account_account_proto_rawDesc = []byte{
	0x0a, 0x30, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1c, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x73, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0x2a, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x77, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x4c, 0x0a, 0x0a, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x68, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2c, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x68, 0x69, 0x70, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x73, 0x68, 0x69, 0x70, 0x2a, 0x46, 0x0a, 0x0e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73,
	0x68, 0x69, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x55, 0x4e, 0x4b, 0x4e, 0x4f,
	0x57, 0x4e, 0x5f, 0x4d, 0x45, 0x4d, 0x42, 0x45, 0x52, 0x53, 0x48, 0x49, 0x50, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x4f, 0x52, 0x4d, 0x41, 0x4c, 0x10, 0x01,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x52, 0x45, 0x4d, 0x49, 0x55, 0x4d, 0x10, 0x02, 0x32, 0x7f, 0x0a,
	0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x74, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x30, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x72,
	0x0a, 0x1f, 0x69, 0x6f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x73, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x42, 0x0c, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x3f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67,
	0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x73, 0x2f, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_grpc_examples_wallet_account_account_proto_rawDescOnce sync.Once
	file_proto_grpc_examples_wallet_account_account_proto_rawDescData = file_proto_grpc_examples_wallet_account_account_proto_rawDesc
)

func file_proto_grpc_examples_wallet_account_account_proto_rawDescGZIP() []byte {
	file_proto_grpc_examples_wallet_account_account_proto_rawDescOnce.Do(func() {
		file_proto_grpc_examples_wallet_account_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_grpc_examples_wallet_account_account_proto_rawDescData)
	})
	return file_proto_grpc_examples_wallet_account_account_proto_rawDescData
}

var file_proto_grpc_examples_wallet_account_account_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_grpc_examples_wallet_account_account_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_grpc_examples_wallet_account_account_proto_goTypes = []interface{}{
	(MembershipType)(0),         // 0: grpc.examples.wallet.account.MembershipType
	(*GetUserInfoRequest)(nil),  // 1: grpc.examples.wallet.account.GetUserInfoRequest
	(*GetUserInfoResponse)(nil), // 2: grpc.examples.wallet.account.GetUserInfoResponse
}
var file_proto_grpc_examples_wallet_account_account_proto_depIdxs = []int32{
	0, // 0: grpc.examples.wallet.account.GetUserInfoResponse.membership:type_name -> grpc.examples.wallet.account.MembershipType
	1, // 1: grpc.examples.wallet.account.Account.GetUserInfo:input_type -> grpc.examples.wallet.account.GetUserInfoRequest
	2, // 2: grpc.examples.wallet.account.Account.GetUserInfo:output_type -> grpc.examples.wallet.account.GetUserInfoResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_grpc_examples_wallet_account_account_proto_init() }
func file_proto_grpc_examples_wallet_account_account_proto_init() {
	if File_proto_grpc_examples_wallet_account_account_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_grpc_examples_wallet_account_account_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserInfoRequest); i {
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
		file_proto_grpc_examples_wallet_account_account_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserInfoResponse); i {
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
			RawDescriptor: file_proto_grpc_examples_wallet_account_account_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_grpc_examples_wallet_account_account_proto_goTypes,
		DependencyIndexes: file_proto_grpc_examples_wallet_account_account_proto_depIdxs,
		EnumInfos:         file_proto_grpc_examples_wallet_account_account_proto_enumTypes,
		MessageInfos:      file_proto_grpc_examples_wallet_account_account_proto_msgTypes,
	}.Build()
	File_proto_grpc_examples_wallet_account_account_proto = out.File
	file_proto_grpc_examples_wallet_account_account_proto_rawDesc = nil
	file_proto_grpc_examples_wallet_account_account_proto_goTypes = nil
	file_proto_grpc_examples_wallet_account_account_proto_depIdxs = nil
}
