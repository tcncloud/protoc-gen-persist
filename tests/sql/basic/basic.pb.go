// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tests/sql/basic/basic.proto

package basic // import "github.com/tcncloud/protoc-gen-persist/tests/sql/basic"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/golang/protobuf/protoc-gen-go/descriptor"
import _ "github.com/golang/protobuf/ptypes/timestamp"
import _ "github.com/tcncloud/protoc-gen-persist/persist"
import test "github.com/tcncloud/protoc-gen-persist/tests/test"

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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_basic_7d6929f464812518, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type IntVal struct {
	Val                  int64    `protobuf:"varint,1,opt,name=val,proto3" json:"val,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IntVal) Reset()         { *m = IntVal{} }
func (m *IntVal) String() string { return proto.CompactTextString(m) }
func (*IntVal) ProtoMessage()    {}
func (*IntVal) Descriptor() ([]byte, []int) {
	return fileDescriptor_basic_7d6929f464812518, []int{1}
}
func (m *IntVal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntVal.Unmarshal(m, b)
}
func (m *IntVal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntVal.Marshal(b, m, deterministic)
}
func (dst *IntVal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntVal.Merge(dst, src)
}
func (m *IntVal) XXX_Size() int {
	return xxx_messageInfo_IntVal.Size(m)
}
func (m *IntVal) XXX_DiscardUnknown() {
	xxx_messageInfo_IntVal.DiscardUnknown(m)
}

var xxx_messageInfo_IntVal proto.InternalMessageInfo

func (m *IntVal) GetVal() int64 {
	if m != nil {
		return m.Val
	}
	return 0
}

type BadReturn struct {
	No                   []*IntVal `protobuf:"bytes,1,rep,name=no,proto3" json:"no,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *BadReturn) Reset()         { *m = BadReturn{} }
func (m *BadReturn) String() string { return proto.CompactTextString(m) }
func (*BadReturn) ProtoMessage()    {}
func (*BadReturn) Descriptor() ([]byte, []int) {
	return fileDescriptor_basic_7d6929f464812518, []int{2}
}
func (m *BadReturn) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BadReturn.Unmarshal(m, b)
}
func (m *BadReturn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BadReturn.Marshal(b, m, deterministic)
}
func (dst *BadReturn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BadReturn.Merge(dst, src)
}
func (m *BadReturn) XXX_Size() int {
	return xxx_messageInfo_BadReturn.Size(m)
}
func (m *BadReturn) XXX_DiscardUnknown() {
	xxx_messageInfo_BadReturn.DiscardUnknown(m)
}

var xxx_messageInfo_BadReturn proto.InternalMessageInfo

