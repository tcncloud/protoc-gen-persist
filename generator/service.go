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
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

type Service struct {
	Desc       *descriptor.ServiceDescriptorProto
	Methods    *Methods
	Package    string // protobuf package
	File       *FileStruct
	AllStructs *StructList
}

func (s *Service) ProcessMethods() {
	for _, m := range s.Desc.GetMethod() {
		s.Methods.AddMethod(m, s)
	}
}

func (s *Service) Process() {
	s.ProcessMethods()
}

func (s *Service) GetName() string {
	return s.Desc.GetName()
}

func (s *Service) GetServiceOption() *persist.TypeMapping {
	if s.Desc.Options != nil && proto.HasExtension(s.Desc.Options, persist.E_Mapping) {
		ext, err := proto.GetExtension(s.Desc.Options, persist.E_Mapping)
		if err == nil {
			return ext.(*persist.TypeMapping)
		}
	}
	return nil
}

func (s *Service) IsSQL() bool {
	for _, m := range *s.Methods {
		if m.IsSQL() {
			return true
		}
	}
	return false
}

// func (s *Service) IsMongo() bool {
// 	for _, m := range *s.Methods {
// 		if m.IsMongo() {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func (s *Service) IsSpanner() bool {
// 	for _, m := range *s.Methods {
// 		if m.IsSpanner() {
// 			return true
// 		}
// 	}
// 	return false
// }

func (s *Service) IsServiceEnabled() bool {
	if s.GetServiceOption() != nil {
		return true
	}
	for _, m := range *s.Methods {
		if m.IsSQL() {
			return true
		}
	}
	return false
}

func (s *Service) ProcessImports() {
	s.File.ImportList.GetOrAddImport("io", "io")
	s.File.ImportList.GetOrAddImport("strings", "strings")
	s.File.ImportList.GetOrAddImport("context", "golang.org/x/net/context")
	s.File.ImportList.GetOrAddImport("grpc", "google.golang.org/grpc")
	s.File.ImportList.GetOrAddImport("codes", "google.golang.org/grpc/codes")
	if opt := s.GetServiceOption(); opt != nil {
		for _, m := range opt.GetTypes() {
			s.File.ImportList.GetOrAddImport(GetGoPackage(m.GetGoPackage()), GetGoPath(m.GetGoPackage()))
		}
	}
	for _, met := range *s.Methods {
		met.ProcessImports()
	}
}

type Services []*Service

func (s *Services) AddService(pkg string, desc *descriptor.ServiceDescriptorProto, allStructs *StructList, file *FileStruct) *Service {
	ret := &Service{
		Package:    pkg,
		Desc:       desc,
		Methods:    &Methods{},
		AllStructs: allStructs,
		File:       file,
	}
	ret.ProcessMethods()
	*s = append(*s, ret)
	return ret
}

func (s *Services) Process() {
	for _, srv := range *s {
		srv.Process()
	}
}