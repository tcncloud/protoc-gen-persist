// Code generated by protoc-gen-go.
// source: examples/spanner/delete/delete_examples.proto
// DO NOT EDIT!

/*
Package delete is a generated protocol buffer package.

It is generated from these files:
	examples/spanner/delete/delete_examples.proto

It has these top-level messages:
	Empty
*/
package delete

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/tcncloud/protoc-gen-persist/persist"

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

// This file is to show what kinds of spanner key ranges are generated
// from different types of delete queries.
type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Empty)(nil), "examples.Empty")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Delete service

type DeleteClient interface {
	// single primary key deletes (easy)
	DeleteEquals(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	DeleteGreater(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	DeleteLess(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	// multi primary key deletes (hard)
	DeleteMultiEquals(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type deleteClient struct {
	cc *grpc.ClientConn
}

func NewDeleteClient(cc *grpc.ClientConn) DeleteClient {
	return &deleteClient{cc}
}

func (c *deleteClient) DeleteEquals(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/examples.Delete/DeleteEquals", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deleteClient) DeleteGreater(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/examples.Delete/DeleteGreater", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deleteClient) DeleteLess(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/examples.Delete/DeleteLess", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deleteClient) DeleteMultiEquals(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/examples.Delete/DeleteMultiEquals", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Delete service

type DeleteServer interface {
	// single primary key deletes (easy)
	DeleteEquals(context.Context, *Empty) (*Empty, error)
	DeleteGreater(context.Context, *Empty) (*Empty, error)
	DeleteLess(context.Context, *Empty) (*Empty, error)
	// multi primary key deletes (hard)
	DeleteMultiEquals(context.Context, *Empty) (*Empty, error)
}

func RegisterDeleteServer(s *grpc.Server, srv DeleteServer) {
	s.RegisterService(&_Delete_serviceDesc, srv)
}

func _Delete_DeleteEquals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeleteServer).DeleteEquals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.Delete/DeleteEquals",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeleteServer).DeleteEquals(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Delete_DeleteGreater_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeleteServer).DeleteGreater(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.Delete/DeleteGreater",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeleteServer).DeleteGreater(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Delete_DeleteLess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeleteServer).DeleteLess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.Delete/DeleteLess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeleteServer).DeleteLess(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Delete_DeleteMultiEquals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeleteServer).DeleteMultiEquals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/examples.Delete/DeleteMultiEquals",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeleteServer).DeleteMultiEquals(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Delete_serviceDesc = grpc.ServiceDesc{
	ServiceName: "examples.Delete",
	HandlerType: (*DeleteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteEquals",
			Handler:    _Delete_DeleteEquals_Handler,
		},
		{
			MethodName: "DeleteGreater",
			Handler:    _Delete_DeleteGreater_Handler,
		},
		{
			MethodName: "DeleteLess",
			Handler:    _Delete_DeleteLess_Handler,
		},
		{
			MethodName: "DeleteMultiEquals",
			Handler:    _Delete_DeleteMultiEquals_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/spanner/delete/delete_examples.proto",
}

func init() { proto.RegisterFile("examples/spanner/delete/delete_examples.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4d, 0xad, 0x48, 0xcc,
	0x2d, 0xc8, 0x49, 0x2d, 0xd6, 0x2f, 0x2e, 0x48, 0xcc, 0xcb, 0x4b, 0x2d, 0xd2, 0x4f, 0x49, 0xcd,
	0x49, 0x2d, 0x49, 0x85, 0x52, 0xf1, 0x30, 0x69, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x0e,
	0x18, 0x5f, 0x4a, 0xb4, 0x20, 0xb5, 0xa8, 0x38, 0xb3, 0xb8, 0x44, 0x3f, 0xbf, 0xa0, 0x24, 0x33,
	0x3f, 0x0f, 0xaa, 0x40, 0x89, 0x9d, 0x8b, 0xd5, 0x35, 0xb7, 0xa0, 0xa4, 0xd2, 0xe8, 0x28, 0x33,
	0x17, 0x9b, 0x0b, 0xd8, 0x0c, 0xa1, 0x4c, 0x2e, 0x1e, 0x08, 0xcb, 0xb5, 0xb0, 0x34, 0x31, 0xa7,
	0x58, 0x88, 0x5f, 0x0f, 0x6e, 0x2a, 0x58, 0xad, 0x14, 0xba, 0x80, 0x92, 0x79, 0xd3, 0x8e, 0x89,
	0x4c, 0x46, 0x5c, 0x06, 0x2e, 0xae, 0x3e, 0xae, 0x21, 0xae, 0x0a, 0x21, 0x8e, 0x4e, 0x3e, 0xae,
	0x1a, 0xea, 0xb9, 0x95, 0x25, 0x89, 0x49, 0x39, 0xa9, 0xea, 0x9a, 0x0a, 0xc1, 0x21, 0x8e, 0x41,
	0x21, 0x1a, 0x86, 0x9a, 0x0a, 0xae, 0x7e, 0x2e, 0x20, 0xca, 0xdb, 0xd3, 0xcf, 0x45, 0xc3, 0xd9,
	0x59, 0x53, 0x28, 0x93, 0x8b, 0x17, 0x62, 0x95, 0x7b, 0x51, 0x6a, 0x62, 0x49, 0x6a, 0x11, 0x11,
	0x76, 0x99, 0x81, 0xec, 0x32, 0xe4, 0xd2, 0x27, 0xca, 0x2e, 0x98, 0x55, 0xfe, 0x9a, 0x42, 0xe9,
	0x5c, 0x5c, 0x10, 0xab, 0x7c, 0x52, 0x8b, 0xa9, 0xe7, 0x27, 0xa8, 0x97, 0x0c, 0xa0, 0x16, 0xf9,
	0x3b, 0x6b, 0x0a, 0x95, 0x72, 0x09, 0x42, 0x2c, 0xf2, 0x2d, 0xcd, 0x29, 0xc9, 0x24, 0x3a, 0x0c,
	0xed, 0x40, 0xf6, 0x59, 0x72, 0x99, 0xe3, 0xb0, 0x0f, 0xea, 0x2d, 0x1d, 0x05, 0x53, 0x03, 0xa8,
	0x9d, 0x10, 0x26, 0x2c, 0x28, 0xa5, 0x58, 0x27, 0xec, 0x98, 0xc8, 0xc4, 0xe8, 0xe4, 0x1e, 0xe5,
	0x9a, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x5f, 0x92, 0x9c, 0x97, 0x9c,
	0x93, 0x5f, 0x9a, 0xa2, 0x0f, 0x8e, 0xed, 0x64, 0xdd, 0xf4, 0xd4, 0x3c, 0x5d, 0x58, 0x42, 0xc0,
	0x91, 0x92, 0xac, 0x21, 0x54, 0x12, 0x1b, 0x58, 0x8b, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x5d,
	0xf8, 0x83, 0xb1, 0x72, 0x02, 0x00, 0x00,
}