func (m *BadReturn) GetNo() []*IntVal {
	if m != nil {
		return m.No
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "tests.Empty")
	proto.RegisterType((*IntVal)(nil), "tests.IntVal")
	proto.RegisterType((*BadReturn)(nil), "tests.BadReturn")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AmazingClient is the client API for Amazing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AmazingClient interface {
	UniarySelect(ctx context.Context, in *test.PartialTable, opts ...grpc.CallOption) (*test.ExampleTable, error)
	UniarySelectWithHooks(ctx context.Context, in *test.PartialTable, opts ...grpc.CallOption) (*test.ExampleTable, error)
	ServerStream(ctx context.Context, in *test.Name, opts ...grpc.CallOption) (Amazing_ServerStreamClient, error)
	ServerStreamWithHooks(ctx context.Context, in *test.Name, opts ...grpc.CallOption) (Amazing_ServerStreamWithHooksClient, error)
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (Amazing_ClientStreamClient, error)
	ClientStreamWithHook(ctx context.Context, opts ...grpc.CallOption) (Amazing_ClientStreamWithHookClient, error)
	UnImplementedPersistMethod(ctx context.Context, in *test.ExampleTable, opts ...grpc.CallOption) (*test.ExampleTable, error)
	NoGenerationForBadReturnTypes(ctx context.Context, in *test.ExampleTable, opts ...grpc.CallOption) (*BadReturn, error)
}

type amazingClient struct {
	cc *grpc.ClientConn
}

func NewAmazingClient(cc *grpc.ClientConn) AmazingClient {
	return &amazingClient{cc}
}

func (c *amazingClient) UniarySelect(ctx context.Context, in *test.PartialTable, opts ...grpc.CallOption) (*test.ExampleTable, error) {
	out := new(test.ExampleTable)
	err := c.cc.Invoke(ctx, "/tests.Amazing/UniarySelect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *amazingClient) UniarySelectWithHooks(ctx context.Context, in *test.PartialTable, opts ...grpc.CallOption) (*test.ExampleTable, error) {
	out := new(test.ExampleTable)
	err := c.cc.Invoke(ctx, "/tests.Amazing/UniarySelectWithHooks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *amazingClient) ServerStream(ctx context.Context, in *test.Name, opts ...grpc.CallOption) (Amazing_ServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Amazing_serviceDesc.Streams[0], "/tests.Amazing/ServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &amazingServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Amazing_ServerStreamClient interface {
	Recv() (*test.ExampleTable, error)
	grpc.ClientStream
}

type amazingServerStreamClient struct {
	grpc.ClientStream
}

func (x *amazingServerStreamClient) Recv() (*test.ExampleTable, error) {
	m := new(test.ExampleTable)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *amazingClient) ServerStreamWithHooks(ctx context.Context, in *test.Name, opts ...grpc.CallOption) (Amazing_ServerStreamWithHooksClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Amazing_serviceDesc.Streams[1], "/tests.Amazing/ServerStreamWithHooks", opts...)
	if err != nil {
		return nil, err
	}
	x := &amazingServerStreamWithHooksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Amazing_ServerStreamWithHooksClient interface {
	Recv() (*test.ExampleTable, error)
	grpc.ClientStream
}

type amazingServerStreamWithHooksClient struct {
	grpc.ClientStream
}

func (x *amazingServerStreamWithHooksClient) Recv() (*test.ExampleTable, error) {
	m := new(test.ExampleTable)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *amazingClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (Amazing_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Amazing_serviceDesc.Streams[2], "/tests.Amazing/ClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &amazingClientStreamClient{stream}
	return x, nil
}

type Amazing_ClientStreamClient interface {
	Send(*test.ExampleTable) error
	CloseAndRecv() (*test.NumRows, error)
	grpc.ClientStream
}

type amazingClientStreamClient struct {
	grpc.ClientStream
}

func (x *amazingClientStreamClient) Send(m *test.ExampleTable) error {
	return x.ClientStream.SendMsg(m)
}

func (x *amazingClientStreamClient) CloseAndRecv() (*test.NumRows, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(test.NumRows)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *amazingClient) ClientStreamWithHook(ctx context.Context, opts ...grpc.CallOption) (Amazing_ClientStreamWithHookClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Amazing_serviceDesc.Streams[3], "/tests.Amazing/ClientStreamWithHook", opts...)
	if err != nil {
		return nil, err
	}
	x := &amazingClientStreamWithHookClient{stream}
	return x, nil
}

type Amazing_ClientStreamWithHookClient interface {
	Send(*test.ExampleTable) error
	CloseAndRecv() (*test.Ids, error)
	grpc.ClientStream
}

type amazingClientStreamWithHookClient struct {
	grpc.ClientStream
}

func (x *amazingClientStreamWithHookClient) Send(m *test.ExampleTable) error {
	return x.ClientStream.SendMsg(m)
}

func (x *amazingClientStreamWithHookClient) CloseAndRecv() (*test.Ids, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(test.Ids)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *amazingClient) UnImplementedPersistMethod(ctx context.Context, in *test.ExampleTable, opts ...grpc.CallOption) (*test.ExampleTable, error) {
	out := new(test.ExampleTable)
	err := c.cc.Invoke(ctx, "/tests.Amazing/UnImplementedPersistMethod", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *amazingClient) NoGenerationForBadReturnTypes(ctx context.Context, in *test.ExampleTable, opts ...grpc.CallOption) (*BadReturn, error) {
	out := new(BadReturn)
	err := c.cc.Invoke(ctx, "/tests.Amazing/NoGenerationForBadReturnTypes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AmazingServer is the server API for Amazing service.
type AmazingServer interface {
	UniarySelect(context.Context, *test.PartialTable) (*test.ExampleTable, error)
	UniarySelectWithHooks(context.Context, *test.PartialTable) (*test.ExampleTable, error)
	ServerStream(*test.Name, Amazing_ServerStreamServer) error
	ServerStreamWithHooks(*test.Name, Amazing_ServerStreamWithHooksServer) error
	ClientStream(Amazing_ClientStreamServer) error
	ClientStreamWithHook(Amazing_ClientStreamWithHookServer) error
	UnImplementedPersistMethod(context.Context, *test.ExampleTable) (*test.ExampleTable, error)
	NoGenerationForBadReturnTypes(context.Context, *test.ExampleTable) (*BadReturn, error)
}

func RegisterAmazingServer(s *grpc.Server, srv AmazingServer) {
	s.RegisterService(&_Amazing_serviceDesc, srv)
}

func _Amazing_UniarySelect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(test.PartialTable)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AmazingServer).UniarySelect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tests.Amazing/UniarySelect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AmazingServer).UniarySelect(ctx, req.(*test.PartialTable))
	}
	return interceptor(ctx, in, info, handler)
}

func _Amazing_UniarySelectWithHooks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(test.PartialTable)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AmazingServer).UniarySelectWithHooks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tests.Amazing/UniarySelectWithHooks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AmazingServer).UniarySelectWithHooks(ctx, req.(*test.PartialTable))
	}
	return interceptor(ctx, in, info, handler)
}

