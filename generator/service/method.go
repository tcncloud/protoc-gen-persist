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

package service

import (
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/tcncloud/protoc-gen-persist/generator/structures"
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

func (m *Method) GetGoTypeNameWithPath(name string) string {
	if t := m.Service.AllStructs.GetStructByProtoName(name); t != nil {
		if t.Package == m.Service.Package {
			return t.GetGoName()
		} else {
			return t.GoPackage + "." + t.GetGoName()
		}
	} else {
		logrus.Fatalf("Can't find type %s ", m.Desc.GetInputType())
		return ""
	}
}

func (m *Method) GetInputType() string {
	return m.GetGoTypeNameWithPath(m.Desc.GetInputType())
}

func (m *Method) GetOutputType() string {
	return m.GetGoTypeNameWithPath(m.Desc.GetOutputType())
}

func (m *Method) GetServiceName() string {
	return m.Service.GetName()
}

func (m *Method) GetAllStructs() *structures.StructList {
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

func (m *Method) IsMongo() bool {
	if opt := m.GetMethodOption(); opt != nil {
		return opt.GetPersist() == persist.PersistenceOptions_MONGO
	}
	return false
}

func (m *Method) IsSpanner() bool {
	if opt := m.GetMethodOption(); opt != nil {
		return opt.GetPersist() == persist.PersistenceOptions_SPANNER
	}
	return false
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

type Methods []*Method

func (m *Methods) AddMethod(desc *descriptor.MethodDescriptorProto, service *Service) {
	*m = append(*m, &Method{Desc: desc, Service: service})
}
