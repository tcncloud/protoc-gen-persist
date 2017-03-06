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
	"github.com/Sirupsen/logrus"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	gen "github.com/golang/protobuf/protoc-gen-go/generator"
)

type GenericDescriptor interface {
	GetName() string
}

type Struct struct {
	Descriptor       GenericDescriptor
	Package          string
	GoPackage        string
	ParentDescriptor *Struct
	IsMessage        bool
	IsInnerType      bool
	IsUsedAsField    bool
	FileOptions      *desc.FileOptions
	EnumDesc         *desc.EnumDescriptorProto
	MsgDesc          *desc.DescriptorProto
}

func (s *Struct) GetGoName() string {
	if s.IsMessage {
		return gen.CamelCase(s.MsgDesc.GetName())
	} else {
		return gen.CamelCase(s.EnumDesc.GetName())
	}
}

func (s *Struct) GetProtoName() string {
	if s.ParentDescriptor == nil {
		return "." + s.Package + "." + s.Descriptor.GetName()
	} else {
		return s.ParentDescriptor.GetProtoName() + "." + s.Descriptor.GetName()
	}
}
func (s *Struct) ProcessFieldUsage(allStructures *StructList) {
	for _, str := range *allStructures {
		if str.IsMessage {
			// check if one of the message fileds uses our s Struct as type
			for _, field := range str.MsgDesc.GetField() {
				if (field.GetType() == desc.FieldDescriptorProto_TYPE_MESSAGE ||
					field.GetType() == desc.FieldDescriptorProto_TYPE_ENUM ||
					field.GetType() == desc.FieldDescriptorProto_TYPE_GROUP) && s.GetProtoName() == field.GetTypeName() {
					s.IsUsedAsField = true
				}
			}
		}
	}
}

type StructList []*Struct

func NewStructList() *StructList {
	return &StructList{}
}

func (s *StructList) GetStructByProtoName(name string) *Struct {
	for _, str := range *s {
		logrus.WithField("StructName", str.GetProtoName()).Debug("Checking with structure")
		if str.GetProtoName() == name {
			return str
		}
	}
	return nil
}

func (s *StructList) AddEnum(enum *desc.EnumDescriptorProto, parent *Struct, pkg string, opts *desc.FileOptions) *Struct {
	str := &Struct{
		IsMessage:        false,
		IsInnerType:      (parent != nil),
		Descriptor:       enum,
		ParentDescriptor: parent,
		Package:          pkg,
		FileOptions:      opts,
		MsgDesc:          nil,
		EnumDesc:         enum,
	}

	*s = append(*s, str)
	return str
}
func (s *StructList) AddMessage(message *desc.DescriptorProto, parent *Struct, pkg string, opts *desc.FileOptions) *Struct {
	str := &Struct{
		IsMessage:        true,
		IsInnerType:      (parent != nil),
		Descriptor:       message,
		ParentDescriptor: parent,
		Package:          pkg,
		FileOptions:      opts,
		MsgDesc:          message,
		EnumDesc:         nil,
	}

	*s = append(*s, str)
	return str
}
