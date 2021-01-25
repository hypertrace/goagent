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

// PropagationFormat represents the propagation formats supported by agents
type PropagationFormat int32

const (
	// B3 propagation format, agents should support both multi and single value formats
	// see https://github.com/openzipkin/b3-propagation
	PropagationFormat_B3 PropagationFormat = 0
	// W3C Propagation format
	// see https://www.w3.org/TR/trace-context/
	PropagationFormat_TRACECONTEXT PropagationFormat = 1
)

// Enum value maps for PropagationFormat.
var (
	PropagationFormat_name = map[int32]string{
		0: "B3",
		1: "TRACECONTEXT",
	}
	PropagationFormat_value = map[string]int32{
		"B3":           0,
		"TRACECONTEXT": 1,
	}
)

func (x PropagationFormat) Enum() *PropagationFormat {
	p := new(PropagationFormat)
	*p = x
	return p
}

func (x PropagationFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PropagationFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_config_proto_enumTypes[0].Descriptor()
}

func (PropagationFormat) Type() protoreflect.EnumType {
	return &file_config_proto_enumTypes[0]
}

func (x PropagationFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PropagationFormat.Descriptor instead.
func (PropagationFormat) EnumDescriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{0}
}

// AgentConfig covers the config for agents.
// The config uses wrappers for primitive types to allow nullable values.
// The nullable values are used for instance to explicitly disable data capture or secure connection.
// Since the wrappers change structure of the objects the marshalling and unmarshalling
// have to be done via protobuf marshallers.
type AgentConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// service_name identifies the service/process running e.g. "my service"
	ServiceName *wrappers.StringValue `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// reporting holds the reporting settings for the agent
	Reporting *Reporting `protobuf:"bytes,2,opt,name=reporting,proto3" json:"reporting,omitempty"`
	// data_capture describes the data being captured by instrumentation
	DataCapture *DataCapture `protobuf:"bytes,3,opt,name=data_capture,json=dataCapture,proto3" json:"data_capture,omitempty"`
	// propagation_formats list the supported propagation formats
	PropagationFormats []PropagationFormat `protobuf:"varint,4,rep,packed,name=propagation_formats,json=propagationFormats,proto3,enum=org.hypertrace.agent.config.PropagationFormat" json:"propagation_formats,omitempty"`
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

func (x *AgentConfig) GetServiceName() *wrappers.StringValue {
	if x != nil {
		return x.ServiceName
	}
	return nil
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

func (x *AgentConfig) GetPropagationFormats() []PropagationFormat {
	if x != nil {
		return x.PropagationFormats
	}
	return nil
}

// Reporting covers the options related to the mechanics for sending data to the
// tracing server o collector.
type Reporting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// endpoint represents the endpoint for reporting the traces e.g. http://api.traceable.ai:9411/api/v2/spans
	Endpoint *wrappers.StringValue `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// when `true`, connects to endpoints over TLS.
	Secure *wrappers.BoolValue `protobuf:"bytes,2,opt,name=secure,proto3" json:"secure,omitempty"`
	// user specific token to access Traceable API
	Token *wrappers.StringValue `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	// opa describes the setting related to the Open Policy Agent
	Opa *Opa `protobuf:"bytes,4,opt,name=opa,proto3" json:"opa,omitempty"`
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

func (x *Reporting) GetEndpoint() *wrappers.StringValue {
	if x != nil {
		return x.Endpoint
	}
	return nil
}

func (x *Reporting) GetSecure() *wrappers.BoolValue {
	if x != nil {
		return x.Secure
	}
	return nil
}

func (x *Reporting) GetToken() *wrappers.StringValue {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *Reporting) GetOpa() *Opa {
	if x != nil {
		return x.Opa
	}
	return nil
}

// Opa covers the options related to the mechanics for getting Open Policy Agent configuration file.
// The client should use secure and token option from reporting settings.
type Opa struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// endpoint represents the endpoint for polling OPA config file e.g. http://opa.traceableai:8181/
	Endpoint *wrappers.StringValue `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// poll period in seconds to query OPA service
	PollPeriodSeconds *wrappers.Int32Value `protobuf:"bytes,2,opt,name=poll_period_seconds,json=pollPeriodSeconds,proto3" json:"poll_period_seconds,omitempty"`
	// when `true` Open Policy Agent evaluation is enabled to block request
	Enabled *wrappers.BoolValue `protobuf:"bytes,3,opt,name=enabled,proto3" json:"enabled,omitempty"`
}

func (x *Opa) Reset() {
	*x = Opa{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Opa) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Opa) ProtoMessage() {}

func (x *Opa) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Opa.ProtoReflect.Descriptor instead.
func (*Opa) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{2}
}