func _Amazing_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(test.Name)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AmazingServer).ServerStream(m, &amazingServerStreamServer{stream})
}

type Amazing_ServerStreamServer interface {
	Send(*test.ExampleTable) error
	grpc.ServerStream
}

type amazingServerStreamServer struct {
	grpc.ServerStream
}

func (x *amazingServerStreamServer) Send(m *test.ExampleTable) error {
	return x.ServerStream.SendMsg(m)
}

func _Amazing_ServerStreamWithHooks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(test.Name)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AmazingServer).ServerStreamWithHooks(m, &amazingServerStreamWithHooksServer{stream})
}

type Amazing_ServerStreamWithHooksServer interface {
	Send(*test.ExampleTable) error
	grpc.ServerStream
}

type amazingServerStreamWithHooksServer struct {
	grpc.ServerStream
}

func (x *amazingServerStreamWithHooksServer) Send(m *test.ExampleTable) error {
	return x.ServerStream.SendMsg(m)
}

func _Amazing_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AmazingServer).ClientStream(&amazingClientStreamServer{stream})
}

type Amazing_ClientStreamServer interface {
	SendAndClose(*test.NumRows) error
	Recv() (*test.ExampleTable, error)
	grpc.ServerStream
}

type amazingClientStreamServer struct {
	grpc.ServerStream
}

func (x *amazingClientStreamServer) SendAndClose(m *test.NumRows) error {
	return x.ServerStream.SendMsg(m)
}

func (x *amazingClientStreamServer) Recv() (*test.ExampleTable, error) {
	m := new(test.ExampleTable)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Amazing_ClientStreamWithHook_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AmazingServer).ClientStreamWithHook(&amazingClientStreamWithHookServer{stream})
}

type Amazing_ClientStreamWithHookServer interface {
	SendAndClose(*test.Ids) error
	Recv() (*test.ExampleTable, error)
	grpc.ServerStream
}

