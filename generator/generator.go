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
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

type Generator struct {
	OriginalRequest *plugin_go.CodeGeneratorRequest
	// Response        *plugin_go.CodeGeneratorResponse

	// structures that need Value() and Scan() implementation
	ImplementedStructures *StructList
	AllStructures         *StructList

	currentFile    *descriptor.FileDescriptorProto    // current processing file
	currentService *descriptor.ServiceDescriptorProto // current processing service

	files *FileList
}

func NewGenerator(request *plugin_go.CodeGeneratorRequest) *Generator {
	ret := new(Generator)
	ret.OriginalRequest = request
	// ret.Response = new(plugin_go.CodeGeneratorResponse)
	ret.ImplementedStructures = NewStructList()
	ret.AllStructures = NewStructList()
	ret.files = NewFileList()
	return ret
}

func (g *Generator) GetResponse() *plugin_go.CodeGeneratorResponse {
	ret := new(plugin_go.CodeGeneratorResponse)
	for _, fileStruct := range g.files.List {
		// format file Content

		ret.File = append(ret.File, &plugin_go.CodeGeneratorResponse_File{
			Content: proto.String(fileStruct.GetContent()),
			Name:    proto.String(fileStruct.Name),
		})
		ret.Error = proto.String(func() string {
			if ret.Error == nil {
				return fileStruct.ErrorList
			} else {
				return *ret.Error + "\n" + fileStruct.ErrorList
			}
		}())
	}
	logrus.WithField("response", ret).Debug("result")
	return ret
}

func (g *Generator) IsDependency(file *descriptor.FileDescriptorProto) bool {
	for _, f := range g.OriginalRequest.ProtoFile {
		for _, name := range f.GetDependency() {
			if name == file.GetName() {
				return true
			}
		}
	}
	return false
}

func (g *Generator) ProcessType(typ string) {
	struc := g.AllStructures.GetEntry(typ)
	if struc != nil {
		// if the structure is not defined into a dependency file
		if !g.IsDependency(struc.File) {
			outFile := g.files.NewOrGetFile(struc.File)
			if struc.IsMessage() {
				// process inner enums
				for _, e := range struc.MessageDescriptor.GetEnumType() {
					if eStruct := g.AllStructures.GetEnum(e, struc.MessageDescriptor, struc.File); eStruct != nil {
						if !g.ImplementedStructures.ContainEnum(eStruct.EnumDescriptor, eStruct.ParentDescriptor, eStruct.File) {
							outFile.P(eStruct.GetValueFunction())
							outFile.P(eStruct.GetScanFunction())
							g.ImplementedStructures.AddStruct(eStruct)
						}
					}
				}
				// process inner messages
				for _, m := range struc.MessageDescriptor.GetNestedType() {
					if mStruct := g.AllStructures.GetMessage(m, struc.MessageDescriptor, struc.File); mStruct != nil {
						if !g.ImplementedStructures.ContainMessage(mStruct.MessageDescriptor, mStruct.ParentDescriptor, mStruct.File) {
							outFile.P(mStruct.GetValueFunction())
							outFile.P(mStruct.GetScanFunction())
							g.ImplementedStructures.AddStruct(mStruct)
						}
					}
				}
				// process fields
				for _, f := range struc.MessageDescriptor.GetField() {
					if f.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM ||
						f.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
						if t := g.AllStructures.GetEntry(f.GetTypeName()); t != nil {
							if g.ImplementedStructures.GetEntry(f.GetTypeName()) == nil {
								if !g.IsDependency(t.File) {
									out := g.files.NewOrGetFile(t.File)
									out.P(t.GetValueFunction())
									out.P(t.GetScanFunction())
									g.ImplementedStructures.AddStruct(t)
								}
							}
						}
					}
				}
			} else {
				// don't process structures that are not messages
			}
		}
	}
}

// process input and output types from service method signatures and implement Value() and Scan() mathods
// for messages and enums used inside of those types
func (g *Generator) ProcessStructs() {
	for _, file := range g.OriginalRequest.ProtoFile {
		if !g.IsDependency(file) {
			for _, service := range file.Service {
				for _, method := range service.GetMethod() {
					g.ProcessType(method.GetInputType())
					g.ProcessType(method.GetOutputType())
				}
			}
		}
	}
}

func (g *Generator) ProcessAllStructures() {
	for _, file := range g.OriginalRequest.ProtoFile {
		g.currentFile = file
		// scan all messages
		for _, m := range file.GetMessageType() {
			for _, im := range m.GetNestedType() {
				g.AllStructures.AddMessage(im, m, file)
			}
			for _, ie := range m.GetEnumType() {
				g.AllStructures.AddEnum(ie, m, file)
			}
			g.AllStructures.AddMessage(m, nil, file)

		}
		// scan all enums
		for _, e := range file.GetEnumType() {
			g.AllStructures.AddEnum(e, nil, file)
		}
	}
}

func (g *Generator) ProcessServices() {

	for _, file := range g.OriginalRequest.ProtoFile {
		if !g.IsDependency(file) {
			outFile := g.files.NewOrGetFile(file)

			for _, service := range file.Service {
				g.currentService = service
				if IsServicePersistEnabled(service) {
					srv := NewService(service, file, g.AllStructures, g.files)
					outFile.P(srv.Generate())
				}
			}
		}
	}
}

// check if a service has at least one method that has the persist.ql extension defined
func IsServicePersistEnabled(service *descriptor.ServiceDescriptorProto) bool {
	if service.Method != nil {
		for _, method := range service.Method {
			if IsMethodEnabled(method) {
				return true
			}
		}
	}
	return false
}

func IsMethodEnabled(method *descriptor.MethodDescriptorProto) bool {
	if method != nil && method.GetOptions() != nil && proto.HasExtension(method.Options, persist.E_Ql) {
		return true
	}
	return false
}

func GetMethodOption(method *descriptor.MethodDescriptorProto) *persist.QLImpl {
	if IsMethodEnabled(method) {
		if ret, err := proto.GetExtension(method.Options, persist.E_Ql); err == nil {
			return ret.(*persist.QLImpl)
		}
	}
	return nil
}

// Process the request
func (g *Generator) ProcessRequest() {
	g.ProcessAllStructures()
	g.ProcessStructs()
	g.ProcessServices()
}