func (x *Opa) GetEndpoint() *wrappers.StringValue {
	if x != nil {
		return x.Endpoint
	}
	return nil
}

func (x *Opa) GetPollPeriodSeconds() *wrappers.Int32Value {
	if x != nil {
		return x.PollPeriodSeconds
	}
	return nil
}

func (x *Opa) GetEnabled() *wrappers.BoolValue {
	if x != nil {
		return x.Enabled
	}
	return nil
}

// Message describes what message should be considered for certain DataCapture option
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// when `false` it disables the capture for the request in a client/request operation
	Request *wrappers.BoolValue `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	// when `false` it disables the capture for the response in a client/request operation
	Response *wrappers.BoolValue `protobuf:"bytes,2,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{3}
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

	// http_headers enables/disables the capture of the request/response headers in HTTP
	HttpHeaders *Message `protobuf:"bytes,1,opt,name=http_headers,json=httpHeaders,proto3" json:"http_headers,omitempty"`
	// http_body enables/disables the capture of the request/response body in HTTP
	HttpBody *Message `protobuf:"bytes,2,opt,name=http_body,json=httpBody,proto3" json:"http_body,omitempty"`
	// rpc_metadata enables/disables the capture of the request/response metadata in RPC
	RpcMetadata *Message `protobuf:"bytes,3,opt,name=rpc_metadata,json=rpcMetadata,proto3" json:"rpc_metadata,omitempty"`
	// rpc_body enables/disables the capture of the request/response body in RPC
	RpcBody *Message `protobuf:"bytes,4,opt,name=rpc_body,json=rpcBody,proto3" json:"rpc_body,omitempty"`
	// maximum size of captured body in bytes. Default should be 131_072 (128 KiB).
	BodyMaxSizeBytes *wrappers.Int32Value `protobuf:"bytes,5,opt,name=body_max_size_bytes,json=bodyMaxSizeBytes,proto3" json:"body_max_size_bytes,omitempty"`
}

func (x *DataCapture) Reset() {
	*x = DataCapture{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataCapture) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataCapture) ProtoMessage() {}

func (x *DataCapture) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[4]
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
	return file_config_proto_rawDescGZIP(), []int{4}
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

func (x *DataCapture) GetBodyMaxSizeBytes() *wrappers.Int32Value {
	if x != nil {
		return x.BodyMaxSizeBytes
	}
	return nil
}

var File_config_proto protoreflect.FileDescriptor

