// Code generated by protoc-gen-go. DO NOT EDIT.
// source: persist/options.proto

package persist

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type PersistenceOptions int32

const (
	// SQL Query
	PersistenceOptions_SQL     PersistenceOptions = 0
	PersistenceOptions_SPANNER PersistenceOptions = 1
)

var PersistenceOptions_name = map[int32]string{
	0: "SQL",
	1: "SPANNER",
}

var PersistenceOptions_value = map[string]int32{
	"SQL":     0,
	"SPANNER": 1,
}

func (x PersistenceOptions) Enum() *PersistenceOptions {
	p := new(PersistenceOptions)
	*p = x
	return p
}

func (x PersistenceOptions) String() string {
	return proto.EnumName(PersistenceOptions_name, int32(x))
}

func (x *PersistenceOptions) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PersistenceOptions_value, data, "PersistenceOptions")
	if err != nil {
		return err
	}
	*x = PersistenceOptions(value)
	return nil
}

func (PersistenceOptions) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6479c8984e34e8f5, []int{0}
}

type QueryOpts struct {
	Queries              []*QLImpl `protobuf:"bytes,1,rep,name=queries" json:"queries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *QueryOpts) Reset()         { *m = QueryOpts{} }
func (m *QueryOpts) String() string { return proto.CompactTextString(m) }
func (*QueryOpts) ProtoMessage()    {}
func (*QueryOpts) Descriptor() ([]byte, []int) {
	return fileDescriptor_6479c8984e34e8f5, []int{0}
}

func (m *QueryOpts) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryOpts.Unmarshal(m, b)
}
func (m *QueryOpts) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryOpts.Marshal(b, m, deterministic)
}
func (m *QueryOpts) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryOpts.Merge(m, src)
}
func (m *QueryOpts) XXX_Size() int {
	return xxx_messageInfo_QueryOpts.Size(m)
}
func (m *QueryOpts) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryOpts.DiscardUnknown(m)
}

var xxx_messageInfo_QueryOpts proto.InternalMessageInfo

func (m *QueryOpts) GetQueries() []*QLImpl {
	if m != nil {
		return m.Queries
	}
	return nil
}

type QLImpl struct {
	// the query string with numeric placeholders for parameters
	// its an array to allow the query to span across multiple lines but it
	// will be joined and used as a single sql query string at generation time
	Query []string `protobuf:"bytes,1,rep,name=query" json:"query,omitempty"`
	// if provided, persist will rewrite the query string in the generated code
	// replacing "@field_name" (no quotes) with "?" or "$<position>"
	// if unprovided, persist will not rewrite the query string
	PmStrategy *string `protobuf:"bytes,2,opt,name=pm_strategy,json=pmStrategy" json:"pm_strategy,omitempty"`
	// name of this query.  must be unique to the service.
	Name *string `protobuf:"bytes,3,req,name=name" json:"name,omitempty"`
	// the message type that matches the parameters
	// Input rpc messages will be converted to this type
	// they will be used in the parameters in the query
	// The INTERFACE of this message will be used for parameters
	// in the generated query function.
	// if absent, this query takes no  parameters.
	// The query does not have to use all the fields of this type as parameters,
	// but it cannot use any parameter NOT listed here.
	In *string `protobuf:"bytes,4,opt,name=in" json:"in,omitempty"`
	// the message type that matches what the query returns.
	// This entity message will be converted to the output type
	// input/output messages on rpc calls will have their fields ignored if they
	// don't match this entity.
	// the generated query function will return this message type
	// if absent, this query returns nothing, and .
	// The query does not have to return a fully populated message,
	// but additional rows returned from the query that do NOT exist on
	// the out message will be ignored.
	Out                  *string  `protobuf:"bytes,5,opt,name=out" json:"out,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QLImpl) Reset()         { *m = QLImpl{} }
func (m *QLImpl) String() string { return proto.CompactTextString(m) }
func (*QLImpl) ProtoMessage()    {}
func (*QLImpl) Descriptor() ([]byte, []int) {
	return fileDescriptor_6479c8984e34e8f5, []int{1}
}

func (m *QLImpl) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QLImpl.Unmarshal(m, b)
}
func (m *QLImpl) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QLImpl.Marshal(b, m, deterministic)
}
func (m *QLImpl) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QLImpl.Merge(m, src)
}
func (m *QLImpl) XXX_Size() int {
	return xxx_messageInfo_QLImpl.Size(m)
}
func (m *QLImpl) XXX_DiscardUnknown() {
	xxx_messageInfo_QLImpl.DiscardUnknown(m)
}

