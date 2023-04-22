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

	desc "google.golang.org/protobuf/types/descriptorpb"
)

type GenericDescriptor interface {
	GetName() string
}

type Struct struct {
	Descriptor       GenericDescriptor
	Package          string
	ParentDescriptor *Struct
	IsMessage        bool
	IsInnerType      bool
	File             *FileStruct // for determine go import path and go package
	EnumDesc         *desc.EnumDescriptorProto
	MsgDesc          *desc.DescriptorProto
	ProtoName        string
	GoName           string

	// private fileds
	messageFieldDesc []*desc.FieldDescriptorProto
}

func (s *Struct) GetGoPath() string {
	if s == nil || s.File == nil || s.File.Desc == nil || s.File.Desc.Options == nil {
		return "__unknown__path__error__"
	}
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
	return s.GoName

}

func (s *Struct) GetProtoName() string {
	return s.ProtoName
	// if s.ParentDescriptor == nil {
	// 	return "." + s.File.Desc.GetPackage() + "." + s.Descriptor.GetName()
	// } else {
	// 	return s.ParentDescriptor.GetProtoName() + "." + s.Descriptor.GetName()
	// }
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

// GetFieldDescriptors returns a slice of FieldDescriptors that exist
// on this message.  If this is not a message, it returns empty slice, false
func (s *Struct) GetFieldDescriptorsIfMessage() ([]*desc.FieldDescriptorProto, bool) {
	// ret := make([]*desc.FieldDescriptorProto, 0)
	if s == nil || !s.IsMessage {
		return s.messageFieldDesc, false
	}

	return s.messageFieldDesc, true
}

type StructList []*Struct

func NewStructList() *StructList {
	return &StructList{}
}

func (s *StructList) GetStructByName(name string) *Struct {
	for _, str := range *s {
		if str.Descriptor != nil && str.Descriptor.GetName() == name {
			return str
		}
	}
	return nil
}
func compareProtoName(name string, protoname string) bool {
	if len(name) == 0 || len(protoname) == 0 {
		return false
	}
	return (name == protoname) || (name == protoname[1:] && protoname[0] == '.')
}

var cache map[string]*Struct = make(map[string]*Struct)

func (s *StructList) GetStructByProtoName(name string) *Struct {
	if str, ok := cache[name]; ok {
		return str
	}

	for _, str := range *s {
		if compareProtoName(name, str.GetProtoName()) {
			cache[name] = str
			return str
		} else if str.GetGoName() == name {
			cache[name] = str
			return str
		}
		// logrus.WithField("protoName", str.GetProtoName()).WithField("goName", str.GetGoName()).Tracef("NOT FOUND %s", name)
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
		messageFieldDesc: make([]*desc.FieldDescriptorProto, 0),
	}

	if str.ParentDescriptor == nil {
		str.ProtoName = "." + str.File.Desc.GetPackage() + "." + str.Descriptor.GetName()
	} else {
		str.ProtoName = str.ParentDescriptor.GetProtoName() + "." + str.Descriptor.GetName()
	}
	if str.IsMessage {
		if str.IsInnerType {
			str.GoName = str.ParentDescriptor.GetGoName() + "_" + GoCamelCase(str.MsgDesc.GetName())
		} else {
			str.GoName = GoCamelCase(str.MsgDesc.GetName())
		}
	} else {
		if str.IsInnerType {
			str.GoName = str.ParentDescriptor.GetGoName() + "_" + GoCamelCase(str.EnumDesc.GetName())
		} else {
			str.GoName = GoCamelCase(str.EnumDesc.GetName())
		}
	}

	*s = append(*s, str)
	str.File.ProcessImportsForType(str.GetGoName())

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
		messageFieldDesc: make([]*desc.FieldDescriptorProto, 0),
	}

	if str.ParentDescriptor == nil {
		str.ProtoName = "." + str.File.Desc.GetPackage() + "." + str.Descriptor.GetName()
	} else {
		str.ProtoName = str.ParentDescriptor.GetProtoName() + "." + str.Descriptor.GetName()
	}
	if str.IsMessage {
		if str.IsInnerType {
			str.GoName = str.ParentDescriptor.GetGoName() + "_" + GoCamelCase(str.MsgDesc.GetName())
		} else {
			str.GoName = GoCamelCase(str.MsgDesc.GetName())
		}
	} else {
		if str.IsInnerType {
			str.GoName = str.ParentDescriptor.GetGoName() + "_" + GoCamelCase(str.EnumDesc.GetName())
		} else {
			str.GoName = GoCamelCase(str.EnumDesc.GetName())
		}
	}
	*s = append(*s, str)
	for _, innerMessage := range message.GetNestedType() {
		s.AddMessage(innerMessage, str, pkg, file)
	}
	for _, innerEnum := range message.GetEnumType() {
		s.AddEnum(innerEnum, str, pkg, file)
	}
	str.File.ProcessImportsForType(str.GetGoName())

	for _, f := range str.MsgDesc.GetField() {
		if f.OneofIndex == nil {
			str.messageFieldDesc = append(str.messageFieldDesc, f)
		}
	}
	return str
}

func (s *StructList) Append(struc *Struct) {
	*s = append(*s, struc)
}

func (s *StructList) GetStructByFieldDesc(fld *desc.FieldDescriptorProto) *Struct {
	for _, str := range *s {
		if str.GetProtoName() == fld.GetName() {
			// logrus.Debugf("the struct name matches. Struct: %s  fld: %s", str.GetProtoName(), fld.GetName())
			return str
		}
	}
	return nil
}
