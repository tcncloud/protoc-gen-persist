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

package files

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/tcncloud/protoc-gen-persist/generator/pkgimport"
	"github.com/tcncloud/protoc-gen-persist/generator/service"
	"github.com/tcncloud/protoc-gen-persist/generator/structures"
)

type FileStruct struct {
	Desc          *descriptor.FileDescriptorProto
	ImportList    *pkgimport.Imports
	Dependency    bool                   // if is dependency
	Structures    *structures.StructList // all structures in the file
	AllStructures *structures.StructList // all structures in all the files
	ServiceList   *service.Services
}

func NewFileStruct(desc *descriptor.FileDescriptorProto, allStructs *structures.StructList) *FileStruct {
	ret := &FileStruct{
		Desc:          desc,
		ImportList:    pkgimport.Empty(),
		Structures:    &structures.StructList{},
		ServiceList:   &service.Services{},
		AllStructures: allStructs,
	}
	return ret
}

func (f *FileStruct) GetOrigName() string {
	return f.Desc.GetName()
}

func (f *FileStruct) GetPackageName() string {
	return f.Desc.GetPackage()
}

func (f *FileStruct) GetFileName() string {
	return strings.Replace(f.Desc.GetName(), ".proto", ".persist.go", -1)
}

func (f *FileStruct) GetServices() *service.Services {
	return f.ServiceList
}

func (f *FileStruct) GetGoPackage() string {
	if f.Desc.Options != nil && f.Desc.GetOptions().GoPackage != nil {
		switch {
		case strings.Contains(f.Desc.GetOptions().GetGoPackage(), ";"):
			idx := strings.LastIndex(f.Desc.GetOptions().GetGoPackage(), ";")
			return f.Desc.GetOptions().GetGoPackage()[idx+1:]
		case strings.Contains(f.Desc.GetOptions().GetGoPackage(), "/"):
			idx := strings.LastIndex(f.Desc.GetOptions().GetGoPackage(), "/")
			return f.Desc.GetOptions().GetGoPackage()[idx+1:]
		default:
			return f.Desc.GetOptions().GetGoPackage()
		}

	} else {
		return strings.Replace(f.Desc.GetPackage(), ".", "_", -1)
	}
}

func (f *FileStruct) ProcessImports() {
	for _, srv := range *f.ServiceList {
		for _, m := range *srv.Methods {
			if s := f.AllStructures.GetStructByProtoName(m.Desc.GetInputType()); s != nil {
				if s.Package == f.Desc.GetPackage() {
					// same file
				} else {
					// s.GoPackage = f.ImportList.GetOrAddImport()
				}

			} else {
				logrus.Fatalf("Can't find structure %s", m.Desc.GetInputType())
			}
		}
	}

	for _, m := range f.Desc.GetMessageType() {
		logrus.WithField("message", m).Debug("message")
		f.Structures.AddMessage(m, nil, f.Desc.GetPackage(), f.Desc.GetOptions())
	}
	for _, e := range f.Desc.GetEnumType() {
		logrus.WithField("enum", e).Debug("enum")
		f.Structures.AddEnum(e, nil, f.Desc.GetPackage(), f.Desc.GetOptions())
	}
}

func (f *FileStruct) Generate() string {
	// f.Process()
	return ExecuteFileTemplate(f)
}

// FileList ----------------

type FileList []*FileStruct

func NewFileList() *FileList {
	return &FileList{}
}

func (fl *FileList) FindFile(desc *descriptor.FileDescriptorProto) *FileStruct {
	for _, f := range *fl {
		if f.Desc.GetName() == desc.GetName() {
			return f
		}
	}
	return nil
}

func (fl *FileList) GetOrCreateFile(desc *descriptor.FileDescriptorProto, allStructs *structures.StructList) *FileStruct {
	if f := fl.FindFile(desc); f != nil {
		return f
	}
	f := NewFileStruct(desc, allStructs)
	*fl = append(*fl, f)
	return f
}