type amazingClientStreamWithHookServer struct {
	grpc.ServerStream
}

func (x *amazingClientStreamWithHookServer) SendAndClose(m *test.Ids) error {
	return x.ServerStream.SendMsg(m)
}

func (x *amazingClientStreamWithHookServer) Recv() (*test.ExampleTable, error) {
	m := new(test.ExampleTable)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Amazing_UnImplementedPersistMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(test.ExampleTable)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AmazingServer).UnImplementedPersistMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tests.Amazing/UnImplementedPersistMethod",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AmazingServer).UnImplementedPersistMethod(ctx, req.(*test.ExampleTable))
	}
	return interceptor(ctx, in, info, handler)
}

func _Amazing_NoGenerationForBadReturnTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(test.ExampleTable)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AmazingServer).NoGenerationForBadReturnTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tests.Amazing/NoGenerationForBadReturnTypes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AmazingServer).NoGenerationForBadReturnTypes(ctx, req.(*test.ExampleTable))
	}
	return interceptor(ctx, in, info, handler)
}

var _Amazing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tests.Amazing",
	HandlerType: (*AmazingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UniarySelect",
			Handler:    _Amazing_UniarySelect_Handler,
		},
		{
			MethodName: "UniarySelectWithHooks",
			Handler:    _Amazing_UniarySelectWithHooks_Handler,
		},
		{
			MethodName: "UnImplementedPersistMethod",
			Handler:    _Amazing_UnImplementedPersistMethod_Handler,
		},
		{
			MethodName: "NoGenerationForBadReturnTypes",
			Handler:    _Amazing_NoGenerationForBadReturnTypes_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStream",
			Handler:       _Amazing_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ServerStreamWithHooks",
			Handler:       _Amazing_ServerStreamWithHooks_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _Amazing_ClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ClientStreamWithHook",
			Handler:       _Amazing_ClientStreamWithHook_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "tests/sql/basic/basic.proto",
}

func init() { proto.RegisterFile("tests/sql/basic/basic.proto", fileDescriptor_basic_7d6929f464812518) }

