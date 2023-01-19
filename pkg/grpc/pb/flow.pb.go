// Refer: https://github.com/devops-metalflow/metalmetrics/blob/master/src/flow/flow.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.5.0
// source: pkg/grpc/pb/flow.proto

package pb

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

// The request message.
type MetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MetricsRequest) Reset() {
	*x = MetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_pb_flow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsRequest) ProtoMessage() {}

func (x *MetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_pb_flow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsRequest.ProtoReflect.Descriptor instead.
func (*MetricsRequest) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_pb_flow_proto_rawDescGZIP(), []int{0}
}

func (x *MetricsRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// The response message.
type MetricsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error  string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Output string `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *MetricsReply) Reset() {
	*x = MetricsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_grpc_pb_flow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsReply) ProtoMessage() {}

func (x *MetricsReply) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_grpc_pb_flow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsReply.ProtoReflect.Descriptor instead.
func (*MetricsReply) Descriptor() ([]byte, []int) {
	return file_pkg_grpc_pb_flow_proto_rawDescGZIP(), []int{1}
}

func (x *MetricsReply) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *MetricsReply) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

var File_pkg_grpc_pb_flow_proto protoreflect.FileDescriptor

var file_pkg_grpc_pb_flow_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x66, 0x6c,
	0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x66, 0x6c, 0x6f, 0x77, 0x22, 0x2a,
	0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x3c, 0x0a, 0x0c, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x32, 0x49, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x39, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x14, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_grpc_pb_flow_proto_rawDescOnce sync.Once
	file_pkg_grpc_pb_flow_proto_rawDescData = file_pkg_grpc_pb_flow_proto_rawDesc
)

func file_pkg_grpc_pb_flow_proto_rawDescGZIP() []byte {
	file_pkg_grpc_pb_flow_proto_rawDescOnce.Do(func() {
		file_pkg_grpc_pb_flow_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_grpc_pb_flow_proto_rawDescData)
	})
	return file_pkg_grpc_pb_flow_proto_rawDescData
}

var file_pkg_grpc_pb_flow_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_grpc_pb_flow_proto_goTypes = []interface{}{
	(*MetricsRequest)(nil), // 0: flow.MetricsRequest
	(*MetricsReply)(nil),   // 1: flow.MetricsReply
}
var file_pkg_grpc_pb_flow_proto_depIdxs = []int32{
	0, // 0: flow.MetricsProto.SendMetrics:input_type -> flow.MetricsRequest
	1, // 1: flow.MetricsProto.SendMetrics:output_type -> flow.MetricsReply
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_grpc_pb_flow_proto_init() }
func file_pkg_grpc_pb_flow_proto_init() {
	if File_pkg_grpc_pb_flow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_grpc_pb_flow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricsRequest); i {
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
		file_pkg_grpc_pb_flow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricsReply); i {
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
			RawDescriptor: file_pkg_grpc_pb_flow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_grpc_pb_flow_proto_goTypes,
		DependencyIndexes: file_pkg_grpc_pb_flow_proto_depIdxs,
		MessageInfos:      file_pkg_grpc_pb_flow_proto_msgTypes,
	}.Build()
	File_pkg_grpc_pb_flow_proto = out.File
	file_pkg_grpc_pb_flow_proto_rawDesc = nil
	file_pkg_grpc_pb_flow_proto_goTypes = nil
	file_pkg_grpc_pb_flow_proto_depIdxs = nil
}