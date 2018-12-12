// Code generated by protoc-gen-go. DO NOT EDIT.
// source: examples/spanner/import_tests/persist_and_go.proto

/*
Package import_tests is a generated protocol buffer package.

It is generated from these files:
	examples/spanner/import_tests/persist_and_go.proto

It has these top-level messages:
	ExampleTable
	AggExampleTables
*/
package import_tests

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/coltonmorris/protoc-gen-persist/persist"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"
import examples_test "github.com/coltonmorris/protoc-gen-persist/examples/test"

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

type ExampleTable struct {
	Id        int64                       `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	StartTime *google_protobuf1.Timestamp `protobuf:"bytes,2,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	Name      string                      `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *ExampleTable) Reset()                    { *m = ExampleTable{} }
func (m *ExampleTable) String() string            { return proto.CompactTextString(m) }
func (*ExampleTable) ProtoMessage()               {}
func (*ExampleTable) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ExampleTable) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ExampleTable) GetStartTime() *google_protobuf1.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *ExampleTable) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type AggExampleTables struct {
	ExampleTables []*ExampleTable `protobuf:"bytes,1,rep,name=example_tables,json=exampleTables" json:"example_tables,omitempty"`
}

func (m *AggExampleTables) Reset()                    { *m = AggExampleTables{} }
func (m *AggExampleTables) String() string            { return proto.CompactTextString(m) }
func (*AggExampleTables) ProtoMessage()               {}
func (*AggExampleTables) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *AggExampleTables) GetExampleTables() []*ExampleTable {
	if m != nil {
		return m.ExampleTables
	}
	return nil
}

func init() {
	proto.RegisterType((*ExampleTable)(nil), "import_tests.ExampleTable")
	proto.RegisterType((*AggExampleTables)(nil), "import_tests.AggExampleTables")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for MySpanner service

type MySpannerClient interface {
	UniaryInsert(ctx context.Context, in *examples_test.ExampleTable, opts ...grpc.CallOption) (*ExampleTable, error)
	ServerStream(ctx context.Context, in *examples_test.ExampleTable, opts ...grpc.CallOption) (MySpanner_ServerStreamClient, error)
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (MySpanner_ClientStreamClient, error)
}

type mySpannerClient struct {
	cc *grpc.ClientConn
}

func NewMySpannerClient(cc *grpc.ClientConn) MySpannerClient {
	return &mySpannerClient{cc}
}

func (c *mySpannerClient) UniaryInsert(ctx context.Context, in *examples_test.ExampleTable, opts ...grpc.CallOption) (*ExampleTable, error) {
	out := new(ExampleTable)
	err := grpc.Invoke(ctx, "/import_tests.MySpanner/UniaryInsert", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mySpannerClient) ServerStream(ctx context.Context, in *examples_test.ExampleTable, opts ...grpc.CallOption) (MySpanner_ServerStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_MySpanner_serviceDesc.Streams[0], c.cc, "/import_tests.MySpanner/ServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &mySpannerServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MySpanner_ServerStreamClient interface {
	Recv() (*ExampleTable, error)
	grpc.ClientStream
}

type mySpannerServerStreamClient struct {
	grpc.ClientStream
}

func (x *mySpannerServerStreamClient) Recv() (*ExampleTable, error) {
	m := new(ExampleTable)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mySpannerClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (MySpanner_ClientStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_MySpanner_serviceDesc.Streams[1], c.cc, "/import_tests.MySpanner/ClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &mySpannerClientStreamClient{stream}
	return x, nil
}

type MySpanner_ClientStreamClient interface {
	Send(*examples_test.ExampleTable) error
	CloseAndRecv() (*AggExampleTables, error)
	grpc.ClientStream
}

type mySpannerClientStreamClient struct {
	grpc.ClientStream
}

func (x *mySpannerClientStreamClient) Send(m *examples_test.ExampleTable) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mySpannerClientStreamClient) CloseAndRecv() (*AggExampleTables, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(AggExampleTables)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for MySpanner service

type MySpannerServer interface {
	UniaryInsert(context.Context, *examples_test.ExampleTable) (*ExampleTable, error)
	ServerStream(*examples_test.ExampleTable, MySpanner_ServerStreamServer) error
	ClientStream(MySpanner_ClientStreamServer) error
}

func RegisterMySpannerServer(s *grpc.Server, srv MySpannerServer) {
	s.RegisterService(&_MySpanner_serviceDesc, srv)
}

func _MySpanner_UniaryInsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(examples_test.ExampleTable)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MySpannerServer).UniaryInsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/import_tests.MySpanner/UniaryInsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MySpannerServer).UniaryInsert(ctx, req.(*examples_test.ExampleTable))
	}
	return interceptor(ctx, in, info, handler)
}

func _MySpanner_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(examples_test.ExampleTable)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MySpannerServer).ServerStream(m, &mySpannerServerStreamServer{stream})
}