var xxx_messageInfo_QLImpl proto.InternalMessageInfo

func (m *QLImpl) GetQuery() []string {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *QLImpl) GetPmStrategy() string {
	if m != nil && m.PmStrategy != nil {
		return *m.PmStrategy
	}
	return ""
}

func (m *QLImpl) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *QLImpl) GetIn() string {
	if m != nil && m.In != nil {
		return *m.In
	}
	return ""
}

func (m *QLImpl) GetOut() string {
	if m != nil && m.Out != nil {
		return *m.Out
	}
	return ""
}

type MOpts struct {
	// must match a name of a QLImpl query in the service.
	Query *string `protobuf:"bytes,1,req,name=query" json:"query,omitempty"`
	// the before function will be called before running any sql code for
	// every input data element and if the return will be not empty/nil and
	// the error will be nil the data returned by this function will be
	// returned by the function skipping the code execution
	Before *bool `protobuf:"varint,10,opt,name=before" json:"before,omitempty"`
	// the after function will be called after running any sql code for
	// every output data element, the return data of this function will be ignored
	After                *bool    `protobuf:"varint,11,opt,name=after" json:"after,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MOpts) Reset()         { *m = MOpts{} }
func (m *MOpts) String() string { return proto.CompactTextString(m) }
func (*MOpts) ProtoMessage()    {}
func (*MOpts) Descriptor() ([]byte, []int) {
	return fileDescriptor_6479c8984e34e8f5, []int{2}
}

func (m *MOpts) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MOpts.Unmarshal(m, b)
}
func (m *MOpts) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MOpts.Marshal(b, m, deterministic)
}
func (m *MOpts) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MOpts.Merge(m, src)
}
func (m *MOpts) XXX_Size() int {
	return xxx_messageInfo_MOpts.Size(m)
}
func (m *MOpts) XXX_DiscardUnknown() {
	xxx_messageInfo_MOpts.DiscardUnknown(m)
}

var xxx_messageInfo_MOpts proto.InternalMessageInfo

func (m *MOpts) GetQuery() string {
	if m != nil && m.Query != nil {
		return *m.Query
	}
	return ""
}

func (m *MOpts) GetBefore() bool {
	if m != nil && m.Before != nil {
		return *m.Before
	}
	return false
}

func (m *MOpts) GetAfter() bool {
	if m != nil && m.After != nil {
		return *m.After
	}
	return false
}

type TypeMapping struct {
	Types                []*TypeMapping_TypeDescriptor `protobuf:"bytes,1,rep,name=types" json:"types,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *TypeMapping) Reset()         { *m = TypeMapping{} }
func (m *TypeMapping) String() string { return proto.CompactTextString(m) }
func (*TypeMapping) ProtoMessage()    {}
func (*TypeMapping) Descriptor() ([]byte, []int) {
	return fileDescriptor_6479c8984e34e8f5, []int{3}
}

func (m *TypeMapping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TypeMapping.Unmarshal(m, b)
}
func (m *TypeMapping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TypeMapping.Marshal(b, m, deterministic)
}
func (m *TypeMapping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TypeMapping.Merge(m, src)
}
func (m *TypeMapping) XXX_Size() int {
	return xxx_messageInfo_TypeMapping.Size(m)
}
func (m *TypeMapping) XXX_DiscardUnknown() {
	xxx_messageInfo_TypeMapping.DiscardUnknown(m)
}

var xxx_messageInfo_TypeMapping proto.InternalMessageInfo

func (m *TypeMapping) GetTypes() []*TypeMapping_TypeDescriptor {
	if m != nil {
		return m.Types
	}
	return nil
}

type TypeMapping_TypeDescriptor struct {
	// if this is not setup the proto_type must be one of the built-in types
	ProtoTypeName *string `protobuf:"bytes,1,opt,name=proto_type_name,json=protoTypeName" json:"proto_type_name,omitempty"`
	// If proto_type_name is set, this need not be set.  If both this and proto_type_name
	// are set, this must be one of TYPE_ENUM, TYPE_MESSAGE
	// TYPE_GROUP is not supported
	ProtoType *descriptor.FieldDescriptorProto_Type `protobuf:"varint,2,opt,name=proto_type,json=protoType,enum=google.protobuf.FieldDescriptorProto_Type" json:"proto_type,omitempty"`
	// if proto_label is not setup we consider any option except LABAEL_REPEATED
	ProtoLabel           *descriptor.FieldDescriptorProto_Label `protobuf:"varint,3,opt,name=proto_label,json=protoLabel,enum=google.protobuf.FieldDescriptorProto_Label" json:"proto_label,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                               `json:"-"`
	XXX_unrecognized     []byte                                 `json:"-"`
	XXX_sizecache        int32                                  `json:"-"`
}

func (m *TypeMapping_TypeDescriptor) Reset()         { *m = TypeMapping_TypeDescriptor{} }
func (m *TypeMapping_TypeDescriptor) String() string { return proto.CompactTextString(m) }
func (*TypeMapping_TypeDescriptor) ProtoMessage()    {}
func (*TypeMapping_TypeDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor_6479c8984e34e8f5, []int{3, 0}
}

func (m *TypeMapping_TypeDescriptor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TypeMapping_TypeDescriptor.Unmarshal(m, b)
}
func (m *TypeMapping_TypeDescriptor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TypeMapping_TypeDescriptor.Marshal(b, m, deterministic)
}
func (m *TypeMapping_TypeDescriptor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TypeMapping_TypeDescriptor.Merge(m, src)
}
func (m *TypeMapping_TypeDescriptor) XXX_Size() int {
	return xxx_messageInfo_TypeMapping_TypeDescriptor.Size(m)
}
func (m *TypeMapping_TypeDescriptor) XXX_DiscardUnknown() {
	xxx_messageInfo_TypeMapping_TypeDescriptor.DiscardUnknown(m)
}

var xxx_messageInfo_TypeMapping_TypeDescriptor proto.InternalMessageInfo

func (m *TypeMapping_TypeDescriptor) GetProtoTypeName() string {
	if m != nil && m.ProtoTypeName != nil {
		return *m.ProtoTypeName
	}
	return ""
}

func (m *TypeMapping_TypeDescriptor) GetProtoType() descriptor.FieldDescriptorProto_Type {
	if m != nil && m.ProtoType != nil {
		return *m.ProtoType
	}
	return descriptor.FieldDescriptorProto_TYPE_DOUBLE
}

func (m *TypeMapping_TypeDescriptor) GetProtoLabel() descriptor.FieldDescriptorProto_Label {
	if m != nil && m.ProtoLabel != nil {
		return *m.ProtoLabel
	}
	return descriptor.FieldDescriptorProto_LABEL_OPTIONAL
}

var E_Pkg = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         560003,
	Name:          "persist.pkg",
	Tag:           "bytes,560003,opt,name=pkg",
	Filename:      "persist/options.proto",
}

var E_Opts = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: (*MOpts)(nil),
	Field:         560004,
	Name:          "persist.opts",
	Tag:           "bytes,560004,opt,name=opts",
	Filename:      "persist/options.proto",
}

var E_Ql = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.ServiceOptions)(nil),
	ExtensionType: (*QueryOpts)(nil),
	Field:         560000,
	Name:          "persist.ql",
	Tag:           "bytes,560000,opt,name=ql",
	Filename:      "persist/options.proto",
}

var E_Mapping = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.ServiceOptions)(nil),
	ExtensionType: (*TypeMapping)(nil),
	Field:         560001,
	Name:          "persist.mapping",
	Tag:           "bytes,560001,opt,name=mapping",
	Filename:      "persist/options.proto",
}

var E_ServiceType = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.ServiceOptions)(nil),
	ExtensionType: (*PersistenceOptions)(nil),
	Field:         560002,
	Name:          "persist.service_type",
	Tag:           "varint,560002,opt,name=service_type,enum=persist.PersistenceOptions",
	Filename:      "persist/options.proto",
}

var E_MappedField = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         560006,
	Name:          "persist.mapped_field",
	Tag:           "varint,560006,opt,name=mapped_field",
	Filename:      "persist/options.proto",
}

func init() {
	proto.RegisterEnum("persist.PersistenceOptions", PersistenceOptions_name, PersistenceOptions_value)
	proto.RegisterType((*QueryOpts)(nil), "persist.QueryOpts")
	proto.RegisterType((*QLImpl)(nil), "persist.QLImpl")
	proto.RegisterType((*MOpts)(nil), "persist.MOpts")
	proto.RegisterType((*TypeMapping)(nil), "persist.TypeMapping")
	proto.RegisterType((*TypeMapping_TypeDescriptor)(nil), "persist.TypeMapping.TypeDescriptor")
	proto.RegisterExtension(E_Pkg)
	proto.RegisterExtension(E_Opts)
	proto.RegisterExtension(E_Ql)
	proto.RegisterExtension(E_Mapping)
	proto.RegisterExtension(E_ServiceType)
	proto.RegisterExtension(E_MappedField)
}

func init() { proto.RegisterFile("persist/options.proto", fileDescriptor_6479c8984e34e8f5) }

var fileDescriptor_6479c8984e34e8f5 = []byte{
	// 581 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0xc5, 0x76, 0xd3, 0x34, 0xd7, 0x25, 0x8d, 0x46, 0x05, 0x59, 0xe5, 0xd1, 0x28, 0x48, 0x28,
	0x04, 0xd5, 0x41, 0x59, 0x54, 0x60, 0x56, 0x54, 0x2d, 0x52, 0x45, 0x92, 0x26, 0x13, 0x56, 0x6c,
	0x22, 0xc7, 0x99, 0xb8, 0x56, 0x6d, 0xcf, 0xc4, 0x1e, 0x23, 0x65, 0xc7, 0x4b, 0xfc, 0x02, 0x1b,
	0x7e, 0x8a, 0x3f, 0x42, 0x33, 0x63, 0x3b, 0x45, 0x89, 0x54, 0x56, 0x9e, 0x39, 0x73, 0xcf, 0xb9,
	0x8f, 0x73, 0x65, 0x78, 0xc0, 0x48, 0x92, 0x06, 0x29, 0xef, 0x52, 0xc6, 0x03, 0x1a, 0xa7, 0x36,
	0x4b, 0x28, 0xa7, 0xa8, 0x9a, 0xc3, 0x47, 0x4d, 0x9f, 0x52, 0x3f, 0x24, 0x5d, 0x09, 0xcf, 0xb2,
	0x45, 0x77, 0x4e, 0x52, 0x2f, 0x09, 0x18, 0xa7, 0x89, 0x0a, 0x6d, 0x9d, 0x42, 0x6d, 0x9c, 0x91,
	0x64, 0x75, 0xc5, 0x78, 0x8a, 0x5e, 0x40, 0x75, 0x99, 0x91, 0x24, 0x20, 0xa9, 0xa5, 0x35, 0x8d,
	0xb6, 0xd9, 0x3b, 0xb0, 0x73, 0x25, 0x7b, 0xdc, 0xbf, 0x8c, 0x58, 0x88, 0x8b, 0xf7, 0x56, 0x06,
	0xbb, 0x0a, 0x42, 0x87, 0x50, 0x11, 0xe0, 0x4a, 0x52, 0x6a, 0x58, 0x5d, 0xd0, 0x31, 0x98, 0x2c,
	0x9a, 0xa6, 0x3c, 0x71, 0x39, 0xf1, 0x57, 0x96, 0xde, 0xd4, 0xda, 0x35, 0x0c, 0x2c, 0x9a, 0xe4,
	0x08, 0x42, 0xb0, 0x13, 0xbb, 0x11, 0xb1, 0x8c, 0xa6, 0xde, 0xae, 0x61, 0x79, 0x46, 0x75, 0xd0,
	0x83, 0xd8, 0xda, 0x91, 0xb1, 0x7a, 0x10, 0xa3, 0x06, 0x18, 0x34, 0xe3, 0x56, 0x45, 0x02, 0xe2,
	0xd8, 0xfa, 0x00, 0x95, 0x81, 0x2c, 0xf5, 0x56, 0x56, 0x7d, 0x9d, 0xf5, 0x21, 0xec, 0xce, 0xc8,
	0x82, 0x26, 0xc4, 0x82, 0xa6, 0xd6, 0xde, 0xc3, 0xf9, 0x4d, 0x44, 0xbb, 0x0b, 0x4e, 0x12, 0xcb,
	0x94, 0xb0, 0xba, 0xb4, 0x7e, 0xeb, 0x60, 0x7e, 0x5c, 0x31, 0x32, 0x70, 0x19, 0x0b, 0x62, 0x1f,
	0xbd, 0x81, 0x0a, 0x5f, 0xb1, 0xb2, 0xf9, 0x67, 0x65, 0xf3, 0xb7, 0x82, 0xe4, 0xf9, 0xbc, 0x9c,
	0x22, 0x56, 0x8c, 0xa3, 0x3f, 0x1a, 0xd4, 0xff, 0x7d, 0x41, 0xcf, 0xe1, 0x40, 0x8e, 0x78, 0x2a,
	0x22, 0xa6, 0xb2, 0x57, 0x4d, 0x36, 0x72, 0x5f, 0xc2, 0x22, 0x7a, 0x28, 0x9a, 0xbe, 0x04, 0x58,
	0xc7, 0xc9, 0x41, 0xd5, 0x7b, 0x1d, 0x5b, 0x19, 0x67, 0x17, 0xc6, 0xd9, 0xef, 0x03, 0x12, 0xce,
	0xd7, 0xea, 0x23, 0x81, 0xcb, 0x5a, 0x70, 0xad, 0x94, 0x43, 0x7d, 0x30, 0x95, 0x54, 0xe8, 0xce,
	0x48, 0x68, 0x19, 0x52, 0xeb, 0xe5, 0xff, 0x69, 0xf5, 0x05, 0x05, 0xab, 0x52, 0xe4, 0xb9, 0xd3,
	0x01, 0x34, 0x52, 0x03, 0x20, 0xb1, 0x47, 0xae, 0xd4, 0x86, 0xa1, 0x2a, 0x18, 0x93, 0x71, 0xbf,
	0x71, 0x0f, 0x99, 0x50, 0x9d, 0x8c, 0xde, 0x0d, 0x87, 0x17, 0xb8, 0xa1, 0x39, 0xaf, 0xc0, 0x60,
	0x37, 0x3e, 0x7a, 0xbc, 0x25, 0x57, 0x58, 0x50, 0xad, 0xef, 0xbf, 0x5a, 0xca, 0x49, 0x76, 0xe3,
	0x3b, 0xe7, 0xb0, 0x43, 0x85, 0x91, 0x4f, 0x37, 0x28, 0x03, 0xc2, 0xaf, 0xe9, 0xbc, 0x20, 0xfd,
	0x90, 0x24, 0xb3, 0x57, 0x2f, 0xdd, 0x90, 0x0b, 0x80, 0x25, 0xdb, 0xb9, 0x00, 0x7d, 0x19, 0xa2,
	0xe3, 0x0d, 0x8d, 0x09, 0x49, 0x3e, 0x07, 0x65, 0xd1, 0xd6, 0x97, 0x5c, 0x04, 0xad, 0xf7, 0xb9,
	0x58, 0x7a, 0xac, 0x2f, 0x43, 0x67, 0x0c, 0xd5, 0x28, 0x5f, 0x82, 0x3b, 0xb5, 0xbe, 0xe6, 0x5a,
	0x87, 0xdb, 0xd6, 0x03, 0x17, 0x3a, 0x8e, 0x0b, 0xfb, 0xa9, 0x22, 0x4a, 0x63, 0xef, 0xd6, 0xfd,
	0x26, 0x75, 0xeb, 0xbd, 0x47, 0xa5, 0xee, 0xe6, 0xf4, 0xb1, 0x99, 0x6b, 0x8a, 0x94, 0xce, 0x19,
	0xec, 0x8b, 0x6c, 0x64, 0x3e, 0x5d, 0x08, 0x47, 0xd1, 0x93, 0xed, 0x4e, 0x17, 0x09, 0x7e, 0xca,
	0x04, 0x7b, 0xd8, 0x54, 0x24, 0xf9, 0x76, 0xf6, 0xfa, 0xd3, 0xa9, 0x1f, 0xf0, 0xeb, 0x6c, 0x66,
	0x7b, 0x34, 0xea, 0x72, 0x2f, 0xf6, 0x42, 0x9a, 0xcd, 0xd5, 0x0f, 0xc3, 0x3b, 0xf1, 0x49, 0x7c,
	0x52, 0xfc, 0x62, 0xf2, 0xef, 0xdb, 0xfc, 0xfb, 0x37, 0x00, 0x00, 0xff, 0xff, 0x3c, 0x1e, 0x6c,
	0x8e, 0x7c, 0x04, 0x00, 0x00,
}
