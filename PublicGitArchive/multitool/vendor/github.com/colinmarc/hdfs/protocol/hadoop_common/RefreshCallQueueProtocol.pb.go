// Code generated by protoc-gen-go.
// source: RefreshCallQueueProtocol.proto
// DO NOT EDIT!

package hadoop_common

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// *
//  Refresh callqueue request.
type RefreshCallQueueRequestProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RefreshCallQueueRequestProto) Reset()         { *m = RefreshCallQueueRequestProto{} }
func (m *RefreshCallQueueRequestProto) String() string { return proto.CompactTextString(m) }
func (*RefreshCallQueueRequestProto) ProtoMessage()    {}

// *
// void response.
type RefreshCallQueueResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RefreshCallQueueResponseProto) Reset()         { *m = RefreshCallQueueResponseProto{} }
func (m *RefreshCallQueueResponseProto) String() string { return proto.CompactTextString(m) }
func (*RefreshCallQueueResponseProto) ProtoMessage()    {}

func init() {
}
