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
	"strconv"
	"strings"

	"github.com/Shrugs/fauxgaux"
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

type Method struct {
	Desc    *descriptor.MethodDescriptorProto
	Service *Service
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

func (m *Method) GetQueryParamString() string {
	if inputTypeStruct := m.GetInputTypeStruct(); inputTypeStruct != nil {
		if opt := m.GetMethodOption(); opt != nil {
			if opt.GetArguments() != nil {
				return "," + strings.Join(fauxgaux.Chain(opt.GetArguments()).Map(func(arg string) string {
					// TODO check if the type is a mapped type
					if fld := inputTypeStruct.GetFieldType(arg); fld != nil {
						if m.IsTypeMapped(fld) {
							return m.GetMappedObject(fld) + ".ToSql(" + "req." + _gen.CamelCase(arg) + ")"
						}
					}

					return "req." + _gen.CamelCase(arg)
				}).ConvertString(), ",")
			}
		}
	}
	return ""
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
	if str.File.GetOrigName() != m.Service.File.GetOrigName() {
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

// GetMappedType return mapped type for a proto name
func (m *Method) GetMappedType(typ *descriptor.FieldDescriptorProto) string {
	if mapping := m.GetTypeMapping(); mapping != nil {
		// if we have a mapping we are going to process it first
		for _, mapp := range mapping.Types {
			logrus.WithField("mapping", mapp).WithField("type", typ).Debug("checking mapping")
			if mapp.GetProtoType() == typ.GetType() &&
				mapp.GetProtoLabel() == typ.GetLabel() &&
				mapp.GetProtoTypeName() == typ.GetTypeName() {
				return "*" + m.Service.File.ImportList.GetImportPkgForPath(GetGoPath(mapp.GetGoPackage())) + "." + mapp.GetGoType()
			}
		}
	}
	switch typ.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		logrus.Fatalf("we currently don't support groups/oneof structures %s", typ.GetName())
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if structure := m.Service.AllStructs.GetStructByProtoName(typ.GetTypeName()); structure != nil {
			if imp := m.Service.File.ImportList.GetGoNameByStruct(structure); imp != nil {
				return "*" + imp.GoPackageName + "." + structure.GetGoName()
			} else {
				return "*" + structure.GetGoName()
			}
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

func (m *Method) GetFieldsWithLocalTypesFor(str *Struct) map[string]string {
	if str.IsMessage {
		ret := map[string]string{}
		// NOTE we don't process oneof fields
		// for _, mapping := range str.MsgDesc.GetOneofDecl() {}
		for _, mp := range str.MsgDesc.GetField() {
			// skip oneof fields
			if mp.OneofIndex == nil {
				ret[_gen.CamelCase(mp.GetName())] = m.GetMappedType(mp)
			}
		}
		return ret
	} else {
		return map[string]string{}
	}
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
	if opt := m.GetMethodOption(); opt != nil {
		return opt.GetPersist() == persist.PersistenceOptions_SQL
	}
	return false
}

// func (m *Method) IsMongo() bool {
// 	if opt := m.GetMethodOption(); opt != nil {
// 		return opt.GetPersist() == persist.PersistenceOptions_MONGO
// 	}
// 	return false
// }
//
// func (m *Method) IsSpanner() bool {
// 	if opt := m.GetMethodOption(); opt != nil {
// 		return opt.GetPersist() == persist.PersistenceOptions_SPANNER
// 	}
// 	return false
// }

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

func (m *Method) Process() {
	logrus.Debugf("Process method %s", m.GetName())
}

func (m *Method) ProcessImports() {
	if m.GetMethodOption() != nil {
		if m.GetMethodOption().GetMapping() != nil {
			for _, mapping := range m.GetMethodOption().GetMapping().GetTypes() {
				m.Service.File.ImportList.GetOrAddImport(GetGoPackage(mapping.GetGoPackage()), GetGoPath(mapping.GetGoPackage()))
			}
		}
	}
}

// -- Methods

type Methods []*Method

func (m *Methods) AddMethod(desc *descriptor.MethodDescriptorProto, service *Service) {
	*m = append(*m, &Method{Desc: desc, Service: service})
}

func (m *Methods) Process() {
	for _, meth := range *m {
		meth.Process()
	}
}
