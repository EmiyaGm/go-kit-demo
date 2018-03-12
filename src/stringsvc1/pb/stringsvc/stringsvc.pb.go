// Code generated by protoc-gen-go. DO NOT EDIT.
// source: stringsvc/stringsvc.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	stringsvc/stringsvc.proto

It has these top-level messages:
	UppercaseRequest
	UppercaseReply
	CreateRequest
	CreateReply
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type UppercaseRequest struct {
	S string `protobuf:"bytes,1,opt,name=s" json:"s,omitempty"`
}

func (m *UppercaseRequest) Reset()                    { *m = UppercaseRequest{} }
func (m *UppercaseRequest) String() string            { return proto.CompactTextString(m) }
func (*UppercaseRequest) ProtoMessage()               {}
func (*UppercaseRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UppercaseRequest) GetS() string {
	if m != nil {
		return m.S
	}
	return ""
}

type UppercaseReply struct {
	V   string `protobuf:"bytes,1,opt,name=v" json:"v,omitempty"`
	Err string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *UppercaseReply) Reset()                    { *m = UppercaseReply{} }
func (m *UppercaseReply) String() string            { return proto.CompactTextString(m) }
func (*UppercaseReply) ProtoMessage()               {}
func (*UppercaseReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UppercaseReply) GetV() string {
	if m != nil {
		return m.V
	}
	return ""
}

