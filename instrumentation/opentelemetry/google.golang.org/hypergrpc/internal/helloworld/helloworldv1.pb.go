package helloworld // import "github.com/hypertrace/goagent/instrumentation/hypertrace/google.golang.org/hypergrpc/examples/helloworld"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type HelloResponse struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,json=name,proto3" json:"product_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloResponse) Reset()         { *m = HelloResponse{} }
func (m *HelloResponse) String() string { return proto.CompactTextString(m) }
func (*HelloResponse) ProtoMessage()    {}

func (m *HelloResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CartItem.Unmarshal(m, b)
}
func (m *HelloResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CartItem.Marshal(b, m, deterministic)
}
func (dst *HelloResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CartItem.Merge(dst, src)
}
func (m *HelloResponse) XXX_Size() int {
	return xxx_messageInfo_CartItem.Size(m)
}
func (m *HelloResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CartItem.DiscardUnknown(m)
}

var xxx_messageInfo_CartItem proto.InternalMessageInfo

func (m *HelloResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}
