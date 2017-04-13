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

type GeneratorStruct interface {
	Generate() string
}

type Generator struct {
	OriginalRequest *plugin_go.CodeGeneratorRequest
	AllStructures   *StructList // all structures present in the files
	Files           *FileList

	crtFile *descriptor.FileDescriptorProto
}

func NewGenerator(request *plugin_go.CodeGeneratorRequest) *Generator {
	ret := new(Generator)
	ret.OriginalRequest = request
	ret.AllStructures = NewStructList()
	ret.Files = NewFileList()
	return ret
}

func (g *Generator) GetResponse() (*plugin_go.CodeGeneratorResponse, error) {
	//logrus.WithField("structs", g.AllStructures).Debug("collected structures")
	ret := new(plugin_go.CodeGeneratorResponse)
	logrus.Debugf("going over %d files\n", len(*g.Files))
	for _, fileStruct := range *g.Files {
		// format file Content

		if !fileStruct.Dependency {
			fileContents, err := fileStruct.Generate()
			if err != nil {
				return nil, err
			}
			ret.File = append(ret.File, &plugin_go.CodeGeneratorResponse_File{
				Content: proto.String(string(FormatCode(fileStruct.GetFileName(), fileContents))),
				Name:    proto.String(fileStruct.GetFileName()),
			})
		}
	}
	//logrus.WithField("response", ret).Debug("result")
	return ret, nil
}

func (g *Generator) Process() error {
	logrus.Debug("processing the generator")
	for _, file := range g.OriginalRequest.ProtoFile {
		dep := func() bool {
			for _, fileName := range g.OriginalRequest.FileToGenerate {
				if fileName == file.GetName() {
					return false
				}
			}
			return true
		}()
		logrus.WithFields(logrus.Fields{
			"fileName": file.GetName(),
			"dependency?": dep,
		}).Debug("about to get or create this file")
		f := g.Files.GetOrCreateFile(file, g.AllStructures, dep)
		err := f.Process()
		if err != nil {
			return err
		}
	}

	for _, f := range *g.Files {
		f.ProcessImports()
	}
	return nil
}