func (m *UppercaseReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type CreateRequest struct {
	ID     string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	FlowID uint32 `protobuf:"varint,2,opt,name=FlowID" json:"FlowID,omitempty"`
	Source string `protobuf:"bytes,3,opt,name=Source" json:"Source,omitempty"`
	Type   string `protobuf:"bytes,4,opt,name=Type" json:"Type,omitempty"`
}

func (m *CreateRequest) Reset()                    { *m = CreateRequest{} }
func (m *CreateRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()               {}
func (*CreateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CreateRequest) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *CreateRequest) GetFlowID() uint32 {
	if m != nil {
		return m.FlowID
	}
	return 0
}

func (m *CreateRequest) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

func (m *CreateRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type CreateReply struct {
	V   string `protobuf:"bytes,1,opt,name=v" json:"v,omitempty"`
	Err string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *CreateReply) Reset()                    { *m = CreateReply{} }
func (m *CreateReply) String() string            { return proto.CompactTextString(m) }
func (*CreateReply) ProtoMessage()               {}
func (*CreateReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CreateReply) GetV() string {
	if m != nil {
		return m.V
	}
	return ""
}

func (m *CreateReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto.RegisterType((*UppercaseRequest)(nil), "pb.UppercaseRequest")
	proto.RegisterType((*UppercaseReply)(nil), "pb.UppercaseReply")
	proto.RegisterType((*CreateRequest)(nil), "pb.createRequest")
	proto.RegisterType((*CreateReply)(nil), "pb.createReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Add service

type AddClient interface {
	// Uppercase 1 string.
	Uppercase(ctx context.Context, in *UppercaseRequest, opts ...grpc.CallOption) (*UppercaseReply, error)
	// create 更新报警记录
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error)
}

type addClient struct {
	cc *grpc.ClientConn
}

func NewAddClient(cc *grpc.ClientConn) AddClient {
	return &addClient{cc}
}

func (c *addClient) Uppercase(ctx context.Context, in *UppercaseRequest, opts ...grpc.CallOption) (*UppercaseReply, error) {
	out := new(UppercaseReply)
	err := grpc.Invoke(ctx, "/pb.Add/Uppercase", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error) {
	out := new(CreateReply)
	err := grpc.Invoke(ctx, "/pb.Add/create", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Add service

type AddServer interface {
	// Uppercase 1 string.
	Uppercase(context.Context, *UppercaseRequest) (*UppercaseReply, error)
	// create 更新报警记录
	Create(context.Context, *CreateRequest) (*CreateReply, error)
}

func RegisterAddServer(s *grpc.Server, srv AddServer) {
	s.RegisterService(&_Add_serviceDesc, srv)
}

func _Add_Uppercase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UppercaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServer).Uppercase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Add/Uppercase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServer).Uppercase(ctx, req.(*UppercaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Add_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Add/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Add_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Add",
	HandlerType: (*AddServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Uppercase",
			Handler:    _Add_Uppercase_Handler,
		},
		{
			MethodName: "create",
			Handler:    _Add_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stringsvc/stringsvc.proto",
}

func init() { proto.RegisterFile("stringsvc/stringsvc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x50, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x7d, 0x93, 0x96, 0x42, 0xef, 0xb3, 0xb5, 0x5e, 0x44, 0x62, 0x71, 0x51, 0x66, 0x25, 0x82,
	0x89, 0xe8, 0xce, 0x9d, 0x50, 0x84, 0xae, 0x84, 0xa8, 0x1f, 0x30, 0x99, 0x5e, 0x62, 0x20, 0x66,
	0xc6, 0x99, 0x69, 0x34, 0x5b, 0x17, 0xfe, 0x80, 0x9f, 0xe6, 0x2f, 0xf8, 0x21, 0x92, 0x69, 0x1a,
	0x6a, 0x57, 0xee, 0xce, 0x39, 0xf7, 0x9c, 0x3b, 0xf7, 0x0c, 0x1c, 0x5b, 0x67, 0xf2, 0x32, 0xb3,
	0x95, 0x8c, 0x3b, 0x14, 0x69, 0xa3, 0x9c, 0xc2, 0x40, 0xa7, 0xd3, 0x93, 0x4c, 0xa9, 0xac, 0xa0,
	0x58, 0xe8, 0x3c, 0x16, 0x65, 0xa9, 0x9c, 0x70, 0xb9, 0x2a, 0xed, 0xda, 0xc1, 0x67, 0x30, 0x79,
	0xd4, 0x9a, 0x8c, 0x14, 0x96, 0x12, 0x7a, 0x59, 0x91, 0x75, 0xb8, 0x07, 0xcc, 0x86, 0x6c, 0xc6,
	0x4e, 0x87, 0x09, 0xb3, 0xfc, 0x02, 0xc6, 0x5b, 0x0e, 0x5d, 0xd4, 0xcd, 0xbc, 0xda, 0xcc, 0x2b,
	0x9c, 0x40, 0x8f, 0x8c, 0x09, 0x03, 0xcf, 0x1b, 0xc8, 0x25, 0x8c, 0xa4, 0x21, 0xe1, 0xba, 0x85,
	0x63, 0x08, 0x16, 0xf3, 0x36, 0x11, 0x2c, 0xe6, 0x78, 0x04, 0x83, 0xdb, 0x42, 0xbd, 0x2e, 0xe6,
	0x3e, 0x35, 0x4a, 0x5a, 0xd6, 0xe8, 0xf7, 0x6a, 0x65, 0x24, 0x85, 0x3d, 0xef, 0x6d, 0x19, 0x22,
	0xf4, 0x1f, 0x6a, 0x4d, 0x61, 0xdf, 0xab, 0x1e, 0xf3, 0x73, 0xf8, 0xbf, 0x79, 0xe4, 0x0f, 0x37,
	0x5d, 0x7e, 0x30, 0xe8, 0xdd, 0x2c, 0x97, 0x78, 0x07, 0xc3, 0xae, 0x0d, 0x1e, 0x46, 0x3a, 0x8d,
	0x76, 0xeb, 0x4f, 0x71, 0x47, 0xd5, 0x45, 0xcd, 0xc3, 0xf7, 0xaf, 0xef, 0xcf, 0x00, 0xf9, 0x28,
	0xa6, 0x37, 0xf1, 0xac, 0x0b, 0x8a, 0x49, 0x3e, 0xa9, 0x6b, 0x76, 0x86, 0x11, 0x0c, 0xd6, 0x77,
	0xe0, 0x41, 0x93, 0xfb, 0x55, 0x7c, 0xba, 0xbf, 0x2d, 0x35, 0x7b, 0xfe, 0xa5, 0x03, 0xff, 0xef,
	0x57, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x77, 0x8d, 0xc7, 0x3c, 0xb6, 0x01, 0x00, 0x00,
}
