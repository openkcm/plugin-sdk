// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: plugin/test/v1/test.proto

package testv1

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

type TestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Request string `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
}

func (x *TestRequest) Reset() {
	*x = TestRequest{}
	mi := &file_plugin_test_v1_test_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRequest) ProtoMessage() {}

func (x *TestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_test_v1_test_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRequest.ProtoReflect.Descriptor instead.
func (*TestRequest) Descriptor() ([]byte, []int) {
	return file_plugin_test_v1_test_proto_rawDescGZIP(), []int{0}
}

func (x *TestRequest) GetRequest() string {
	if x != nil {
		return x.Request
	}
	return ""
}

type TestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *TestResponse) Reset() {
	*x = TestResponse{}
	mi := &file_plugin_test_v1_test_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResponse) ProtoMessage() {}

func (x *TestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_test_v1_test_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResponse.ProtoReflect.Descriptor instead.
func (*TestResponse) Descriptor() ([]byte, []int) {
	return file_plugin_test_v1_test_proto_rawDescGZIP(), []int{1}
}

func (x *TestResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

var File_plugin_test_v1_test_proto protoreflect.FileDescriptor

var file_plugin_test_v1_test_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x76, 0x31,
	0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x22, 0x27, 0x0a, 0x0b, 0x54,
	0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x2a, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0x50, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x41, 0x0a, 0x04, 0x54, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x74, 0x65,
	0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0xb4, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x54, 0x65, 0x73, 0x74, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x6b, 0x63, 0x6d, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x74, 0x65, 0x73, 0x74, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x50, 0x54, 0x58, 0xaa, 0x02, 0x0e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e, 0x50, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x5c, 0x54, 0x65, 0x73, 0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1a, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x5c, 0x54, 0x65, 0x73, 0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x10, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x3a,
	0x3a, 0x54, 0x65, 0x73, 0x74, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_plugin_test_v1_test_proto_rawDescOnce sync.Once
	file_plugin_test_v1_test_proto_rawDescData = file_plugin_test_v1_test_proto_rawDesc
)

func file_plugin_test_v1_test_proto_rawDescGZIP() []byte {
	file_plugin_test_v1_test_proto_rawDescOnce.Do(func() {
		file_plugin_test_v1_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_test_v1_test_proto_rawDescData)
	})
	return file_plugin_test_v1_test_proto_rawDescData
}

var file_plugin_test_v1_test_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_plugin_test_v1_test_proto_goTypes = []any{
	(*TestRequest)(nil),  // 0: plugin.test.v1.TestRequest
	(*TestResponse)(nil), // 1: plugin.test.v1.TestResponse
}
var file_plugin_test_v1_test_proto_depIdxs = []int32{
	0, // 0: plugin.test.v1.TestService.Test:input_type -> plugin.test.v1.TestRequest
	1, // 1: plugin.test.v1.TestService.Test:output_type -> plugin.test.v1.TestResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_plugin_test_v1_test_proto_init() }
func file_plugin_test_v1_test_proto_init() {
	if File_plugin_test_v1_test_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_plugin_test_v1_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_test_v1_test_proto_goTypes,
		DependencyIndexes: file_plugin_test_v1_test_proto_depIdxs,
		MessageInfos:      file_plugin_test_v1_test_proto_msgTypes,
	}.Build()
	File_plugin_test_v1_test_proto = out.File
	file_plugin_test_v1_test_proto_rawDesc = nil
	file_plugin_test_v1_test_proto_goTypes = nil
	file_plugin_test_v1_test_proto_depIdxs = nil
}
