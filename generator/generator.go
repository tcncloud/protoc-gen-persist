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
	"github.com/tcncloud/protoc-gen-persist/generator/files"
	"github.com/tcncloud/protoc-gen-persist/generator/structures"
)

type GeneratorStruct interface {
	Generate() string
}

type Generator struct {
	OriginalRequest *plugin_go.CodeGeneratorRequest
	AllStructures   *structures.StructList // all structures present in the files
	Files           *files.FileList

	crtFile *descriptor.FileDescriptorProto
}

func NewGenerator(request *plugin_go.CodeGeneratorRequest) *Generator {
	ret := new(Generator)
	ret.OriginalRequest = request
	ret.AllStructures = structures.NewStructList()
	ret.Files = files.NewFileList()
	return ret
}

func (g *Generator) GetResponse() *plugin_go.CodeGeneratorResponse {
	ret := new(plugin_go.CodeGeneratorResponse)
	for _, fileStruct := range *g.Files {
		// format file Content

		ret.File = append(ret.File, &plugin_go.CodeGeneratorResponse_File{
			Content: proto.String(fileStruct.Generate()),
			Name:    proto.String(fileStruct.GetFileName()),
		})
		// ret.Error = proto.String(func() string {
		// 	if ret.Error == nil {
		// 		return fileStruct.ErrorList
		// 	} else {
		// 		return *ret.Error + "\n" + fileStruct.ErrorList
		// 	}
		// }())
	}
	logrus.WithField("response", ret).Debug("result")
	return ret
}

// Process the request
func (g *Generator) ProcessRequest() {
	// create a list of structures
	// g.ProcessStructs()
	// g.ProcessServices()
	// create structures
	for _, f := range g.OriginalRequest.ProtoFile {
		logrus.WithField("f", f.GetName()).Debug("Processing file")
		g.crtFile = f
		// determine if this file is just imported
		dependency := func() bool {
			for _, fileName := range g.OriginalRequest.FileToGenerate {
				if fileName == f.GetName() {
					return true
				}
			}
			return false
		}()

		if dependency {
			file := g.Files.GetOrCreateFile(f, g.AllStructures)
			logrus.WithField("new file name", file.GetFileName()).Debug("new file name")
			file.Process()
		}
		for _, m := range f.GetMessageType() {
			g.ProcessMessage(m, nil, f.GetOptions())
		}
		for _, e := range f.GetEnumType() {
			g.ProcessEnum(e, nil, f.GetOptions())
		}
	}
	for _, x := range *g.AllStructures {
		x.ProcessFieldUsage(g.AllStructures)
		logrus.Debugf("%s inner %b", x.GetProtoName(), x.IsInnerType)
	}
	for _, file := range *g.Files {
		file.ProcessImports()
	}
}

func (g *Generator) ProcessMessage(msg *descriptor.DescriptorProto, parent *structures.Struct, opts *descriptor.FileOptions) {
	// add the current message to the list
	logrus.WithField("message", msg.GetName()).Debug("processing message")
	m := g.AllStructures.AddMessage(msg, parent, g.crtFile.GetPackage(), opts)
	for _, message := range msg.GetNestedType() {
		g.ProcessMessage(message, m, opts)
	}
	for _, enum := range msg.GetEnumType() {
		g.ProcessEnum(enum, m, opts)
	}
}
func (g *Generator) ProcessEnum(enum *descriptor.EnumDescriptorProto, parent *structures.Struct, opts *descriptor.FileOptions) {
	// add the current message to the list
	g.AllStructures.AddEnum(enum, parent, g.crtFile.GetPackage(), opts)
}
