// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: config/config.proto

package config

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Recording struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// recordRequestBody enables/disables the recording of the request body
	RecordRequestBody *bool `protobuf:"varint,1,opt,name=recordRequestBody,def=1" json:"recordRequestBody,omitempty"`
	// recordRequestBody enables/disables the recording of the request headers
	RecordRequestHeaders *bool `protobuf:"varint,2,opt,name=recordRequestHeaders,def=1" json:"recordRequestHeaders,omitempty"`
	// recordRequestBody enables/disables the recording of the response body
	RecordResponseBody *bool `protobuf:"varint,3,opt,name=recordResponseBody,def=1" json:"recordResponseBody,omitempty"`
	// recordRequestBody enables/disables the recording of the response headers
	RecordResponseHeaders *bool `protobuf:"varint,4,opt,name=recordResponseHeaders,def=1" json:"recordResponseHeaders,omitempty"`
}

// Default values for Recording fields.
const (
	Default_Recording_RecordRequestBody     = bool(true)
	Default_Recording_RecordRequestHeaders  = bool(true)
	Default_Recording_RecordResponseBody    = bool(true)
	Default_Recording_RecordResponseHeaders = bool(true)
)

func (x *Recording) Reset() {
	*x = Recording{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Recording) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Recording) ProtoMessage() {}

func (x *Recording) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Recording.ProtoReflect.Descriptor instead.
func (*Recording) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0}
}

func (x *Recording) GetRecordRequestBody() bool {
	if x != nil && x.RecordRequestBody != nil {
		return *x.RecordRequestBody
	}
	return Default_Recording_RecordRequestBody
}

func (x *Recording) GetRecordRequestHeaders() bool {
	if x != nil && x.RecordRequestHeaders != nil {
		return *x.RecordRequestHeaders
	}
	return Default_Recording_RecordRequestHeaders
}

func (x *Recording) GetRecordResponseBody() bool {
	if x != nil && x.RecordResponseBody != nil {
		return *x.RecordResponseBody
	}
	return Default_Recording_RecordResponseBody
}

func (x *Recording) GetRecordResponseHeaders() bool {
	if x != nil && x.RecordResponseHeaders != nil {
		return *x.RecordResponseHeaders
	}
	return Default_Recording_RecordResponseHeaders
}

type Reporting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host   *string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	ApiKey *string `protobuf:"bytes,2,opt,name=apiKey" json:"apiKey,omitempty"`
}

func (x *Reporting) Reset() {
	*x = Reporting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reporting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reporting) ProtoMessage() {}

func (x *Reporting) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reporting.ProtoReflect.Descriptor instead.
func (*Reporting) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{1}
}

func (x *Reporting) GetHost() string {
	if x != nil && x.Host != nil {
		return *x.Host
	}
	return ""
}

func (x *Reporting) GetApiKey() string {
	if x != nil && x.ApiKey != nil {
		return *x.ApiKey
	}
	return ""
}

type Instrumentation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceName *string `protobuf:"bytes,1,req,name=serviceName" json:"serviceName,omitempty"`
}

func (x *Instrumentation) Reset() {
	*x = Instrumentation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Instrumentation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Instrumentation) ProtoMessage() {}

func (x *Instrumentation) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Instrumentation.ProtoReflect.Descriptor instead.
func (*Instrumentation) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{2}
}

func (x *Instrumentation) GetServiceName() string {
	if x != nil && x.ServiceName != nil {
		return *x.ServiceName
	}
	return ""
}

type AgentConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instrumentation *Instrumentation `protobuf:"bytes,1,req,name=instrumentation" json:"instrumentation,omitempty"`
	Recording       *Recording       `protobuf:"bytes,2,req,name=recording" json:"recording,omitempty"`
	Reporting       *Reporting       `protobuf:"bytes,3,req,name=reporting" json:"reporting,omitempty"`
}

func (x *AgentConfig) Reset() {
	*x = AgentConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentConfig) ProtoMessage() {}

func (x *AgentConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentConfig.ProtoReflect.Descriptor instead.
func (*AgentConfig) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{3}
}

func (x *AgentConfig) GetInstrumentation() *Instrumentation {
	if x != nil {
		return x.Instrumentation
	}
	return nil
}

func (x *AgentConfig) GetRecording() *Recording {
	if x != nil {
		return x.Recording
	}
	return nil
}

func (x *AgentConfig) GetReporting() *Reporting {
	if x != nil {
		return x.Reporting
	}
	return nil
}

var File_config_config_proto protoreflect.FileDescriptor

var file_config_config_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xeb, 0x01, 0x0a, 0x09, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x69, 0x6e, 0x67, 0x12, 0x32, 0x0a, 0x11, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x04,
	0x74, 0x72, 0x75, 0x65, 0x52, 0x11, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x38, 0x0a, 0x14, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x04, 0x74, 0x72, 0x75, 0x65, 0x52, 0x14, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x12, 0x34, 0x0a, 0x12, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x04, 0x74,
	0x72, 0x75, 0x65, 0x52, 0x12, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x3a, 0x0a, 0x15, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x04, 0x74, 0x72, 0x75, 0x65, 0x52, 0x15, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x22, 0x37, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x68, 0x6f, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x22, 0x33, 0x0a, 0x0f,
	0x49, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x20, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x02, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0x9d, 0x01, 0x0a, 0x0b, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x3a, 0x0a, 0x0f, 0x69, 0x6e, 0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x02, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x69, 0x6e,
	0x73, 0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a,
	0x09, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x02, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x09, 0x72, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x28, 0x0a, 0x09, 0x72, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x02, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x09, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e,
	0x67, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x69, 0x2f, 0x67, 0x6f, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
}

var (
	file_config_config_proto_rawDescOnce sync.Once
	file_config_config_proto_rawDescData = file_config_config_proto_rawDesc
)

func file_config_config_proto_rawDescGZIP() []byte {
	file_config_config_proto_rawDescOnce.Do(func() {
		file_config_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_config_proto_rawDescData)
	})
	return file_config_config_proto_rawDescData
}

var file_config_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_config_config_proto_goTypes = []interface{}{
	(*Recording)(nil),       // 0: Recording
	(*Reporting)(nil),       // 1: Reporting
	(*Instrumentation)(nil), // 2: Instrumentation
	(*AgentConfig)(nil),     // 3: AgentConfig
}
var file_config_config_proto_depIdxs = []int32{
	2, // 0: AgentConfig.instrumentation:type_name -> Instrumentation
	0, // 1: AgentConfig.recording:type_name -> Recording
	1, // 2: AgentConfig.reporting:type_name -> Reporting
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_config_config_proto_init() }
func file_config_config_proto_init() {
	if File_config_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_config_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Recording); i {
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
		file_config_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reporting); i {
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
		file_config_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Instrumentation); i {
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
		file_config_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AgentConfig); i {
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
			RawDescriptor: file_config_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_config_proto_goTypes,
		DependencyIndexes: file_config_config_proto_depIdxs,
		MessageInfos:      file_config_config_proto_msgTypes,
	}.Build()
	File_config_config_proto = out.File
	file_config_config_proto_rawDesc = nil
	file_config_config_proto_goTypes = nil
	file_config_config_proto_depIdxs = nil
}
