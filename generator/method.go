// Copyright 2017, TCN Inc.
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of TCN Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package generator

import (
	"fmt"
	"github.com/Shrugs/fauxgaux"
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/tcncloud/protoc-gen-persist/persist"
	"strconv"
	"strings"
)

type Method struct {
	Desc    *descriptor.MethodDescriptorProto
	Service *Service
	Spanner *SpannerHelper
}


func NewMethod(desc *descriptor.MethodDescriptorProto, srv *Service) (*Method, error) {
	meth := &Method{Desc: desc, Service: srv}
	return meth, nil
}

func (m *Method) String() string {
	if m == nil {
		return "METHOD: <nil>"
	}
	isSql := fmt.Sprintf("%t", m.IsSQL())
	isSpanner := fmt.Sprintf("%t", m.IsSpanner())
	name := m.Desc.GetName()
	input := m.Desc.GetInputType()
	output := m.Desc.GetOutputType()
	return fmt.Sprintf("Method:\n\tName: %s\n\tisSql: %s\n\tisSpanner: %s\n\tinput: %s\n\toutput: %s\n\tSpanner: %s\n\n",
		name, isSql, isSpanner, input, output, m.Spanner)
}

func (m *Method) GetMethodOption() *persist.QLImpl {
	if m.Desc.Options != nil && proto.HasExtension(m.Desc.Options, persist.E_Ql) {
		ext, err := proto.GetExtension(m.Desc.Options, persist.E_Ql)
		if err == nil {
			return ext.(*persist.QLImpl)
		}
	}
	return nil
}

func (m *Method) GetQueryParamString(comma bool) string {
	c := ""
	if comma {
		c = ","
	}
	if inputTypeStruct := m.GetInputTypeStruct(); inputTypeStruct != nil {
		if opt := m.GetMethodOption(); opt != nil {
			if opt.GetArguments() != nil {
				return c + strings.Join(fauxgaux.Chain(opt.GetArguments()).Map(func(arg string) string {
					// TODO check if the type is a mapped type
					if fld := inputTypeStruct.GetFieldType(arg); fld != nil {
						if m.IsTypeMapped(fld) {
							return m.GetMappedObject(fld) + "{}.ToSql(" + "req." + _gen.CamelCase(arg) + ")"
						}
					}

					return "req." + _gen.CamelCase(arg)
				}).ConvertString(), ",")
			}
		}
	}
	return ""
}

func (m *Method) GetFieldsWithLocalTypesFor(st *Struct) map[string]string {
	if st == nil {
		return nil
	}
	// The Fields on the struct
	mapping := make(map[string]string)
	//ranges over the proto fields
	for _, field := range st.MsgDesc.GetField() {
		// dont support oneof fields yet
		if field.Name != nil && field.OneofIndex == nil {
			name := _gen.CamelCase(*field.Name)
			if m.IsTypeMapped(field) {
				mapping[name] = m.GetMappedType(field)
			} else {
				mapping[name] = m.DefaultMapping(field)
			}
		}
	}
	return mapping
}

func (m *Method) GetTypeStructByProtoName(proto string) *Struct {
	return m.Service.AllStructs.GetStructByProtoName(proto)
}

func (m *Method) GetInputTypeStruct() *Struct {
	return m.GetTypeStructByProtoName(m.Desc.GetInputType())
}
func (m *Method) GetOutputTypeStruct() *Struct {
	return m.GetTypeStructByProtoName(m.Desc.GetOutputType())
}

func (m *Method) GetQuery() string {
	if opt := m.GetMethodOption(); opt != nil {
		return strconv.Quote(opt.GetQuery())
	}
	return ""
}

func (m *Method) GetGoTypeName(typ string) string {
	str := m.GetAllStructs().GetStructByProtoName(typ)
	if m.Service.File.GetPackageName() != str.File.GetPackageName() {
		if imp := m.Service.File.ImportList.GetGoNameByStruct(str); imp != nil {
			return imp.GoPackageName + "." + str.GetGoName()
		} else {
			logrus.WithField("struct", str).Fatal("Can't find struct in import list")
			return "__unknown__import__path__"
		}
	} else {
		return str.GetGoName()
	}
}

func (m *Method) GetInputType() string {
	return m.GetGoTypeName(m.Desc.GetInputType())
}

func (m *Method) GetOutputType() string {
	return m.GetGoTypeName(m.Desc.GetOutputType())
}

func (m *Method) GetTypeMapping() *persist.TypeMapping {
	if opt := m.GetMethodOption(); opt != nil {
		if opt.GetMapping() != nil {
			return opt.GetMapping()
		}
	}
	if opt := m.Service.GetServiceOption(); opt != nil {
		return opt
	}
	return nil

}

