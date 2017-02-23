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

// Process the request
func (g *Generator) ProcessRequest() {
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

	logrus.WithField("All Structures", g.AllStructures).Debug("All structures found")
	for _, file := range g.OriginalRequest.ProtoFile {
		if !g.IsDependency(file) {
			outFile := g.files.NewOrGetFile(file)
			g.currentFile = file
			// implement file
			for _, str := range g.AllStructures.List {
				if str.File.GetName() == file.GetName() {
					logrus.WithField("Structure", str).Debug("Implementing structure")
					outFile.P(str.GetValueFunction())
					outFile.P(str.GetScanFunction())
				}
			}

			for _, service := range file.Service {
				g.currentService = service
				if IsServicePersistEnabled(service) {
					for _, method := range service.Method {
						data := GetMethodExtensionData(method)
						if data != nil {
							logrus.WithFields(logrus.Fields{
								"method":             method.GetName(),
								"Query":              data.GetQuery(),
								"Arguments":          data.GetArguments(),
								"Persistence Module": data.GetPersist(),
								"Variable Mapping":   data.GetMapping(),
							}).Debug("implementing method")
							if msg := g.AllStructures.GetEntry(method.GetInputType()); msg != nil {
								// we need to check if we are in the same file
								g.ImplementedStructures.AddStruct(msg)
							} else {
								logrus.Fatalf("Input type %s for method %s in file %s is missing!", method.GetInputType(), method.GetName(), file.GetName())
							}
							if msg := g.AllStructures.GetEntry(method.GetOutputType()); msg != nil {
								// we need to check if we are in the same file
								g.ImplementedStructures.AddStruct(msg)
							} else {
								logrus.Fatalf("Output type %s for method %s in file %s is missing!", method.GetOutputType(), method.GetName(), file.GetName())
							}

							// implement function body
							switch {
							// unary function
							case !method.GetClientStreaming() && !method.GetServerStreaming():
							// client streaming function
							case method.GetClientStreaming() && !method.GetServerStreaming():
							// server streaming function
							case !method.GetClientStreaming() && method.GetServerStreaming():
							// both streaming
							case method.GetClientStreaming() && method.GetServerStreaming():

							}
						}
					}
				}
			}

		}
	}
}
