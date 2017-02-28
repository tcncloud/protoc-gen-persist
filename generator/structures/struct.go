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

package structures

import (
	"html/template"

	"github.com/Sirupsen/logrus"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var (
	valueTemplate *template.Template
	scanTemplate  *template.Template
)

const (
	valueFuncTemplate = `
func (s {{.GetGoName}}) Value() (driver.Value, error) {
{{if .IsMessage}}
	marshaler := jsonpb.Marshaler{}
	buf, err := marshaler.MarshalToString(&s)
	if err != nil {
		return driver.Value(""), fmt.Errorf("Can't serialize to json structure %+v", s)
	}
	return driver.Value(buf), nil
{{end}}
{{if .IsEnum}}
return driver.Value(int32(s)), nil
{{end}}
}`

	scanFuncTemplate = `
func (s *{{.GetGoName}}) Scan(src interface{}) error {
{{if .IsMessage}}
switch src.(type) {
case string:
	err := jsonpb.UnmarshalString(src.(string), s)
	return err
case []byte:
	err := jsonpb.UnmarshalString(string(src.([]byte)), s)
	return err
default:
	return fmt.Errorf("Unsupported type for %+v deserializing {{.GetGoName}}", src)
}
{{end}}
{{if .IsEnum}}
switch src.(type) {
case int32:
*s = src.({{.GetGoName}})
return nil
default:
return fmt.Errorf("can't convert %+v to {{.GetGoName}}",src)
}
{{end}}
}`
)

func init() {
	var err error
	valueTemplate, err = template.New("ValueFunc").Parse(valueFuncTemplate)
	if err != nil {
		logrus.Fatal("Error parsing value function template!")
	}
	scanTemplate, err = template.New("ScanFunc").Parse(scanFuncTemplate)
	if err != nil {
		logrus.Fatal("Error parsing scan function template!")
	}

}

type GenericDescriptor interface {
	GetName() string
}

type Struct struct {
	Descriptor       GenericDescriptor
	Package          string
	ParentDescriptor *Struct
	IsMessage        bool
	IsInnerType      bool
}

type StructList []*Struct

func NewStructList() *StructList {
	return &StructList{}
}

func (s *StructList) AddEnum(enum *desc.EnumDescriptorProto, parent *Struct, pkg string) *Struct {
	str := &Struct{
		IsMessage:        false,
		IsInnerType:      (parent != nil),
		Descriptor:       enum,
		ParentDescriptor: parent,
		Package:          pkg,
	}

	*s = append(*s, str)
	return str
}
func (s *StructList) AddMessage(message *desc.DescriptorProto, parent *Struct, pkg string) *Struct {
	str := &Struct{
		IsMessage:        true,
		IsInnerType:      (parent != nil),
		Descriptor:       message,
		ParentDescriptor: parent,
		Package:          pkg,
	}

	*s = append(*s, str)
	return str
}