type MySpanner_ServerStreamServer interface {
	Send(*ExampleTable) error
	grpc.ServerStream
}

type mySpannerServerStreamServer struct {
	grpc.ServerStream
}

func (x *mySpannerServerStreamServer) Send(m *ExampleTable) error {
	return x.ServerStream.SendMsg(m)
}

func _MySpanner_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MySpannerServer).ClientStream(&mySpannerClientStreamServer{stream})
}

type MySpanner_ClientStreamServer interface {
	SendAndClose(*AggExampleTables) error
	Recv() (*examples_test.ExampleTable, error)
	grpc.ServerStream
}

type mySpannerClientStreamServer struct {
	grpc.ServerStream
}

func (x *mySpannerClientStreamServer) SendAndClose(m *AggExampleTables) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mySpannerClientStreamServer) Recv() (*examples_test.ExampleTable, error) {
	m := new(examples_test.ExampleTable)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _MySpanner_serviceDesc = grpc.ServiceDesc{
	ServiceName: "import_tests.MySpanner",
	HandlerType: (*MySpannerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UniaryInsert",
			Handler:    _MySpanner_UniaryInsert_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStream",
			Handler:       _MySpanner_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _MySpanner_ClientStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "examples/spanner/import_tests/persist_and_go.proto",
}

func init() { proto.RegisterFile("examples/spanner/import_tests/persist_and_go.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 608 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x4f, 0x6b, 0x13, 0x4d,
	0x18, 0xef, 0x6c, 0xdf, 0xb7, 0xd0, 0x69, 0x1a, 0xeb, 0x80, 0x18, 0x22, 0xe8, 0xb0, 0x20, 0xa4,
	0x25, 0xdd, 0x2d, 0x15, 0x44, 0xa9, 0x35, 0xa4, 0x75, 0xa5, 0x8a, 0x45, 0xdd, 0x6c, 0x15, 0x7a,
	0x30, 0x4c, 0x92, 0x27, 0xeb, 0x60, 0x76, 0x66, 0xdd, 0x99, 0x2d, 0xe6, 0xaa, 0x07, 0xd1, 0x83,
	0x18, 0x8f, 0x7e, 0x06, 0xaf, 0x42, 0x6e, 0x7e, 0x08, 0xbf, 0x90, 0xcc, 0x6e, 0x52, 0x37, 0x41,
	0x8a, 0x8a, 0x97, 0xdd, 0x79, 0x66, 0x9e, 0xf9, 0xfd, 0xdb, 0x99, 0xc5, 0xdb, 0xf0, 0x8a, 0x45,
	0xf1, 0x00, 0x94, 0xab, 0x62, 0x26, 0x04, 0x24, 0x2e, 0x8f, 0x62, 0x99, 0xe8, 0xb6, 0x06, 0xa5,
	0x95, 0x1b, 0x43, 0xa2, 0xb8, 0xd2, 0x6d, 0x26, 0x7a, 0xed, 0x50, 0x3a, 0x71, 0x22, 0xb5, 0x24,
	0xa5, 0x62, 0x4b, 0xf5, 0xc2, 0xa4, 0xc7, 0x95, 0xb1, 0xe6, 0x52, 0xa8, 0xbc, 0xa9, 0x7a, 0x25,
	0x94, 0x32, 0x1c, 0x80, 0x9b, 0x55, 0x9d, 0xb4, 0xef, 0x6a, 0x1e, 0x81, 0xd2, 0x2c, 0x8a, 0x27,
	0x0d, 0x95, 0x53, 0x66, 0x83, 0x93, 0x3d, 0xf2, 0x15, 0x3b, 0xc2, 0x25, 0x2f, 0x5f, 0x0b, 0x58,
	0x67, 0x00, 0xa4, 0x8c, 0x2d, 0xde, 0xab, 0x20, 0x8a, 0x6a, 0x8b, 0xbe, 0xc5, 0x7b, 0xe4, 0x26,
	0xc6, 0x4a, 0x33, 0x23, 0x80, 0x47, 0x50, 0xb1, 0x28, 0xaa, 0xad, 0x6c, 0x57, 0x9d, 0x9c, 0xcf,
	0x99, 0xf2, 0x39, 0xc1, 0x94, 0xcf, 0x5f, 0xce, 0xba, 0x4d, 0x4d, 0x08, 0xfe, 0x4f, 0xb0, 0x08,
	0x2a, 0x8b, 0x14, 0xd5, 0x96, 0xfd, 0x6c, 0x6c, 0x1f, 0xe1, 0xb5, 0x66, 0x18, 0x16, 0x19, 0x15,
	0x69, 0xe2, 0xf2, 0x44, 0x5e, 0x5b, 0x67, 0x33, 0x15, 0x44, 0x17, 0x33, 0x9a, 0xa2, 0x77, 0xa7,
	0xb8, 0xc9, 0x5f, 0x85, 0x22, 0xc4, 0xf6, 0xb7, 0xff, 0xf1, 0xf2, 0xe1, 0xb0, 0x95, 0xa7, 0x4a,
	0xbe, 0x23, 0x5c, 0x3a, 0x12, 0x9c, 0x25, 0xc3, 0x7b, 0x42, 0x41, 0xa2, 0xc9, 0x25, 0x67, 0xea,
	0xdf, 0xc9, 0xac, 0x17, 0xa1, 0xaa, 0x67, 0xd0, 0xd8, 0xef, 0xd0, 0xeb, 0xf1, 0xc8, 0x7a, 0x83,
	0xf0, 0xfd, 0x1c, 0x89, 0x72, 0xa1, 0x25, 0x9d, 0x51, 0x4b, 0x6b, 0xbc, 0x57, 0xa7, 0x3f, 0x33,
	0xaa, 0x53, 0x63, 0x74, 0x9d, 0xd2, 0x27, 0x6c, 0x90, 0x82, 0xa2, 0xb5, 0x46, 0x9d, 0x36, 0xea,
	0xd4, 0xee, 0x30, 0xc1, 0x04, 0x53, 0xf6, 0x3a, 0x31, 0xb1, 0x16, 0x42, 0xf5, 0x09, 0x2e, 0x07,
	0xa0, 0xf4, 0x1e, 0xf4, 0x65, 0x02, 0x07, 0x52, 0xbe, 0x20, 0x0b, 0xc7, 0xe7, 0xf1, 0xaa, 0x99,
	0x6b, 0xf6, 0x35, 0x24, 0xf9, 0x14, 0xf9, 0x80, 0x70, 0xa9, 0x05, 0xc9, 0x09, 0x24, 0x2d, 0x9d,
	0x00, 0x8b, 0xfe, 0xde, 0xd5, 0xbe, 0x31, 0x75, 0x1b, 0x6f, 0xb5, 0xbc, 0x07, 0xde, 0x7e, 0x40,
	0x37, 0xe8, 0x5d, 0xff, 0xe1, 0xe1, 0x9c, 0xab, 0xa7, 0x07, 0x9e, 0xef, 0x15, 0x6c, 0xd1, 0x5b,
	0xb4, 0x51, 0x54, 0xbd, 0x85, 0xc8, 0x57, 0x84, 0x4b, 0xfb, 0x03, 0x0e, 0x42, 0xff, 0x8e, 0xa0,
	0xcb, 0xb3, 0x82, 0xe6, 0x4f, 0x81, 0xdd, 0x37, 0xa2, 0x18, 0xbe, 0x71, 0xf4, 0xe8, 0x4e, 0x33,
	0xf0, 0xe6, 0xd4, 0xb4, 0xbc, 0xa0, 0xa0, 0x65, 0xb7, 0x91, 0x87, 0xbc, 0xdb, 0x98, 0xc8, 0xe4,
	0xbd, 0xdd, 0x19, 0x71, 0x24, 0x3b, 0x6c, 0x26, 0xe8, 0xe3, 0x73, 0x78, 0xa5, 0x19, 0x86, 0x3e,
	0xbc, 0x4c, 0x0d, 0x2d, 0x59, 0xa8, 0xa1, 0xea, 0x27, 0xf4, 0x7e, 0x3c, 0xb2, 0xde, 0x22, 0xfc,
	0x0c, 0x9f, 0x71, 0x92, 0xd7, 0x56, 0xec, 0xa5, 0xc3, 0xa1, 0x29, 0x37, 0xae, 0x87, 0x5c, 0x3f,
	0x4f, 0x3b, 0x4e, 0x57, 0x46, 0xae, 0xee, 0x8a, 0xee, 0x40, 0xa6, 0xbd, 0xfc, 0xae, 0x75, 0x37,
	0x43, 0x10, 0x9b, 0xd3, 0x4b, 0x79, 0x7a, 0xc9, 0xa2, 0xa1, 0x51, 0x81, 0xaf, 0xe2, 0x8b, 0x73,
	0x89, 0x98, 0x6f, 0xea, 0x89, 0x34, 0x5a, 0x2b, 0x1b, 0x70, 0x33, 0xda, 0x58, 0xf8, 0x38, 0x1e,
	0x59, 0x68, 0xef, 0x0b, 0xfa, 0x3c, 0x1e, 0x59, 0xc1, 0x9f, 0x32, 0xfd, 0xf2, 0x47, 0xa2, 0x20,
	0x39, 0xe1, 0x5d, 0xd8, 0x99, 0xbc, 0x8f, 0x1f, 0xff, 0x0b, 0xd4, 0x9d, 0x62, 0xd1, 0x59, 0xca,
	0xb6, 0x5f, 0xfb, 0x11, 0x00, 0x00, 0xff, 0xff, 0xf2, 0x21, 0x66, 0x88, 0xd3, 0x04, 0x00, 0x00,
}