var file_config_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b,
	0x6f, 0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc2, 0x02, 0x0a, 0x0b,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3f, 0x0a, 0x0c, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x44, 0x0a, 0x09,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x26, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65,
	0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x09, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69,
	0x6e, 0x67, 0x12, 0x4b, 0x0a, 0x0c, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x63, 0x61, 0x70, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x68,
	0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x43, 0x61, 0x70, 0x74, 0x75,
	0x72, 0x65, 0x52, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x12,
	0x5f, 0x0a, 0x13, 0x70, 0x72, 0x6f, 0x70, 0x61, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x6f,
	0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x50, 0x72, 0x6f, 0x70, 0x61,
	0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x12, 0x70, 0x72,
	0x6f, 0x70, 0x61, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x73,
	0x22, 0xe1, 0x01, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x38,
	0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x75,
	0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x12, 0x32, 0x0a, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x32, 0x0a, 0x03, 0x6f, 0x70, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x6f, 0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4f, 0x70, 0x61, 0x52,
	0x03, 0x6f, 0x70, 0x61, 0x22, 0xc2, 0x01, 0x0a, 0x03, 0x4f, 0x70, 0x61, 0x12, 0x38, 0x0a, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x65, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x4b, 0x0a, 0x13, 0x70, 0x6f, 0x6c, 0x6c, 0x5f, 0x70,
	0x65, 0x72, 0x69, 0x6f, 0x64, 0x5f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x11, 0x70, 0x6f, 0x6c, 0x6c, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x53, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x22, 0x77, 0x0a, 0x07, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x42,
	0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0xef, 0x02, 0x0a, 0x0b, 0x44, 0x61, 0x74, 0x61, 0x43, 0x61, 0x70, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x47, 0x0a, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x68,
	0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x0b,
	0x68, 0x74, 0x74, 0x70, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x41, 0x0a, 0x09, 0x68,
	0x74, 0x74, 0x70, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x68, 0x74, 0x74, 0x70, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x47,
	0x0a, 0x0c, 0x72, 0x70, 0x63, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x0b, 0x72, 0x70, 0x63, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3f, 0x0a, 0x08, 0x72, 0x70, 0x63, 0x5f, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6f, 0x72, 0x67, 0x2e,
	0x68, 0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x07, 0x72, 0x70, 0x63, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x4a, 0x0a, 0x13, 0x62, 0x6f, 0x64, 0x79,
	0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x10, 0x62, 0x6f, 0x64, 0x79, 0x4d, 0x61, 0x78, 0x53, 0x69, 0x7a, 0x65, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x2a, 0x2d, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x70, 0x61, 0x67, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x06, 0x0a, 0x02, 0x42, 0x33, 0x10,
	0x00, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x52, 0x41, 0x43, 0x45, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x58,
	0x54, 0x10, 0x01, 0x42, 0x48, 0x0a, 0x1b, 0x6f, 0x72, 0x67, 0x2e, 0x68, 0x79, 0x70, 0x65, 0x72,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68,
	0x79, 0x70, 0x65, 0x72, 0x74, 0x72, 0x61, 0x63, 0x65, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2d,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_config_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_config_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_config_proto_goTypes = []interface{}{
	(PropagationFormat)(0),       // 0: org.hypertrace.agent.config.PropagationFormat
	(*AgentConfig)(nil),          // 1: org.hypertrace.agent.config.AgentConfig
	(*Reporting)(nil),            // 2: org.hypertrace.agent.config.Reporting
	(*Opa)(nil),                  // 3: org.hypertrace.agent.config.Opa
	(*Message)(nil),              // 4: org.hypertrace.agent.config.Message
	(*DataCapture)(nil),          // 5: org.hypertrace.agent.config.DataCapture
	(*wrappers.StringValue)(nil), // 6: google.protobuf.StringValue
	(*wrappers.BoolValue)(nil),   // 7: google.protobuf.BoolValue
	(*wrappers.Int32Value)(nil),  // 8: google.protobuf.Int32Value
}
var file_config_proto_depIdxs = []int32{
	6,  // 0: org.hypertrace.agent.config.AgentConfig.service_name:type_name -> google.protobuf.StringValue
	2,  // 1: org.hypertrace.agent.config.AgentConfig.reporting:type_name -> org.hypertrace.agent.config.Reporting
	5,  // 2: org.hypertrace.agent.config.AgentConfig.data_capture:type_name -> org.hypertrace.agent.config.DataCapture
	0,  // 3: org.hypertrace.agent.config.AgentConfig.propagation_formats:type_name -> org.hypertrace.agent.config.PropagationFormat
	6,  // 4: org.hypertrace.agent.config.Reporting.endpoint:type_name -> google.protobuf.StringValue
	7,  // 5: org.hypertrace.agent.config.Reporting.secure:type_name -> google.protobuf.BoolValue
	6,  // 6: org.hypertrace.agent.config.Reporting.token:type_name -> google.protobuf.StringValue
	3,  // 7: org.hypertrace.agent.config.Reporting.opa:type_name -> org.hypertrace.agent.config.Opa
	6,  // 8: org.hypertrace.agent.config.Opa.endpoint:type_name -> google.protobuf.StringValue
	8,  // 9: org.hypertrace.agent.config.Opa.poll_period_seconds:type_name -> google.protobuf.Int32Value
	7,  // 10: org.hypertrace.agent.config.Opa.enabled:type_name -> google.protobuf.BoolValue
	7,  // 11: org.hypertrace.agent.config.Message.request:type_name -> google.protobuf.BoolValue
	7,  // 12: org.hypertrace.agent.config.Message.response:type_name -> google.protobuf.BoolValue
	4,  // 13: org.hypertrace.agent.config.DataCapture.http_headers:type_name -> org.hypertrace.agent.config.Message
	4,  // 14: org.hypertrace.agent.config.DataCapture.http_body:type_name -> org.hypertrace.agent.config.Message
	4,  // 15: org.hypertrace.agent.config.DataCapture.rpc_metadata:type_name -> org.hypertrace.agent.config.Message
	4,  // 16: org.hypertrace.agent.config.DataCapture.rpc_body:type_name -> org.hypertrace.agent.config.Message
	8,  // 17: org.hypertrace.agent.config.DataCapture.body_max_size_bytes:type_name -> google.protobuf.Int32Value
	18, // [18:18] is the sub-list for method output_type
	18, // [18:18] is the sub-list for method input_type
	18, // [18:18] is the sub-list for extension type_name
	18, // [18:18] is the sub-list for extension extendee
	0,  // [0:18] is the sub-list for field type_name
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
			switch v := v.(*Opa); i {
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
		file_config_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_proto_goTypes,
		DependencyIndexes: file_config_proto_depIdxs,
		EnumInfos:         file_config_proto_enumTypes,
		MessageInfos:      file_config_proto_msgTypes,
	}.Build()
	File_config_proto = out.File
	file_config_proto_rawDesc = nil
	file_config_proto_goTypes = nil
	file_config_proto_depIdxs = nil
}
