// Code generated by protoc-gen-go. DO NOT EDIT.
// source: persist/options.proto

/*
Package persist is a generated protocol buffer package.

It is generated from these files:
	persist/options.proto

It has these top-level messages:
	QLImpl
	TypeMapping
*/
package persist

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PersistenceOptions int32

const (
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
func (PersistenceOptions) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type QLImpl struct {
	Query            []string     `protobuf:"bytes,1,rep,name=query" json:"query,omitempty"`
	Arguments        []string     `protobuf:"bytes,2,rep,name=arguments" json:"arguments,omitempty"`
	Mapping          *TypeMapping `protobuf:"bytes,4,opt,name=mapping" json:"mapping,omitempty"`
	Before           *bool        `protobuf:"varint,10,opt,name=before" json:"before,omitempty"`
	After            *bool        `protobuf:"varint,11,opt,name=after" json:"after,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *QLImpl) Reset()                    { *m = QLImpl{} }
func (m *QLImpl) String() string            { return proto.CompactTextString(m) }
func (*QLImpl) ProtoMessage()               {}
func (*QLImpl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *QLImpl) GetQuery() []string {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *QLImpl) GetArguments() []string {
	if m != nil {
		return m.Arguments
	}
	return nil
}

func (m *QLImpl) GetMapping() *TypeMapping {
	if m != nil {
		return m.Mapping
	}
	return nil
}

func (m *QLImpl) GetBefore() bool {
	if m != nil && m.Before != nil {
		return *m.Before
	}
	return false
}

func (m *QLImpl) GetAfter() bool {
	if m != nil && m.After != nil {
		return *m.After
	}
	return false
}

type TypeMapping struct {
	Types            []*TypeMapping_TypeDescriptor `protobuf:"bytes,1,rep,name=types" json:"types,omitempty"`
	XXX_unrecognized []byte                        `json:"-"`
}

func (m *TypeMapping) Reset()                    { *m = TypeMapping{} }
func (m *TypeMapping) String() string            { return proto.CompactTextString(m) }
func (*TypeMapping) ProtoMessage()               {}
func (*TypeMapping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TypeMapping) GetTypes() []*TypeMapping_TypeDescriptor {
	if m != nil {
		return m.Types
	}
	return nil
}

type TypeMapping_TypeDescriptor struct {
	ProtoTypeName    *string                                     `protobuf:"bytes,1,opt,name=proto_type_name,json=protoTypeName" json:"proto_type_name,omitempty"`
	ProtoType        *google_protobuf.FieldDescriptorProto_Type  `protobuf:"varint,2,opt,name=proto_type,json=protoType,enum=google.protobuf.FieldDescriptorProto_Type" json:"proto_type,omitempty"`
	ProtoLabel       *google_protobuf.FieldDescriptorProto_Label `protobuf:"varint,3,opt,name=proto_label,json=protoLabel,enum=google.protobuf.FieldDescriptorProto_Label" json:"proto_label,omitempty"`
	GoType           *string                                     `protobuf:"bytes,4,req,name=go_type,json=goType" json:"go_type,omitempty"`
	GoPackage        *string                                     `protobuf:"bytes,5,req,name=go_package,json=goPackage" json:"go_package,omitempty"`
	XXX_unrecognized []byte                                      `json:"-"`
}

func (m *TypeMapping_TypeDescriptor) Reset()                    { *m = TypeMapping_TypeDescriptor{} }
func (m *TypeMapping_TypeDescriptor) String() string            { return proto.CompactTextString(m) }
func (*TypeMapping_TypeDescriptor) ProtoMessage()               {}
func (*TypeMapping_TypeDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *TypeMapping_TypeDescriptor) GetProtoTypeName() string {
	if m != nil && m.ProtoTypeName != nil {
		return *m.ProtoTypeName
	}
	return ""
}

func (m *TypeMapping_TypeDescriptor) GetProtoType() google_protobuf.FieldDescriptorProto_Type {
	if m != nil && m.ProtoType != nil {
		return *m.ProtoType
	}
	return google_protobuf.FieldDescriptorProto_TYPE_DOUBLE
}

func (m *TypeMapping_TypeDescriptor) GetProtoLabel() google_protobuf.FieldDescriptorProto_Label {
	if m != nil && m.ProtoLabel != nil {
		return *m.ProtoLabel
	}
	return google_protobuf.FieldDescriptorProto_LABEL_OPTIONAL
}

func (m *TypeMapping_TypeDescriptor) GetGoType() string {
	if m != nil && m.GoType != nil {
		return *m.GoType
	}
	return ""
}

func (m *TypeMapping_TypeDescriptor) GetGoPackage() string {
	if m != nil && m.GoPackage != nil {
		return *m.GoPackage
	}
	return ""
}

var E_Pkg = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FileOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         560003,
	Name:          "persist.pkg",
	Tag:           "bytes,560003,opt,name=pkg",
	Filename:      "persist/options.proto",
}

var E_Ql = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*QLImpl)(nil),
	Field:         560000,
	Name:          "persist.ql",
	Tag:           "bytes,560000,opt,name=ql",
	Filename:      "persist/options.proto",
}

var E_Mapping = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.ServiceOptions)(nil),
	ExtensionType: (*TypeMapping)(nil),
	Field:         560001,
	Name:          "persist.mapping",
	Tag:           "bytes,560001,opt,name=mapping",
	Filename:      "persist/options.proto",
}

var E_ServiceType = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.ServiceOptions)(nil),
	ExtensionType: (*PersistenceOptions)(nil),
	Field:         560002,
	Name:          "persist.service_type",
	Tag:           "varint,560002,opt,name=service_type,json=serviceType,enum=persist.PersistenceOptions",
	Filename:      "persist/options.proto",
}

func init() {
	proto.RegisterType((*QLImpl)(nil), "persist.QLImpl")
	proto.RegisterType((*TypeMapping)(nil), "persist.TypeMapping")
	proto.RegisterType((*TypeMapping_TypeDescriptor)(nil), "persist.TypeMapping.TypeDescriptor")
	proto.RegisterEnum("persist.PersistenceOptions", PersistenceOptions_name, PersistenceOptions_value)
	proto.RegisterExtension(E_Pkg)
	proto.RegisterExtension(E_Ql)
	proto.RegisterExtension(E_Mapping)
	proto.RegisterExtension(E_ServiceType)
}

func init() { proto.RegisterFile("persist/options.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 514 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0x4e, 0x93, 0xe0, 0x31, 0xa4, 0xd5, 0xaa, 0xc0, 0xaa, 0x14, 0xb0, 0x82, 0x84, 0xac,
	0xa0, 0x3a, 0x28, 0x07, 0x04, 0xe1, 0x54, 0x04, 0x48, 0x95, 0xd2, 0x90, 0x6c, 0x38, 0x71, 0x89,
	0x1c, 0x67, 0xb2, 0xb5, 0x6a, 0x7b, 0x37, 0xb6, 0x83, 0x94, 0x1b, 0x1f, 0x27, 0x7e, 0x01, 0x12,
	0xbf, 0x8b, 0x1f, 0x84, 0xbc, 0xeb, 0x38, 0xa0, 0x54, 0xd0, 0xd3, 0xee, 0xbc, 0x99, 0xf7, 0x66,
	0xfc, 0x3c, 0x0b, 0x77, 0x24, 0xa6, 0x59, 0x98, 0xe5, 0x5d, 0x21, 0xf3, 0x50, 0x24, 0x99, 0x27,
	0x53, 0x91, 0x0b, 0xd2, 0x2c, 0xe1, 0x23, 0x87, 0x0b, 0xc1, 0x23, 0xec, 0x2a, 0x78, 0xb6, 0x5a,
	0x74, 0xe7, 0x98, 0x05, 0x69, 0x28, 0x73, 0x91, 0xea, 0xd2, 0xf6, 0x4f, 0x03, 0x1a, 0xe3, 0xc1,
	0x59, 0x2c, 0x23, 0x72, 0x08, 0xf5, 0xe5, 0x0a, 0xd3, 0x35, 0x35, 0x9c, 0x9a, 0x6b, 0x31, 0x1d,
	0x90, 0x63, 0xb0, 0xfc, 0x94, 0xaf, 0x62, 0x4c, 0xf2, 0x8c, 0x9a, 0x2a, 0xb3, 0x05, 0x88, 0x07,
	0xcd, 0xd8, 0x97, 0x32, 0x4c, 0x38, 0xdd, 0x73, 0x0c, 0xd7, 0xee, 0x1d, 0x7a, 0x65, 0x6f, 0xef,
	0xc3, 0x5a, 0xe2, 0xb9, 0xce, 0xb1, 0x4d, 0x11, 0xb9, 0x0b, 0x8d, 0x19, 0x2e, 0x44, 0x8a, 0x14,
	0x1c, 0xc3, 0xbd, 0xc9, 0xca, 0xa8, 0xe8, 0xed, 0x2f, 0x72, 0x4c, 0xa9, 0xad, 0x60, 0x1d, 0xb4,
	0x7f, 0x99, 0x60, 0xff, 0x21, 0x43, 0x5e, 0x42, 0x3d, 0x5f, 0x4b, 0xcc, 0xd4, 0x84, 0x76, 0xef,
	0xf1, 0x55, 0xbd, 0xd4, 0xfd, 0x4d, 0xf5, 0x99, 0x4c, 0x33, 0x8e, 0xbe, 0x9b, 0xd0, 0xfa, 0x3b,
	0x43, 0x9e, 0xc0, 0xbe, 0xf2, 0x60, 0x5a, 0x54, 0x4c, 0x13, 0x3f, 0x46, 0x6a, 0x38, 0x86, 0x6b,
	0xb1, 0xdb, 0x0a, 0x2e, 0xaa, 0x87, 0x7e, 0x8c, 0xe4, 0x0c, 0x60, 0x5b, 0x47, 0x4d, 0xc7, 0x70,
	0x5b, 0xbd, 0x8e, 0xa7, 0x9d, 0xf5, 0x36, 0xce, 0x7a, 0xef, 0x42, 0x8c, 0xe6, 0x5b, 0xf5, 0x51,
	0x81, 0xab, 0x59, 0x98, 0x55, 0xc9, 0x91, 0x01, 0xd8, 0x5a, 0x2a, 0xf2, 0x67, 0x18, 0xd1, 0x9a,
	0xd2, 0x7a, 0x7a, 0x3d, 0xad, 0x41, 0x41, 0x61, 0x7a, 0x14, 0x75, 0x27, 0xf7, 0xa0, 0xc9, 0xcb,
	0xa9, 0xf6, 0x1c, 0xd3, 0xb5, 0x58, 0x83, 0xeb, 0x36, 0x0f, 0x00, 0xb8, 0x98, 0x4a, 0x3f, 0xb8,
	0xf4, 0x39, 0xd2, 0xba, 0xca, 0x59, 0x5c, 0x8c, 0x34, 0xd0, 0xe9, 0x00, 0x19, 0x69, 0xe3, 0x30,
	0x09, 0xf0, 0xbd, 0x5e, 0x1d, 0xd2, 0x84, 0xda, 0x64, 0x3c, 0x38, 0xb8, 0x41, 0x6c, 0x68, 0x4e,
	0x46, 0xa7, 0xc3, 0xe1, 0x5b, 0x76, 0x60, 0xf4, 0x9f, 0x41, 0x4d, 0x5e, 0x72, 0x72, 0x7c, 0xc5,
	0x8c, 0xd1, 0x86, 0x4a, 0xbf, 0xfd, 0x68, 0x2b, 0xe3, 0x8a, 0xd2, 0xfe, 0x29, 0x98, 0xcb, 0x88,
	0x3c, 0xdc, 0x21, 0x9c, 0x63, 0x7e, 0x21, 0xe6, 0x1b, 0xca, 0x67, 0x45, 0xb1, 0x7b, 0xfb, 0xd5,
	0x3f, 0xd4, 0x5b, 0xc8, 0xcc, 0x65, 0xd4, 0x1f, 0x57, 0x5b, 0x45, 0x1e, 0xed, 0xe8, 0x4c, 0x30,
	0xfd, 0x14, 0x56, 0x63, 0xd3, 0x2f, 0xa5, 0xd0, 0xbf, 0x17, 0xaf, 0xef, 0xc3, 0xad, 0x4c, 0x13,
	0x95, 0x61, 0xff, 0xd7, 0xfd, 0xaa, 0x74, 0x5b, 0xbd, 0xfb, 0x95, 0xee, 0xae, 0x67, 0xcc, 0x2e,
	0x35, 0x8b, 0x96, 0xaf, 0x5f, 0x7c, 0x7c, 0xce, 0xc3, 0xfc, 0x62, 0x35, 0xf3, 0x02, 0x11, 0x77,
	0xf3, 0x20, 0x09, 0x22, 0xb1, 0x9a, 0xeb, 0xb7, 0x17, 0x9c, 0x70, 0x4c, 0x4e, 0x36, 0xaf, 0xb5,
	0x3c, 0x5f, 0x95, 0xe7, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb6, 0x66, 0x50, 0x87, 0xc7, 0x03,
	0x00, 0x00,
}