func (m *Method) IsTypeMapped(typ *descriptor.FieldDescriptorProto) bool {
	if mapping := m.GetTypeMapping(); mapping != nil {
		// if we have a mapping we are going to process it first
		for _, mapp := range mapping.Types {
			logrus.WithField("mapping", mapp).WithField("type", typ).Debug("checking mapping")
			if mapp.GetProtoType() == typ.GetType() &&
				mapp.GetProtoLabel() == typ.GetLabel() &&
				mapp.GetProtoTypeName() == typ.GetTypeName() {
				return true
			}
		}
	}
	return false
}

func (m *Method) GetMappedObject(typ *descriptor.FieldDescriptorProto) string {
	if mapping := m.GetTypeMapping(); mapping != nil {
		// if we have a mapping we are going to process it first
		for _, mapp := range mapping.Types {
			logrus.WithField("mapping", mapp).WithField("type", typ).Debug("checking mapping")
			if mapp.GetProtoType() == typ.GetType() &&
				mapp.GetProtoLabel() == typ.GetLabel() &&
				mapp.GetProtoTypeName() == typ.GetTypeName() {
				return m.Service.File.ImportList.GetImportPkgForPath(GetGoPath(mapp.GetGoPackage())) + "." + mapp.GetGoType()
			}
		}
	}
	return ""
}

func (m *Method) GetTypeNameMinusPackage(ty *descriptor.FieldDescriptorProto) string {
	if structure := m.Service.AllStructs.GetStructByProtoName(ty.GetTypeName()); structure != nil {
		if imp := m.Service.File.ImportList.GetGoNameByStruct(structure); imp != nil {
			return imp.GoPackageName + "." + structure.GetGoName()
		} else {
			return structure.GetGoName()
		}
	}
	return ""
}

func (m *Method) DefaultMapping(typ *descriptor.FieldDescriptorProto) string {
	switch typ.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		logrus.Fatalf("we currently don't support groups/oneof structures %s", typ.GetName())
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if ret := m.GetTypeNameMinusPackage(typ); ret != "" {
			return ret
		}
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]bool"
		} else {
			return "bool"
		}
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[][]byte"
		} else {
			return "[]byte"
		}
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]float64"
		} else {
			return "float64"
		}
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint32"
		} else {
			return "uint32"
		}
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint64"
		} else {
			return "uint64"
		}
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]float32"
		} else {
			return "float32"
		}
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int32"
		} else {
			return "int32"
		}
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int64"
		} else {
			return "int64"
		}
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int32"
		} else {
			return "int32"
		}
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int64"
		} else {
			return "int64"
		}
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int32"
		} else {
			return "int32"
		}
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int64"
		} else {
			return "int64"
		}
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]string"
		} else {
			return "string"
		}
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint32"
		} else {
			return "uint32"
		}
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint64"
		} else {
			return "uint64"
		}
	}
	return "__type__"
	//default mapping
}

// GetMappedType return mapped type for a proto name
func (m *Method) GetMappedType(typ *descriptor.FieldDescriptorProto) string {
	if mapping := m.GetTypeMapping(); mapping != nil {
		// if we have a mapping we are going to process it first
		for _, mapp := range mapping.Types {
			logrus.WithField("mapping", mapp).WithField("type", typ).Debug("checking mapping")
			if mapp.GetProtoType() == typ.GetType() &&
				mapp.GetProtoLabel() == typ.GetLabel() &&
				mapp.GetProtoTypeName() == typ.GetTypeName() {
				return m.Service.File.ImportList.GetImportPkgForPath(GetGoPath(mapp.GetGoPackage())) + "." + mapp.GetGoType()
			}
		}
	}
	logrus.Debug("returning default mapping")
	return m.DefaultMapping(typ)
}

func (m *Method) GetMapping(typ *descriptor.FieldDescriptorProto) *persist.TypeMapping_TypeDescriptor {
	if mapping := m.GetTypeMapping(); mapping != nil {
		// if we have a mapping we are going to process it first
		for _, mapp := range mapping.Types {
			if mapp.GetProtoType() == typ.GetType() &&
				mapp.GetProtoLabel() == typ.GetLabel() &&
				mapp.GetProtoTypeName() == typ.GetTypeName() {
				return mapp
			}
		}
	}
	return nil
}

type TypeDesc struct {
	Name       string // ex. StartTime
	ProtoName  string // start_time
	GoName     string // mytime.MyTime (if it is mapped)
	OrigGoName string // Timestamp
	Struct     *Struct
	Mapping    *persist.TypeMapping_TypeDescriptor
	EnumName   string // Timestamp
	IsMapped   bool
	IsEnum     bool
	IsMessage  bool
	ResultHook bool
}

type ResultHook interface {
	AddResult(req, row interface{}) error
}

