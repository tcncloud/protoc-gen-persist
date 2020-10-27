// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/extra_persist_service_in_pkg.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	imports "github.com/tcncloud/protoc-gen-persistexamples/user_sql/pb/imports"
	_ "github.com/tcncloud/protoc-gen-persist/persist"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EUser struct {
	Id                   int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Friends              *Friends             `protobuf:"bytes,3,opt,name=friends,proto3" json:"friends,omitempty"`
	CreatedOn            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=created_on,json=createdOn,proto3" json:"created_on,omitempty"`
	Id2                  int32                `protobuf:"varint,5,opt,name=id2,proto3" json:"id2,omitempty"`
	Counts               []int64              `protobuf:"varint,7,rep,packed,name=counts,proto3" json:"counts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *EUser) Reset()         { *m = EUser{} }
func (m *EUser) String() string { return proto.CompactTextString(m) }
func (*EUser) ProtoMessage()    {}
func (*EUser) Descriptor() ([]byte, []int) {
	return fileDescriptor_42d08e3827593ce4, []int{0}
}

func (m *EUser) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EUser.Unmarshal(m, b)
}
func (m *EUser) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EUser.Marshal(b, m, deterministic)
}
func (m *EUser) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EUser.Merge(m, src)
}
func (m *EUser) XXX_Size() int {
	return xxx_messageInfo_EUser.Size(m)
}
func (m *EUser) XXX_DiscardUnknown() {
	xxx_messageInfo_EUser.DiscardUnknown(m)
}

var xxx_messageInfo_EUser proto.InternalMessageInfo

func (m *EUser) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *EUser) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EUser) GetFriends() *Friends {
	if m != nil {
		return m.Friends
	}
	return nil
}

func (m *EUser) GetCreatedOn() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedOn
	}
	return nil
}

func (m *EUser) GetId2() int32 {
	if m != nil {
		return m.Id2
	}
	return 0
}

func (m *EUser) GetCounts() []int64 {
	if m != nil {
		return m.Counts
	}
	return nil
}

func init() {
	proto.RegisterType((*EUser)(nil), "pb.EUser")
}

func init() {
	proto.RegisterFile("pb/extra_persist_service_in_pkg.proto", fileDescriptor_42d08e3827593ce4)
}

var fileDescriptor_42d08e3827593ce4 = []byte{
	// 720 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xdf, 0x4e, 0xdb, 0x48,
	0x14, 0xc6, 0xb1, 0x63, 0xfe, 0xe4, 0x44, 0x8b, 0xa2, 0xd9, 0x5d, 0xd6, 0xeb, 0x0b, 0x18, 0x45,
	0x8b, 0x36, 0xac, 0x42, 0x82, 0xb2, 0x6d, 0xa5, 0x5e, 0x20, 0xe1, 0x84, 0x01, 0x2c, 0x42, 0x40,
	0x8e, 0x69, 0x8b, 0xaa, 0xca, 0xb5, 0xe3, 0x21, 0x1a, 0xd5, 0xb1, 0x2d, 0x7b, 0x82, 0xca, 0x6d,
	0x6f, 0x2a, 0xf5, 0xaa, 0x79, 0x81, 0x5e, 0xf4, 0x25, 0x9a, 0xeb, 0x3e, 0x50, 0x9f, 0xa1, 0xb2,
	0xc7, 0x24, 0xa0, 0x02, 0x2a, 0xbd, 0xf2, 0x99, 0xf1, 0x99, 0xef, 0xfb, 0xcd, 0x39, 0x67, 0x60,
	0x3d, 0x72, 0x1b, 0xf4, 0x2d, 0x8f, 0x1d, 0x3b, 0xa2, 0x71, 0xc2, 0x12, 0x6e, 0x27, 0x34, 0xbe,
	0x60, 0x7d, 0x6a, 0xb3, 0xc0, 0x8e, 0xde, 0x0c, 0xea, 0x51, 0x1c, 0xf2, 0x10, 0xc9, 0x91, 0xab,
	0xfd, 0x99, 0x67, 0x34, 0xc2, 0x88, 0xb3, 0x30, 0x48, 0xc4, 0x2f, 0x6d, 0x6d, 0x10, 0x86, 0x03,
	0x9f, 0x36, 0xb2, 0x95, 0x3b, 0x3a, 0x6f, 0x70, 0x36, 0xa4, 0x09, 0x77, 0x86, 0x51, 0x9e, 0xb0,
	0xc2, 0x86, 0x51, 0x18, 0xf3, 0xa4, 0x21, 0xbe, 0xd4, 0xcb, 0xf7, 0x7f, 0x8b, 0xdc, 0xc6, 0x28,
	0xa1, 0xb1, 0x58, 0x56, 0xbe, 0x48, 0x30, 0x4f, 0x4e, 0x13, 0x1a, 0xa3, 0x65, 0x90, 0x99, 0xa7,
	0x4a, 0x58, 0xaa, 0x16, 0x4c, 0x99, 0x79, 0x08, 0x81, 0x12, 0x38, 0x43, 0xaa, 0xca, 0x58, 0xaa,
	0x16, 0xcd, 0x2c, 0x46, 0xeb, 0xb0, 0x78, 0x1e, 0x33, 0x1a, 0x78, 0x89, 0x5a, 0xc0, 0x52, 0xb5,
	0xd4, 0x2c, 0xd5, 0x23, 0xb7, 0xbe, 0x27, 0xb6, 0xcc, 0xab, 0x7f, 0xe8, 0x29, 0x40, 0x3f, 0xa6,
	0x0e, 0xa7, 0x9e, 0x1d, 0x06, 0xaa, 0x92, 0x65, 0x6a, 0x75, 0x41, 0x5c, 0xbf, 0x22, 0xae, 0x5b,
	0x57, 0xc4, 0x66, 0x31, 0xcf, 0x3e, 0x0e, 0x50, 0x19, 0x0a, 0xcc, 0x6b, 0xaa, 0xf3, 0x58, 0xaa,
	0xce, 0x9b, 0x69, 0x88, 0x56, 0x60, 0xa1, 0x1f, 0x8e, 0x02, 0x9e, 0xa8, 0x8b, 0xb8, 0x50, 0x2d,
	0x98, 0xf9, 0xaa, 0xf9, 0x75, 0x09, 0x8a, 0x24, 0xad, 0x61, 0x8f, 0xc6, 0x17, 0x68, 0x17, 0x4a,
	0xed, 0x4c, 0xc4, 0x72, 0x5c, 0x9f, 0xa2, 0x62, 0xca, 0x45, 0x86, 0x11, 0xbf, 0xd4, 0x66, 0x61,
	0x65, 0xed, 0xf3, 0x64, 0x2c, 0x6b, 0xa0, 0x8a, 0xd2, 0x0b, 0x53, 0x3b, 0x2d, 0x46, 0x62, 0xf3,
	0xec, 0x58, 0x17, 0x14, 0xcb, 0xb2, 0x2c, 0x74, 0x0f, 0xec, 0x83, 0xf4, 0x36, 0x41, 0x21, 0x84,
	0x90, 0x9f, 0xd4, 0x43, 0x07, 0xa0, 0x18, 0x86, 0x61, 0xa0, 0x3f, 0xd2, 0xad, 0xbc, 0x7f, 0x75,
	0x23, 0xef, 0xdf, 0x83, 0x8c, 0x77, 0xa1, 0xb4, 0x4f, 0xb9, 0xee, 0xfb, 0x69, 0x6b, 0x93, 0xeb,
	0xe5, 0x58, 0x4a, 0xc3, 0x74, 0xb7, 0xb2, 0x9a, 0x8a, 0xfc, 0x0d, 0xbf, 0x0b, 0x91, 0x01, 0xe5,
	0xb6, 0xe3, 0xfb, 0x42, 0xe5, 0x44, 0x7a, 0x21, 0x6d, 0x49, 0x68, 0x17, 0x96, 0x7b, 0xd4, 0xa7,
	0x7d, 0x9e, 0xe6, 0xb7, 0x2e, 0x0d, 0x0f, 0x4d, 0x4f, 0xff, 0xa8, 0xf3, 0x97, 0xd0, 0x49, 0xb2,
	0x03, 0x99, 0x8c, 0xed, 0x5e, 0xda, 0xcc, 0x43, 0x1b, 0xb0, 0x7c, 0x1a, 0x79, 0x0e, 0xa7, 0xba,
	0xef, 0x77, 0x9d, 0x21, 0xbd, 0x03, 0x67, 0x6e, 0x4b, 0xd2, 0xbe, 0x29, 0xef, 0x26, 0x63, 0xf9,
	0x93, 0x02, 0x13, 0x09, 0x0e, 0xdb, 0x26, 0xd1, 0x2d, 0x82, 0x2d, 0xbd, 0xd5, 0x21, 0x38, 0xc3,
	0xaa, 0x32, 0x0f, 0xb3, 0x80, 0xd3, 0x01, 0x8d, 0xf1, 0x89, 0x69, 0x1c, 0xe9, 0xe6, 0x19, 0x3e,
	0x24, 0x67, 0x35, 0x9c, 0x0e, 0x27, 0x7e, 0xa6, 0x9b, 0xed, 0x03, 0xdd, 0xac, 0x3e, 0xde, 0xda,
	0xa8, 0xe1, 0x7c, 0x14, 0x71, 0xeb, 0xcc, 0x22, 0x7a, 0x0d, 0x9e, 0xcc, 0xe6, 0xf1, 0x66, 0x1e,
	0xf3, 0x9a, 0xb8, 0x77, 0xa4, 0x77, 0x3a, 0x46, 0xd7, 0xaa, 0x61, 0x31, 0x5a, 0xf8, 0xe5, 0xab,
	0x96, 0xb1, 0x6f, 0x74, 0xad, 0x0d, 0x24, 0xfd, 0xa3, 0xdd, 0x59, 0xe4, 0xca, 0xec, 0x1a, 0xff,
	0xcd, 0x42, 0x78, 0x0d, 0x8f, 0x7a, 0xa4, 0x43, 0xda, 0x16, 0x66, 0x9e, 0x20, 0x9b, 0xd2, 0xd4,
	0xf0, 0x0c, 0x43, 0x58, 0xef, 0x99, 0xc7, 0x47, 0xe2, 0x6e, 0xda, 0x6d, 0x7d, 0xb8, 0x6e, 0x31,
	0xad, 0x14, 0x8c, 0xa0, 0xfd, 0x2b, 0x0e, 0xf8, 0xf9, 0x01, 0x31, 0x09, 0x66, 0x1e, 0xde, 0xc6,
	0x3b, 0xcc, 0x4b, 0xef, 0x76, 0x57, 0xcf, 0x2a, 0x53, 0xb3, 0x6b, 0xb6, 0xef, 0x25, 0xd8, 0x14,
	0x8d, 0xcc, 0xf5, 0x12, 0xca, 0x45, 0xe5, 0xb7, 0xf1, 0x4e, 0xf6, 0xbd, 0x69, 0x80, 0xe1, 0x5f,
	0x93, 0x58, 0xa7, 0x66, 0xd7, 0xe8, 0xee, 0xdf, 0x4f, 0x9a, 0xa2, 0xac, 0x08, 0x94, 0x51, 0x66,
	0x20, 0x50, 0xd2, 0xfc, 0xdb, 0x48, 0x3e, 0x4c, 0xc6, 0xf2, 0x0e, 0xac, 0xc2, 0x3d, 0xaf, 0xaa,
	0x5c, 0x02, 0x15, 0xb2, 0x67, 0xd4, 0xf3, 0x59, 0x9f, 0xf6, 0x78, 0xcc, 0x82, 0xc1, 0x89, 0x13,
	0x3b, 0xc3, 0x72, 0x09, 0x94, 0x72, 0x41, 0x2d, 0x7c, 0x9c, 0x8c, 0xe5, 0x39, 0x77, 0x21, 0x3b,
	0xfa, 0xff, 0xf7, 0x00, 0x00, 0x00, 0xff, 0xff, 0xda, 0x8c, 0xe5, 0x71, 0x91, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ExtraServClient is the client API for ExtraServ service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ExtraServClient interface {
	CreateTable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	TTTT(ctx context.Context, in *timestamp.Timestamp, opts ...grpc.CallOption) (*Empty, error)
	EEEE(ctx context.Context, in *timestamp.Timestamp, opts ...grpc.CallOption) (*Empty, error)
	IIII(ctx context.Context, in *imports.Imported, opts ...grpc.CallOption) (*Empty, error)
	GetAllUsers(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ExtraServ_GetAllUsersClient, error)
	SelectUserById(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	UpdateAllNames(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ExtraServ_UpdateAllNamesClient, error)
}

type extraServClient struct {
	cc *grpc.ClientConn
}

func NewExtraServClient(cc *grpc.ClientConn) ExtraServClient {
	return &extraServClient{cc}
}

func (c *extraServClient) CreateTable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.ExtraServ/CreateTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *extraServClient) TTTT(ctx context.Context, in *timestamp.Timestamp, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.ExtraServ/TTTT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *extraServClient) EEEE(ctx context.Context, in *timestamp.Timestamp, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.ExtraServ/EEEE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *extraServClient) IIII(ctx context.Context, in *imports.Imported, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.ExtraServ/IIII", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *extraServClient) GetAllUsers(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ExtraServ_GetAllUsersClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ExtraServ_serviceDesc.Streams[0], "/pb.ExtraServ/GetAllUsers", opts...)
	if err != nil {
		return nil, err
	}
	x := &extraServGetAllUsersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExtraServ_GetAllUsersClient interface {
	Recv() (*User, error)
	grpc.ClientStream
}

type extraServGetAllUsersClient struct {
	grpc.ClientStream
}

func (x *extraServGetAllUsersClient) Recv() (*User, error) {
	m := new(User)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *extraServClient) SelectUserById(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/pb.ExtraServ/SelectUserById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *extraServClient) UpdateAllNames(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ExtraServ_UpdateAllNamesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ExtraServ_serviceDesc.Streams[1], "/pb.ExtraServ/UpdateAllNames", opts...)
	if err != nil {
		return nil, err
	}
	x := &extraServUpdateAllNamesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExtraServ_UpdateAllNamesClient interface {
	Recv() (*User, error)
	grpc.ClientStream
}

type extraServUpdateAllNamesClient struct {
	grpc.ClientStream
}

func (x *extraServUpdateAllNamesClient) Recv() (*User, error) {
	m := new(User)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExtraServServer is the server API for ExtraServ service.
type ExtraServServer interface {
	CreateTable(context.Context, *Empty) (*Empty, error)
	TTTT(context.Context, *timestamp.Timestamp) (*Empty, error)
	EEEE(context.Context, *timestamp.Timestamp) (*Empty, error)
	IIII(context.Context, *imports.Imported) (*Empty, error)
	GetAllUsers(*Empty, ExtraServ_GetAllUsersServer) error
	SelectUserById(context.Context, *User) (*User, error)
	UpdateAllNames(*Empty, ExtraServ_UpdateAllNamesServer) error
}

func RegisterExtraServServer(s *grpc.Server, srv ExtraServServer) {
	s.RegisterService(&_ExtraServ_serviceDesc, srv)
}

func _ExtraServ_CreateTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtraServServer).CreateTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ExtraServ/CreateTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtraServServer).CreateTable(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExtraServ_TTTT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(timestamp.Timestamp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtraServServer).TTTT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ExtraServ/TTTT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtraServServer).TTTT(ctx, req.(*timestamp.Timestamp))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExtraServ_EEEE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(timestamp.Timestamp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtraServServer).EEEE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ExtraServ/EEEE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtraServServer).EEEE(ctx, req.(*timestamp.Timestamp))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExtraServ_IIII_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(imports.Imported)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtraServServer).IIII(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ExtraServ/IIII",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtraServServer).IIII(ctx, req.(*imports.Imported))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExtraServ_GetAllUsers_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExtraServServer).GetAllUsers(m, &extraServGetAllUsersServer{stream})
}

type ExtraServ_GetAllUsersServer interface {
	Send(*User) error
	grpc.ServerStream
}

type extraServGetAllUsersServer struct {
	grpc.ServerStream
}

func (x *extraServGetAllUsersServer) Send(m *User) error {
	return x.ServerStream.SendMsg(m)
}

func _ExtraServ_SelectUserById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtraServServer).SelectUserById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.ExtraServ/SelectUserById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtraServServer).SelectUserById(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExtraServ_UpdateAllNames_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExtraServServer).UpdateAllNames(m, &extraServUpdateAllNamesServer{stream})
}

type ExtraServ_UpdateAllNamesServer interface {
	Send(*User) error
	grpc.ServerStream
}

type extraServUpdateAllNamesServer struct {
	grpc.ServerStream
}

func (x *extraServUpdateAllNamesServer) Send(m *User) error {
	return x.ServerStream.SendMsg(m)
}

var _ExtraServ_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.ExtraServ",
	HandlerType: (*ExtraServServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTable",
			Handler:    _ExtraServ_CreateTable_Handler,
		},
		{
			MethodName: "TTTT",
			Handler:    _ExtraServ_TTTT_Handler,
		},
		{
			MethodName: "EEEE",
			Handler:    _ExtraServ_EEEE_Handler,
		},
		{
			MethodName: "IIII",
			Handler:    _ExtraServ_IIII_Handler,
		},
		{
			MethodName: "SelectUserById",
			Handler:    _ExtraServ_SelectUserById_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAllUsers",
			Handler:       _ExtraServ_GetAllUsers_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UpdateAllNames",
			Handler:       _ExtraServ_UpdateAllNames_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pb/extra_persist_service_in_pkg.proto",
}
