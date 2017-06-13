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
	"strings"

	"fmt"

	"github.com/Sirupsen/logrus"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	gen "github.com/golang/protobuf/protoc-gen-go/generator"
)

type GenericDescriptor interface {
	GetName() string
}

type Struct struct {
	Descriptor GenericDescriptor
	Package    string
	// GoPackage        string
	ParentDescriptor *Struct
	IsMessage        bool
	IsInnerType      bool
	IsUsedAsField    bool
	File             *FileStruct // for determine go import path and go package
	EnumDesc         *desc.EnumDescriptorProto
	MsgDesc          *desc.DescriptorProto
}

func (s *Struct) String() string {
	return fmt.Sprintf("[desc: %s, pkg: %s, file: %s, proto_pkg: %s, get_proto_name: %s]",
		s.Descriptor.GetName(),
		s.Package,
		s.File.Desc.GetName(),
		s.File.Desc.GetPackage(),
		s.GetProtoName(),
	)
}

func (s *Struct) GetGoPath() string {
	if s.File.Desc.Options != nil {
		if s.File.Desc.GetOptions().GoPackage != nil {
			pkg := s.File.Desc.GetOptions().GetGoPackage()
			if strings.Contains(pkg, ";") {
				idx := strings.LastIndex(pkg, ";")
				return pkg[0:idx]
			} else if strings.Contains(pkg, "/") {
				return pkg
			} else {
				return strings.Replace(pkg, ".", "_", -1)
			}
		} else {
			// return the package name
			return strings.Replace(s.Package, ".", "_", -1)
		}
	} else {
		return "__unknown__path__error__"
	}
}

func (s *Struct) GetFieldType(field string) *desc.FieldDescriptorProto {
	for _, f := range s.MsgDesc.Field {
		if f.GetName() == field {
			return f
		}
	}
	return nil
}
func (s *Struct) GetGoName() string {
	if s.IsMessage {
		if s.IsInnerType {
			return s.ParentDescriptor.GetGoName() + "_" + gen.CamelCase(s.MsgDesc.GetName())
		} else {
			return gen.CamelCase(s.MsgDesc.GetName())
		}
	} else {
		if s.IsInnerType {
			return s.ParentDescriptor.GetGoName() + "_" + gen.CamelCase(s.EnumDesc.GetName())
		} else {
			return gen.CamelCase(s.EnumDesc.GetName())
		}
	}
}

func (s *Struct) GetProtoName() string {
	if s.ParentDescriptor == nil {
		return "." + s.File.Desc.GetPackage() + "." + s.Descriptor.GetName()
	} else {
		return s.ParentDescriptor.GetProtoName() + "." + s.Descriptor.GetName()
	}
}

func (s *Struct) GetImportedFiles() *FileList {
	fl := NewFileList()
	fl.Append(s.File)
	if s.IsMessage {
		for _, field := range s.MsgDesc.GetField() {
			if str := s.File.AllStructures.GetStructByProtoName(field.GetName()); str != nil {
				fl.Append(str.File)
			}
		}
	}
	return fl
}

type StructList []*Struct

func NewStructList() *StructList {
	return &StructList{}
}

func (s *StructList) String() string {
	ret := ""
	for _, st := range *s {
		ret += st.String() + "\n"

	}
	return ret
}

func (s *StructList) GetStructByProtoName(name string) *Struct {
	for _, str := range *s {
		if str.GetProtoName() == name {
			return str
		}
	}
	return nil
}

func (s *StructList) AddEnum(enum *desc.EnumDescriptorProto, parent *Struct, pkg string, file *FileStruct) *Struct {
	str := &Struct{
		IsMessage:        false,
		IsInnerType:      (parent != nil),
		Descriptor:       enum,
		ParentDescriptor: parent,
		Package:          pkg,
		MsgDesc:          nil,
		EnumDesc:         enum,
		File:             file,
	}

	*s = append(*s, str)
	return str
}

func (s *StructList) AddMessage(message *desc.DescriptorProto, parent *Struct, pkg string, file *FileStruct) *Struct {
	str := &Struct{
		IsMessage:        true,
		IsInnerType:      (parent != nil),
		Descriptor:       message,
		ParentDescriptor: parent,
		Package:          pkg,
		MsgDesc:          message,
		EnumDesc:         nil,
		File:             file,
	}

	*s = append(*s, str)
	for _, innerMessage := range message.GetNestedType() {
		s.AddMessage(innerMessage, str, pkg, file)
	}
	for _, innerEnum := range message.GetEnumType() {
		s.AddEnum(innerEnum, str, pkg, file)
	}
	return str
}

func (s *StructList) Append(struc *Struct) {
	*s = append(*s, struc)
}

//TODO  str.GetProtoName  will never  be equal to fld.GetName
func (s *StructList) GetStructByFieldDesc(fld *desc.FieldDescriptorProto) *Struct {
	//logrus.Debugf("finding struct: %s", fld.GetName())
	for _, str := range *s {
		//logrus.Debugf("checking str:  %s", str.GetProtoName())
		if str.GetProtoName() == fld.GetName() {
			logrus.Debugf("the struct name matches. Struct: %s  fld: %s", str.GetProtoName, fld.GetName())
			return str
		}
	}
	return nil
}
