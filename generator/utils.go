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
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/sirupsen/logrus"
	"github.com/tcncloud/protoc-gen-persist/persist"

	"golang.org/x/tools/imports"
)

var reduceEmptyLines = regexp.MustCompile("(\n)+")

// GetGoPath get a go import url under the following formats
// github.com/path/project/dir;package
// github.com/path/project/dir
// project/dir;package
// project/dir
// and will return the path portion from url:
// github.com/path/project/dir
// project/dir
func GetGoPath(url string) string {
	idx := strings.LastIndex(url, ";")
	switch {
	case idx >= 0:
		return url[0:idx]
	default:
		return url
	}
}

// GetGoPackage get a go import url under the following formats
// github.com/path/project/dir;package
// github.com/path/project/dir
// project/dir;package
// project/dir
// and will return the package name from url
// package
// dir
// package
// dir
func GetGoPackage(url string) string {
	switch {
	case strings.Contains(url, ";"):
		idx := strings.LastIndex(url, ";")
		return url[idx+1:]
	case strings.Contains(url, "/"):
		idx := strings.LastIndex(url, "/")
		return url[idx+1:]
	default:
		return url
	}
}

func FormatCode(filename string, buffer []byte) []byte {
	// reduce the empty lines
	tmp := reduceEmptyLines.ReplaceAll(buffer, []byte{'\n'})
	buf, err := imports.Process(filename, tmp, nil)
	if err != nil {
		logrus.WithError(err).Errorf("Error processing file %s", filename)
		return tmp
	}
	return buf
}
func getGoNamesForTypeMapping(tm *persist.TypeMapping_TypeDescriptor, file *FileStruct) (string, string) {
	name := file.GetGoTypeName(tm.GetProtoTypeName())
	nameParts := strings.Split(name, ".")
	for i, v := range nameParts {
		nameParts[i] = strings.Title(v)
	}
	titled := strings.Join(nameParts, "")
	return name, titled
}

func needsExtraStar(tm *persist.TypeMapping_TypeDescriptor) (bool, string) {
	if tm.GetProtoType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return true, "*"
	}
	return false, ""
}

func convertedMsgTypeByProtoName(protoName string, f *FileStruct) string {
	str := f.AllStructures.GetStructByProtoName(protoName)
	if str == nil {
		return ""
	}
	if imp := f.ImportList.GetGoNameByStruct(str); imp != nil {
		if f.NotSameAsMyPackage(imp.GoImportPath) {
			return imp.GoPackageName + "." + str.GetGoName()
		}
	}
	return str.GetGoName()
}
func defaultMapping(typ *descriptor.FieldDescriptorProto, file *FileStruct) (string, error) {
	switch typ.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return "int32", nil
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return "__unsupported__type__", fmt.Errorf("one of is unsupported")
		//logrus.Fatalf("we currently don't support groups/oneof structures %s", typ.GetName())
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:

		if ret := file.GetGoTypeName(typ.GetTypeName()); ret != "" {
			if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				return "[]*" + ret, nil
			} else {
				return "*" + ret, nil
			}
		}
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]bool", nil
		} else {
			return "bool", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[][]byte", nil
		} else {
			return "[]byte", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]float64", nil
		} else {
			return "float64", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint32", nil
		} else {
			return "uint32", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint64", nil
		} else {
			return "uint64", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]float32", nil
		} else {
			return "float32", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int32", nil
		} else {
			return "int32", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int64", nil
		} else {
			return "int64", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int32", nil
		} else {
			return "int32", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int64", nil
		} else {
			return "int64", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int32", nil
		} else {
			return "int32", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]int64", nil
		} else {
			return "int64", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]string", nil
		} else {
			return "string", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint32", nil
		} else {
			return "uint32", nil
		}
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		if typ.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return "[]uint64", nil
		} else {
			return "uint64", nil
		}
	}
	return "__type__", fmt.Errorf("unknown type")
}

type Printer struct {
	str string
}

func P(args ...interface{}) string {
	printer := &Printer{}
	printer.Q(args...)

	return printer.String()
}

func (p *Printer) P(formatString string, args ...interface{}) {
	p.str += fmt.Sprintf(formatString, args...)
}

func (p *Printer) Q(args ...interface{}) {
	for _, arg := range args {
		p.str += fmt.Sprintf("%v", arg)
	}
}
func (p *Printer) PA(formatStrings []string, args ...interface{}) {
	s := strings.Join(formatStrings, "")
	p.P(s, args...)
}

func (p *Printer) PTemplate(t string, dot interface{}) {
	var buff bytes.Buffer

	tem, err := template.New("printTemplate").Parse(t)
	if err != nil {
		p.P("\nPARSE ERROR:<%v>\nPARSING:<%s>\n", err, t)
		return
	}
	if err := tem.Execute(&buff, dot); err != nil {
		p.P("\nEXEC ERROR:<%v>\nEXECUTING:<%s>\n", err, t)
		return
	}
	p.P("%s", buff.String())
}

func (p *Printer) String() string {
	return p.str
}
