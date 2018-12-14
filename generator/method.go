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
	"regexp"
	"strings"

	"bytes"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/sirupsen/logrus"
	"github.com/coltonmorris/protoc-gen-persist/generator/parser"
	"github.com/coltonmorris/protoc-gen-persist/persist"
)

type Method struct {
	Desc    *descriptor.MethodDescriptorProto
	Service *Service
	Query   parser.Query
	backend BackendStringer
}

func NewMethod(desc *descriptor.MethodDescriptorProto, srv *Service) (*Method, error) {
	meth := &Method{Desc: desc, Service: srv}
	if meth.Service.IsSpanner() {
		meth.backend = &SpannerStringer{method: meth}
	} else {
		meth.backend = &SqlStringer{method: meth}
	}
	return meth, nil
}

func (m *Method) String() string {
	if m.IsUnary() {
		return NewUnaryStringer(m, m.backend).String()
	} else if m.IsClientStreaming() {
		return NewClientStreamStringer(m, m.backend).String()
	} else if m.IsServerStreaming() {
		return NewServerStreamStringer(m, m.backend).String()
	} else {
		return NewBidiStreamStringer(m, m.backend).String()
	}
}
func (m *Method) IsSelect() bool {
	if m.Query != nil && m.Query.Type() == parser.SELECT_QUERY {
		return true
	}
	return false
}
func (m *Method) GetMethodOption() *persist.MOpts {
	if m.Desc.Options != nil && proto.HasExtension(m.Desc.Options, persist.E_Opts) {
		ext, err := proto.GetExtension(m.Desc.Options, persist.E_Opts)
		if err == nil {
			return ext.(*persist.MOpts)
		}
	}
	return nil
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

// GetQuery returns the doctored query string for this method
// by applying the pm strategy to the joined query string found
// on this methods service.
func (m *Method) GetQuery() (string, []TypeDesc, error) {
	mopt := m.GetMethodOption()
	query, err := m.Service.GetUndoctoredQueryByName(mopt.GetQuery())
	if err != nil {
		return "", nil, err
	}
	orig := strings.Join(query.GetQuery(), " ")

	pmStrat := query.GetPmStrategy()

	nextParamMarker := func() func(string) string {
		var count int
		return func(req string) string {
			if pmStrat == "$" {
				count++
				return fmt.Sprintf("$%d", count)
			} else if pmStrat == "?" {
				return "?"
			}
			return req
		}
	}()

	newQuery := ""
	r := regexp.MustCompile("@[a-zA-Z0-9_]*")
	potentialFieldNames := r.FindAllString(orig, -1)
	fieldsMap := m.GetTypeDescForFieldsInStructSnakeCase(m.GetInputTypeStruct())
	params := make([]TypeDesc, 0)

	for _, pf := range potentialFieldNames {
		start := strings.Index(orig, pf)
		stop := start + len(pf)
		// index into the map (removing the "@")
		td, exists := fieldsMap[pf[1:]]
		// eat up to the field name
		newQuery += orig[:start]
		if !exists { // it was just part of the query, not a field on input
			newQuery += pf
		} else { // it was a field, mark it
			newQuery += nextParamMarker(pf)
			params = append(params, td)
		}
		// remove the already written stuff
		orig = orig[stop:]
	}
	newQuery += orig

	return newQuery, params, nil
}

// helper method for getting a files package for stream calls
// if the service.pb.go and the persist.go are in different packages
// it will return the import prefix+.  of the package,  otherwise it returns
// the empty string
func (m *Method) GetFilePackage() string {
	if !m.Service.File.DifferentImpl() {
		return ""
	}

	imp := m.Service.File.ImportList.GetImportPkgForPath(m.Service.File.GetFullGoPackage())
	if imp != "__invalid__import__" {
		return imp + "."
	}
	return ""

}

func (m *Method) GetGoTypeName(typ string) string {
	return m.Service.File.GetGoTypeName(typ)
}
func (m *Method) GetGoTypeNameByFieldDesc(ty *descriptor.FieldDescriptorProto) string {
	return m.GetGoTypeName(ty.GetTypeName())
}

func (m *Method) GetInputType() string {
	return m.GetGoTypeName(m.Desc.GetInputType())
}

// returns the last element of the type.  So instead of test.ExampleTable,
// it returns ExampleTable
func (m *Method) GetOutputTypeMinusPackage() string {
	strs := strings.Split(m.GetOutputType(), ".")
	return strs[len(strs)-1]
}

// returns the last element of the type.  So instead of test.ExampleTable,
// it returns ExampleTable
func (m *Method) GetInputTypeMinusPackage() string {
	strs := strings.Split(m.GetInputType(), ".")
	return strs[len(strs)-1]
}

func (m *Method) GetOutputType() string {
	return m.GetGoTypeName(m.Desc.GetOutputType())
}

func (m *Method) DefaultMapping(typ *descriptor.FieldDescriptorProto) string {
	switch typ.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return "__unsupported__type__"
		//logrus.Fatalf("we currently don't support groups/oneof structures %s", typ.GetName())
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if ret := m.GetGoTypeNameByFieldDesc(typ); ret != "" {
			if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				return "[]*" + ret
			} else {
				return "*" + ret
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

func (m *Method) GetMapping(typ *descriptor.FieldDescriptorProto) *persist.TypeMapping_TypeDescriptor {
	if mapping := m.Service.GetTypeMapping(); mapping != nil {
		// if we have a mapping we are going to process it first
		for _, mapp := range mapping.Types {
			ptn := mapp.GetProtoTypeName()
			if mapp.GetProtoType() == typ.GetType() &&
				mapp.GetProtoLabel() == typ.GetLabel() &&
				((ptn == typ.GetTypeName()) || ("."+ptn == typ.GetTypeName())) {
				return mapp
			}
		}
	}
	return nil
}

type TypeDesc struct {
	Name      string // ex. StartTime
	ProtoName string // start_time
	// Is default Mapping  ex: string, []float64
	// or, if is a message type, then *pb.TestMessage  []*TestMessage
	GoName string
	// just the Name part of the type no [], or *
	GoTypeName      string
	OrigGoName      string // Timestamp
	Struct          *Struct
	Mapping         *persist.TypeMapping_TypeDescriptor
	IsMapped        bool
	IsRepeated      bool
	IsEnum          bool
	IsMessage       bool
	FieldDescriptor *descriptor.FieldDescriptorProto
	// spanner.GenericColumnValue, spanner.NullString, spanner.NullInt64
	// or if just the GoName
	SpannerType string
	// name used as field in the spanner.Null* types. ex: StringVal, NullInt64
	SpannerTypeFieldName string
	// if our spannerType != GoName, we need to convert our message
	NeedsSpannerConversion bool
}

func SpannerType(t TypeDesc) string {
	if t.IsMapped {
		return "spanner.GenericColumnValue"
	}
	switch t.GoName {
	case "string":
		return "spanner.NullString"
	case "[]string":
		return "[]spanner.NullString"
	case "int64":
		return "spanner.NullInt64"
	case "[]int64":
		return "[]spanner.NullInt64"
	case "bool":
		return "spanner.NullBool"
	case "[]bool":
		return "[]spanner.NullBool"
	case "float64":
		return "spanner.NullFloat64"
	case "[]float64":
		return "[]spanner.NullFloat64"
	}

	return t.GoName
}

func SpannerTypeFieldName(t TypeDesc) string {
	switch t.GoName {
	case "string", "[]string":
		return "StringVal"
	case "int64", "[]int64":
		return "Int64"
	case "float64", "[]float64":
		return "Float64"
	case "bool", "[]bool":
		return "Bool"
	}
	return ""
}

func (m *Method) GetTypeDescArrayForStruct(str *Struct) []TypeDesc {
	ret := make([]TypeDesc, 0)
	if str == nil || !str.IsMessage {
		return ret
	}
	for _, mp := range str.MsgDesc.GetField() {
		if mp.OneofIndex != nil {
			continue
		}
		typeDesc := TypeDesc{
			Name:   _gen.CamelCase(mp.GetName()),
			Struct: m.Service.AllStructs.GetStructByFieldDesc(mp),
			GoName: func() string {
				if m.GetMapping(mp) != nil {
					typName, _ := getGoNamesForTypeMapping(m.GetMapping(mp), m.Service.File)
					return typName
				}
				return m.DefaultMapping(mp)
			}(),
			ProtoName:       mp.GetName(),
			GoTypeName:      m.GetGoTypeNameByFieldDesc(mp),
			OrigGoName:      m.DefaultMapping(mp),
			Mapping:         m.GetMapping(mp),
			IsMapped:        (m.GetMapping(mp) != nil),
			FieldDescriptor: mp,
			IsRepeated:      (mp.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED),
			IsEnum: (mp.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM &&
				m.GetMapping(mp) == nil),
			IsMessage: (mp.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE &&
				m.GetMapping(mp) == nil),
		}
		//TODO refactor typeDesc into using a NewTypeDesc method
		typeDesc.SpannerType = SpannerType(typeDesc)
		typeDesc.SpannerTypeFieldName = SpannerTypeFieldName(typeDesc)
		typeDesc.NeedsSpannerConversion = (typeDesc.SpannerType != typeDesc.GoName)

		ret = append(ret, typeDesc)
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

func (m *Method) GetName() string {
	return m.Desc.GetName()
}

func (m *Method) IsSQL() bool {
	return m.Service.IsSQL() && m.GetMethodOption() != nil
}
func (m *Method) IsSpanner() bool {
	return m.Service.IsSpanner() && m.GetMethodOption() != nil
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

func (m *Method) GetBeforeHookName() string {
	return P(m.GetName(), "BeforeHook")
}
func (m *Method) GetAfterHookName() string {
	return P(m.GetName(), "AfterHook")
}

func (m *Method) Process() error {
	logrus.Debugf("Process method %s", m.GetName())
	if m.IsSpanner() {
		query, tds, err := m.GetQuery()
		if err != nil {
			return fmt.Errorf("error processing spanner query: %v", err)
		}
		reader := bytes.NewBufferString(query)
		// TODO remove
		p := parser.NewParser(reader)
		parsedQuery, err := p.Parse()
		if err != nil {
			return fmt.Errorf("%s\n  method: %s", err, m.GetName())
		}
		m.Query = parsedQuery
		for _, t := range tds {
			// this needs to be @, otherwise it will not be found
			m.Query.AddParam("@"+t.ProtoName, fmt.Sprintf("req.Get%s()", t.Name))
		}
	}
	return nil
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

func (m *Methods) PreGenerate() error {
	for _, meth := range *m {
		if err := meth.Process(); err != nil {
			return err
		}
	}
	return nil
}
