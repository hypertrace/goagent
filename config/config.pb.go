// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: config.proto

package config

import (
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

// AgentConfig covers the config for agents
type AgentConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// serviceName identifies the service/process running
	ServiceName string `protobuf:"bytes,1,opt,name=serviceName,proto3" json:"serviceName,omitempty"`
	// reporting holds the reporting settings for the agent
	Reporting *Reporting `protobuf:"bytes,2,opt,name=reporting,proto3" json:"reporting,omitempty"`
	// dataCapture describes the data being captured by instrumentation
	DataCapture *DataCapture `protobuf:"bytes,3,opt,name=dataCapture,proto3" json:"dataCapture,omitempty"`
}

func (x *AgentConfig) Reset() {
	*x = AgentConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AgentConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentConfig) ProtoMessage() {}

func (x *AgentConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[0]
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
	return file_config_proto_rawDescGZIP(), []int{0}
}

func (x *AgentConfig) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *AgentConfig) GetReporting() *Reporting {
	if x != nil {
		return x.Reporting
	}
	return nil
}

func (x *AgentConfig) GetDataCapture() *DataCapture {
	if x != nil {
		return x.DataCapture
	}
	return nil
}

// Reporting covers the options related to the mechanics for sending data to the
// tracing server o collector.
type Reporting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// address represents the host for reporting the traes e.g. api.traceable.ai
	Address *wrappers.StringValue `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// secure when false, permits connecting to the trace endpoint without a certificate
	Secure *wrappers.BoolValue `protobuf:"bytes,2,opt,name=secure,proto3" json:"secure,omitempty"`
}

func (x *Reporting) Reset() {
	*x = Reporting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reporting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reporting) ProtoMessage() {}

func (x *Reporting) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[1]
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
	return file_config_proto_rawDescGZIP(), []int{1}
}

func (x *Reporting) GetAddress() *wrappers.StringValue {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *Reporting) GetSecure() *wrappers.BoolValue {
	if x != nil {
		return x.Secure
	}
	return nil
}

// Message describes what message should be considered for certain DataCapture option
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// request describes the outgoing/incoming message in a client/request operation
	Request *wrappers.BoolValue `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	// response describes the incoming/outgoing message in a client/request operation
	Response *wrappers.BoolValue `protobuf:"bytes,2,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{2}
}

func (x *Message) GetRequest() *wrappers.BoolValue {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *Message) GetResponse() *wrappers.BoolValue {
	if x != nil {
		return x.Response
	}
	return nil
}

// DataCapture describes the elements to be captured by the agent instrumentation
type DataCapture struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// httpHeaders enables/disables the capture of the request/response headers in HTTP
	HttpHeaders *Message `protobuf:"bytes,1,opt,name=httpHeaders,proto3" json:"httpHeaders,omitempty"`
	// httpBody enables/disables the capture of the request/response body in HTTP
	HttpBody *Message `protobuf:"bytes,2,opt,name=httpBody,proto3" json:"httpBody,omitempty"`
	// rpcMetadata enables/disables the capture of the request/response metadata in RPC
	RpcMetadata *Message `protobuf:"bytes,3,opt,name=rpcMetadata,proto3" json:"rpcMetadata,omitempty"`
	// rpcBody enables/disables the capture of the request/response body in RPC
	RpcBody *Message `protobuf:"bytes,4,opt,name=rpcBody,proto3" json:"rpcBody,omitempty"`
}

func (x *DataCapture) Reset() {
	*x = DataCapture{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataCapture) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataCapture) ProtoMessage() {}

func (x *DataCapture) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataCapture.ProtoReflect.Descriptor instead.
func (*DataCapture) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{3}
}

func (x *DataCapture) GetHttpHeaders() *Message {
	if x != nil {
		return x.HttpHeaders
	}
	return nil
}

func (x *DataCapture) GetHttpBody() *Message {
	if x != nil {
		return x.HttpBody
	}
	return nil
}

func (x *DataCapture) GetRpcMetadata() *Message {
	if x != nil {
		return x.RpcMetadata
	}
	return nil
}

func (x *DataCapture) GetRpcBody() *Message {
	if x != nil {
		return x.RpcBody
	}
	return nil
}

var File_config_proto protoreflect.FileDescriptor

var file_config_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89,
	0x01, 0x0a, 0x0b, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x20,
	0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x28, 0x0a, 0x09, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x52,
	0x09, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x2e, 0x0a, 0x0b, 0x64, 0x61,
	0x74, 0x61, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x52, 0x0b, 0x64,
	0x61, 0x74, 0x61, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x22, 0x77, 0x0a, 0x09, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x36, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x32, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x73, 0x65, 0x63,
	0x75, 0x72, 0x65, 0x22, 0x77, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x34,
	0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x07, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xaf, 0x01, 0x0a,
	0x0b, 0x44, 0x61, 0x74, 0x61, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x12, 0x2a, 0x0a, 0x0b,
	0x68, 0x74, 0x74, 0x70, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x0b, 0x68, 0x74, 0x74,
	0x70, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x24, 0x0a, 0x08, 0x68, 0x74, 0x74, 0x70,
	0x42, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x68, 0x74, 0x74, 0x70, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x2a,
	0x0a, 0x0b, 0x72, 0x70, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x0b, 0x72,
	0x70, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x22, 0x0a, 0x07, 0x72, 0x70,
	0x63, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x72, 0x70, 0x63, 0x42, 0x6f, 0x64, 0x79, 0x42, 0x48,
	0x0a, 0x1b, 0x6f, 0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65,
	0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5a, 0x29, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x79, 0x70, 0x65, 0x72, 0x74,
	0x72, 0x61, 0x63, 0x65, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_proto_rawDescOnce sync.Once
	file_config_proto_rawDescData = file_config_proto_rawDesc
)

func file_config_proto_rawDescGZIP() []byte {
	file_config_proto_rawDescOnce.Do(func() {
		file_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_proto_rawDescData)
	})
	return file_config_proto_rawDescData
}

var file_config_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_config_proto_goTypes = []interface{}{
	(*AgentConfig)(nil),          // 0: AgentConfig
	(*Reporting)(nil),            // 1: Reporting
	(*Message)(nil),              // 2: Message
	(*DataCapture)(nil),          // 3: DataCapture
	(*wrappers.StringValue)(nil), // 4: google.protobuf.StringValue
	(*wrappers.BoolValue)(nil),   // 5: google.protobuf.BoolValue
}
var file_config_proto_depIdxs = []int32{
	1,  // 0: AgentConfig.reporting:type_name -> Reporting
	3,  // 1: AgentConfig.dataCapture:type_name -> DataCapture
	4,  // 2: Reporting.address:type_name -> google.protobuf.StringValue
	5,  // 3: Reporting.secure:type_name -> google.protobuf.BoolValue
	5,  // 4: Message.request:type_name -> google.protobuf.BoolValue
	5,  // 5: Message.response:type_name -> google.protobuf.BoolValue
	2,  // 6: DataCapture.httpHeaders:type_name -> Message
	2,  // 7: DataCapture.httpBody:type_name -> Message
	2,  // 8: DataCapture.rpcMetadata:type_name -> Message
	2,  // 9: DataCapture.rpcBody:type_name -> Message
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_config_proto_init() }
func file_config_proto_init() {
	if File_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataCapture); i {
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
			RawDescriptor: file_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_proto_goTypes,
		DependencyIndexes: file_config_proto_depIdxs,
		MessageInfos:      file_config_proto_msgTypes,
	}.Build()
	File_config_proto = out.File
	file_config_proto_rawDesc = nil
	file_config_proto_goTypes = nil
	file_config_proto_depIdxs = nil
}