func (m *Method) GetTypeDescArrayForStruct(str *Struct) []TypeDesc {
	ret := make([]TypeDesc, 0)
	if str != nil && str.IsMessage {
		for _, mp := range str.MsgDesc.GetField() {
			logrus.Debugf("mp name: %s\n", mp.GetName())
			if mp.OneofIndex == nil {
				typeDesc := TypeDesc{
					Name:       _gen.CamelCase(mp.GetName()),
					Struct:     m.Service.AllStructs.GetStructByFieldDesc(mp),
					ProtoName:  mp.GetName(),
					GoName:     m.GetMappedType(mp),
					OrigGoName: m.DefaultMapping(mp),
					Mapping:    m.GetMapping(mp),
					EnumName:   m.GetTypeNameMinusPackage(mp),
					IsMapped:   (m.GetMapping(mp) != nil),
					IsEnum:     (mp.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM),
					IsMessage:  (mp.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE && m.GetMapping(mp) == nil),
				}
				ret = append(ret, typeDesc)
			}
		}
	}
	return ret
}

func (m *Method) GetTypeDescForFieldsInStruct(str *Struct) map[string]TypeDesc {
	ret := map[string]TypeDesc{}
	for _, typeDesc := range m.GetTypeDescArrayForStruct(str) {
		ret[typeDesc.Name] = typeDesc
	}
	return ret
}

func (m *Method) GetTypeDescForFieldsInStructSnakeCase(str *Struct) map[string]TypeDesc {
	ret := map[string]TypeDesc{}
	for _, typeDesc := range m.GetTypeDescArrayForStruct(str) {
		ret[typeDesc.ProtoName] = typeDesc
	}
	return ret
}

func (m *Method) GetServiceName() string {
	return m.Service.GetName()
}

func (m *Method) GetAllStructs() *StructList {
	return m.Service.AllStructs
}

func (m *Method) GetName() string {
	return m.Desc.GetName()
}

func (m *Method) IsEnabled() bool {
	if m.GetMethodOption() != nil {
		return true
	}
	return false
}

func (m *Method) IsSQL() bool {
	return m.Service.IsSQL()
}

// func (m *Method) IsMongo() bool {
// 	if opt := m.GetMethodOption(); opt != nil {
// 		return opt.GetPersist() == persist.PersistenceOptions_MONGO
// 	}
// 	return false
// }

func (m *Method) IsSpanner() bool {
	return m.Service.IsSpanner()
}

func (m *Method) IsUnary() bool {
	return !m.Desc.GetClientStreaming() && !m.Desc.GetServerStreaming()
}

func (m *Method) IsClientStreaming() bool {
	return m.Desc.GetClientStreaming() && !m.Desc.GetServerStreaming()
}

func (m *Method) IsServerStreaming() bool {
	return !m.Desc.GetClientStreaming() && m.Desc.GetServerStreaming()
}

func (m *Method) IsBidiStreaming() bool {
	return m.Desc.GetClientStreaming() && m.Desc.GetServerStreaming()
}

func (m *Method) Process() error {
	logrus.Debug("Process method %s", m.GetName())
	if m.IsSpanner() {
		logrus.Debug("We are a spanner method")
		s, err := NewSpannerHelper(m)
		if err != nil {
			return err
		}
		m.Spanner = s
	} else if m.IsSQL() {
		logrus.Debug("we are a sql method")
	} else {
		logrus.Debug("we are neither?")
	}
	return nil
}

func (m *Method) ProcessImports() {
	if m.GetMethodOption() != nil {
		if m.GetMethodOption().GetMapping() != nil {
			for _, mapping := range m.GetMethodOption().GetMapping().GetTypes() {
				m.Service.File.ImportList.GetOrAddImport(GetGoPackage(mapping.GetGoPackage()), GetGoPath(mapping.GetGoPackage()))
			}
		}
		// if CallbackFunction options exist,  import the packages
		// name string, package string
		beforeOpt := m.GetMethodOption().GetBefore()
		afterOpt := m.GetMethodOption().GetAfter()
		if beforeOpt != nil {
			m.Service.File.ImportList.GetOrAddImport(GetGoPackage(beforeOpt.GetPackage()), GetGoPath(beforeOpt.GetPackage()))
		}
		if afterOpt != nil {
			m.Service.File.ImportList.GetOrAddImport(GetGoPackage(afterOpt.GetPackage()), GetGoPath(afterOpt.GetPackage()))
		}
	}
}

func (m *Method) GetGoPackage(path string) string {
	return GetGoPackage(path)
}

func (m *Method) GeGoPath(path string) string {
	return GetGoPath(path)
}

// -- Methods

type Methods []*Method

func (m *Methods) AddMethod(desc *descriptor.MethodDescriptorProto, service *Service) error {
	meth, err := NewMethod(desc, service)
	if err != nil {
		return err
	}
	*m = append(*m, meth)
	return nil
}

func (m *Methods) String() string {
	ret := "Methods:\n"
	for i, met := range *m {
		ret += fmt.Sprintf("\ti:%d val: %v", i, met)
	}
	return ret
}

func (m *Methods) PreGenerate() error {
	for _, meth := range *m {
		err := meth.Process()
		if err != nil {
			return err
		}
	}
	return nil
}