var fileDescriptor_basic_7d6929f464812518 = []byte{
	// 683 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x51, 0x4f, 0x1a, 0x4b,
	0x14, 0xc7, 0x5d, 0x88, 0x7a, 0x1d, 0xd1, 0xcb, 0x1d, 0x25, 0xd7, 0x6e, 0x63, 0x3b, 0xd9, 0xf4,
	0x01, 0x89, 0x82, 0xb1, 0x8f, 0xad, 0x06, 0xb5, 0x6b, 0x25, 0x51, 0xb4, 0x0b, 0x6a, 0xd3, 0x17,
	0x3a, 0xec, 0x8e, 0x30, 0xe9, 0xee, 0xcc, 0x76, 0x66, 0xb0, 0xa5, 0x8f, 0x26, 0xed, 0x83, 0x4f,
	0xe5, 0x2b, 0xf4, 0x3b, 0x34, 0xe1, 0xe3, 0x35, 0xb3, 0x0b, 0xba, 0xd0, 0x40, 0xd2, 0xbe, 0x9c,
	0x0c, 0x9c, 0x73, 0x7e, 0xff, 0xff, 0x4e, 0xce, 0x1c, 0xf0, 0x58, 0x11, 0xa9, 0x64, 0x49, 0x7e,
	0xf4, 0x4b, 0x4d, 0x2c, 0xa9, 0x1b, 0xc7, 0x62, 0x28, 0xb8, 0xe2, 0x70, 0x36, 0x4a, 0x9a, 0xb9,
	0x90, 0x08, 0x49, 0xa5, 0x2a, 0xf1, 0x50, 0x51, 0xce, 0x64, 0x9c, 0x35, 0x9f, 0xb6, 0x38, 0x6f,
	0xf9, 0xa4, 0x14, 0xfd, 0x6a, 0x76, 0xae, 0x4b, 0x8a, 0x06, 0x44, 0x2a, 0x1c, 0x84, 0x83, 0x02,
	0x34, 0x5e, 0xe0, 0x11, 0xe9, 0x0a, 0x1a, 0x2a, 0x2e, 0x06, 0x15, 0xb9, 0x58, 0x5d, 0xc7, 0x28,
	0xc4, 0x7f, 0x5b, 0xf3, 0x60, 0xd6, 0x0e, 0x42, 0xd5, 0xb5, 0x4c, 0x30, 0x57, 0x61, 0xea, 0x12,
	0xfb, 0x30, 0x0b, 0xd2, 0x37, 0xd8, 0x5f, 0x33, 0x90, 0x91, 0x4f, 0x3b, 0xfa, 0x68, 0x15, 0xc0,
	0xc2, 0x01, 0xf6, 0x1c, 0xa2, 0x3a, 0x82, 0xc1, 0x75, 0x90, 0x62, 0x7c, 0xcd, 0x40, 0xe9, 0xfc,
	0xe2, 0xce, 0x52, 0x31, 0xa2, 0x16, 0xe3, 0x4e, 0x27, 0xc5, 0xf8, 0xce, 0xed, 0x02, 0x98, 0xdf,
	0x0f, 0xf0, 0x17, 0xca, 0x5a, 0xf0, 0x0a, 0x64, 0x2e, 0x18, 0xc5, 0xa2, 0x5b, 0x23, 0x3e, 0x71,
	0x15, 0x5c, 0x1b, 0x94, 0x47, 0xfa, 0xe7, 0x58, 0x28, 0x8a, 0xfd, 0x3a, 0x6e, 0xfa, 0xc4, 0x1c,
	0xc9, 0xd8, 0x9f, 0x71, 0x10, 0xfa, 0x24, 0xca, 0x58, 0x2b, 0x3f, 0xfa, 0xbd, 0xd4, 0x32, 0xc8,
	0xc8, 0x08, 0xd1, 0x68, 0x76, 0x1b, 0xd4, 0x83, 0x1e, 0xc8, 0x25, 0xc1, 0x57, 0x54, 0xb5, 0x8f,
	0x39, 0xff, 0x20, 0xff, 0x4a, 0xe1, 0x7f, 0xad, 0x00, 0x47, 0x15, 0xce, 0x8d, 0xb7, 0x06, 0xac,
	0x81, 0x4c, 0x8d, 0x88, 0x1b, 0x22, 0x6a, 0x4a, 0x10, 0x1c, 0xc0, 0x6c, 0x12, 0x51, 0xc5, 0xc1,
	0x34, 0x68, 0x4e, 0x43, 0xb3, 0x60, 0xf9, 0x01, 0xca, 0x70, 0x40, 0xb6, 0x0d, 0xf8, 0x1e, 0xe4,
	0x92, 0xd0, 0x07, 0xeb, 0x7f, 0x42, 0x7f, 0xa4, 0xe9, 0xab, 0xe3, 0x74, 0x6d, 0x7a, 0xdb, 0x80,
	0x55, 0x90, 0x39, 0xf4, 0x29, 0x61, 0x6a, 0x60, 0x7b, 0x22, 0xc6, 0x5c, 0x19, 0x91, 0xec, 0x04,
	0x0e, 0xff, 0x24, 0xad, 0x25, 0xcd, 0xfe, 0x07, 0xcc, 0x51, 0x26, 0x89, 0x50, 0x79, 0x03, 0x5e,
	0x82, 0xd5, 0x24, 0x6f, 0xe8, 0x78, 0x0a, 0xf7, 0xdf, 0x64, 0xa6, 0xe2, 0x49, 0xeb, 0x3f, 0xcd,
	0xcc, 0x0c, 0x99, 0xda, 0x67, 0xde, 0x80, 0x0e, 0x30, 0x2f, 0x58, 0x45, 0x77, 0x05, 0x84, 0x29,
	0xe2, 0x9d, 0xc7, 0xa3, 0x7f, 0x4a, 0x54, 0x9b, 0x7b, 0x53, 0xe8, 0x93, 0xaf, 0x65, 0x06, 0x9e,
	0x82, 0xf5, 0x2a, 0x7f, 0x4d, 0x18, 0x11, 0x58, 0xbf, 0x9f, 0x23, 0x2e, 0xee, 0x07, 0xb7, 0xde,
	0x0d, 0x89, 0x9c, 0x82, 0x1d, 0xde, 0xff, 0x7d, 0x83, 0x35, 0x63, 0xfe, 0x4c, 0xdf, 0xf6, 0x7b,
	0xa9, 0xaf, 0x69, 0x70, 0x67, 0x80, 0xc3, 0x9a, 0x7d, 0x62, 0x1f, 0xd6, 0x51, 0x01, 0x5d, 0x0b,
	0x1e, 0x20, 0x12, 0xf7, 0x36, 0x94, 0x6e, 0x46, 0x57, 0x6d, 0x22, 0x08, 0xa2, 0xde, 0x6e, 0x99,
	0x7a, 0x68, 0xbf, 0xfa, 0x0a, 0x49, 0x85, 0x85, 0x6a, 0xe8, 0x67, 0xba, 0x57, 0x7e, 0x38, 0x43,
	0xe3, 0x99, 0x39, 0x32, 0x67, 0xd6, 0xc4, 0x81, 0x2d, 0x4c, 0x74, 0x0a, 0x02, 0xb0, 0x79, 0xef,
	0xe5, 0xc8, 0x39, 0x3b, 0x1d, 0xf7, 0x72, 0x6c, 0x3b, 0x36, 0xd2, 0x03, 0xb1, 0x5b, 0xd6, 0x51,
	0x8b, 0x8e, 0x4d, 0x8a, 0xf5, 0xdb, 0xb0, 0x4d, 0x91, 0xfb, 0x66, 0x80, 0x37, 0x95, 0x6a, 0xcd,
	0x76, 0xea, 0xa8, 0x52, 0xad, 0x9f, 0x8d, 0xa9, 0xe5, 0xa9, 0xb7, 0x99, 0xf8, 0xda, 0xcd, 0x48,
	0x78, 0x03, 0x5d, 0xee, 0x9f, 0x5c, 0xd8, 0x35, 0x94, 0x2f, 0xeb, 0x74, 0x39, 0x99, 0x8f, 0x3c,
	0x6d, 0x68, 0x53, 0x83, 0x71, 0xb0, 0x26, 0x4a, 0x17, 0x32, 0x83, 0x4c, 0xb4, 0xa2, 0xee, 0xfa,
	0xbd, 0x14, 0x02, 0x4f, 0x80, 0x59, 0x8c, 0x97, 0x5d, 0x71, 0xb8, 0xec, 0x8a, 0xf5, 0xe1, 0x36,
	0xcc, 0x2e, 0x7e, 0xef, 0xf7, 0x52, 0x33, 0x07, 0x7b, 0xef, 0x5e, 0xb6, 0xa8, 0x6a, 0x77, 0x9a,
	0x45, 0x97, 0x07, 0x25, 0xe5, 0x32, 0xd7, 0xe7, 0x1d, 0x2f, 0xde, 0x8e, 0xee, 0x56, 0x8b, 0xb0,
	0xad, 0xe1, 0x9e, 0x1d, 0xdb, 0xc9, 0x2f, 0xa2, 0xd8, 0x9c, 0x8b, 0x2a, 0x9f, 0xff, 0x0a, 0x00,
	0x00, 0xff, 0xff, 0xee, 0xac, 0x59, 0xe1, 0xb3, 0x05, 0x00, 0x00,
}
