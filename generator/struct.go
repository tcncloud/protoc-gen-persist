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
	"text/template"

	"bytes"

	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
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
	valueTemplate, err = template.New("valueFunc").Parse(valueFuncTemplate)
	if err != nil {
		logrus.Fatal("Error parsing value function template!")
	}
	scanTemplate, err = template.New("scanFunc").Parse(scanFuncTemplate)
	if err != nil {
		logrus.Fatal("Error parsing scan function template!")
	}

}

type Struct struct {
	EnumDescriptor    *descriptor.EnumDescriptorProto
	MessageDescriptor *descriptor.DescriptorProto
	ParentDescriptor  *descriptor.DescriptorProto
	File              *descriptor.FileDescriptorProto
}

func (s *Struct) String() string {
	if s.EnumDescriptor == nil {
		if s.ParentDescriptor != nil {
			return "." + s.File.GetPackage() + "." + s.ParentDescriptor.GetName() + "." + s.MessageDescriptor.GetName()
		} else {
			return "." + s.File.GetPackage() + "." + s.MessageDescriptor.GetName()
		}

	} else if s.MessageDescriptor == nil {
		if s.ParentDescriptor != nil {
			return "." + s.File.GetPackage() + "." + s.ParentDescriptor.GetName() + "." + s.EnumDescriptor.GetName()
		} else {
			return "." + s.File.GetPackage() + "." + s.EnumDescriptor.GetName()
		}
	} else {
		return "FATAL ERROR"
	}

}

func (s *Struct) GetGoPath() string {
	if s.File.Options != nil && s.File.GetOptions().GoPackage != nil {
		return s.File.GetOptions().GetGoPackage()
	} else {
		return fmt.Sprintf("%s;%s",
			s.File.GetName()[0:strings.LastIndex(s.File.GetName(), "/")],
			strings.Replace(s.File.GetPackage(), ".", "_", -1))
	}
}

func (s *Struct) IsInnerType() bool {
	return (s.ParentDescriptor != nil)
}

func (s *Struct) GetValueFunction() string {
	var tpl bytes.Buffer
	if err := valueTemplate.Execute(&tpl, s); err != nil {
		logrus.Fatal("Error processing value function temaplate")
	}
	return tpl.String()
}

func (s *Struct) GetScanFunction() string {
	var tpl bytes.Buffer
	if err := scanTemplate.Execute(&tpl, s); err != nil {
		logrus.Fatal("Error processing scan function temaplate")
	}
	return tpl.String()
}

func (s *Struct) IsMessage() bool {
	return s.MessageDescriptor != nil && s.EnumDescriptor == nil
}

func (s *Struct) IsEnum() bool {
	return s.MessageDescriptor == nil && s.EnumDescriptor != nil
}

func (s *Struct) GetName() string {
	return s.String()
}

func (s *Struct) GetGoName() string {
	if s.ParentDescriptor != nil {
		if s.IsMessage() {

			return _gen.CamelCaseSlice([]string{s.ParentDescriptor.GetName(), s.MessageDescriptor.GetName()})
		} else {
			return _gen.CamelCaseSlice([]string{s.ParentDescriptor.GetName(), s.EnumDescriptor.GetName()})
		}
	} else {
		if s.IsMessage() {
			return _gen.CamelCase(s.MessageDescriptor.GetName())
		} else {
			return _gen.CamelCase(s.EnumDescriptor.GetName())
		}
	}
}

func (s *Struct) Equal(obj *Struct) bool {
	return ((s.ParentDescriptor == nil && obj.ParentDescriptor == nil) || (s.ParentDescriptor.GetName() == obj.ParentDescriptor.GetName())) &&
		((s.EnumDescriptor == nil && obj.EnumDescriptor == nil) || (s.EnumDescriptor.GetName() == obj.EnumDescriptor.GetName())) &&
		((s.MessageDescriptor == nil && obj.MessageDescriptor == nil) || (s.MessageDescriptor.GetName() == obj.MessageDescriptor.GetName())) &&
		(s.File.GetName() == obj.File.GetName())
}

type StructList struct {
	List []*Struct
}

func NewStructList() *StructList {
	return new(StructList)
}

func (sl *StructList) GetMessage(msg, parent *descriptor.DescriptorProto, file *descriptor.FileDescriptorProto) *Struct {
	for _, m := range sl.List {
		if m.IsMessage() {
			if (m.ParentDescriptor == nil && parent == nil) || (m.ParentDescriptor.GetName() == parent.GetName()) {
				if msg.GetName() == m.MessageDescriptor.GetName() &&
					file.GetPackage() == m.File.GetPackage() {
					return m
				}
			}
		}
	}
	return nil
}

func (sl *StructList) GetEnum(enum *descriptor.EnumDescriptorProto, parent *descriptor.DescriptorProto, file *descriptor.FileDescriptorProto) *Struct {
	for _, m := range sl.List {
		if m.IsEnum() {
			if (m.ParentDescriptor == nil && parent == nil) || (m.ParentDescriptor.GetName() == parent.GetName()) {
				if enum.GetName() == m.EnumDescriptor.GetName() &&
					file.GetPackage() == m.File.GetPackage() {
					return m
				}
			}
		}
	}
	return nil
}

func (sl *StructList) ContainMessage(msg, parent *descriptor.DescriptorProto, file *descriptor.FileDescriptorProto) bool {
	for _, m := range sl.List {
		if m.IsMessage() {
			if (m.ParentDescriptor == nil && parent == nil) || (m.ParentDescriptor.GetName() == parent.GetName()) {
				if msg.GetName() == m.MessageDescriptor.GetName() &&
					file.GetPackage() == m.File.GetPackage() {
					return true
				}
			}
		}
	}
	return false
}

func (sl *StructList) ContainEnum(enum *descriptor.EnumDescriptorProto, parent *descriptor.DescriptorProto, file *descriptor.FileDescriptorProto) bool {
	for _, m := range sl.List {
		if m.IsEnum() {
			if (m.ParentDescriptor == nil && parent == nil) || (m.ParentDescriptor.GetName() == parent.GetName()) {
				if enum.GetName() == m.EnumDescriptor.GetName() &&
					file.GetPackage() == m.File.GetPackage() {
					return true
				}
			}
		}
	}
	return false
}

func (sl *StructList) AddMessage(msg, parent *descriptor.DescriptorProto, file *descriptor.FileDescriptorProto) {
	if !sl.ContainMessage(msg, parent, file) {
		sl.List = append(sl.List, &Struct{
			MessageDescriptor: msg,
			ParentDescriptor:  parent,
			File:              file,
		})
	}
}

func (sl *StructList) AddEnum(enum *descriptor.EnumDescriptorProto, parent *descriptor.DescriptorProto, file *descriptor.FileDescriptorProto) {
	if !sl.ContainEnum(enum, parent, file) {
		sl.List = append(sl.List, &Struct{
			EnumDescriptor:   enum,
			ParentDescriptor: parent,
			File:             file,
		})
	}
}

func (sl *StructList) GetEntry(name string) *Struct {
	for _, entry := range sl.List {
		if entry.GetName() == name {
			return entry
		}
	}
	return nil
}

func (sl *StructList) AddStruct(s *Struct) {
	logrus.WithField("structure", s).Debug("adding structure")
	for _, x := range sl.List {
		if x.Equal(s) {
			return
		}
	}
	sl.List = append(sl.List, s)
	logrus.Debug("added!")
}
