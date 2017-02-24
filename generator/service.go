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
	"os"
	"text/template"

	"strings"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"

	_gen "github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

var (
	unaryTemplate        *template.Template
	serverStreamTemplate *template.Template
	clientStreamTemplate *template.Template
	bidirStreamTemplate  *template.Template
)

const (
	unary = `
func (s *{{.GetServiceImplName}}) {{.GetMethod}}(ctx context.Context, req *{{.GetInputType}}) (*{{.GetOutputType}}, error) {
	var (
		{{range $field := .GetSafeResponseFields}}
		{{$field.K}} {{$field.V}} {{end}}
	)
	err := s.DB.QueryRow(
		"{{.GetQuery}}",{{range $qParam := .GetQueryParams}}
		_utils.ToSafeType(req.{{$qParam}}),
		{{end}}).
		Scan({{range $field := .GetSafeResponseFields}} &{{$field.K}},
		{{end}})

	if err != nil {
		return nil, ConvertError(err, req)
	}
	result := &{{.GetOutputType}}{}
	{{range $local, $go := .GetResponseFieldsMap}}
	_utils.AssignTo(&result.{{$go}}, {{$local}}) {{end}}

	return result, nil
}
`
	server = `
func (s *{{.GetServiceImplName}}) {{.GetMethod}}(req *{{.GetInputType}}, stream {{.GetStreamType}}), error {
	rows, err := s.DB.Query("{{.GetQuery}}", {{range $qParam := .GetQueryParams}}
		ToSafeType(req.{{$qParam}}),
	{{end}})
	if err != nil {
		return ConvertError(err, req)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return ConvertError(err, req)
		}

		var (
			{{range $field := .GetSafeResponseFields}}
			{{$field.K}} {{$field.V}} {{end}}
		)

		err := rows.Scan({{range $field := .GetSafeResponseFields}}
			&{{$field.K}},{{end}}
		)
		if err != nil {
			return ConvertError(err, req)
		}

		result := &{{.GetOutputType}}{}
		{{ range $local, $go := .GetResponseFieldsMap}}
		AssignTo(&result.{{$go}}, {{$local}}) {{end}}
	}
	return result, nil
}
`
	client = `
func (s *{{.GetServiceImplName}}) {{.GetMethod}}(stream {{.GetStreamType}}), error {
	stmt, err := s.DB.Prepare("{{.GetQuery}}")
	if err != nil {
		return ConvertError(err, nil)
	}
	tx, err := s.db.Begin()
	if err != nil {
		return ConvertError(err, nil)
	}
	totalAffected := int64(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return ConvertError(err, req)
		}

		affected, err := tx.Stmt(stmt).Exec({{range $local := .GetQueryParams}}
		ToSafeType(req.{{$local}}),
		{{end}})
		if err != nil {
			tx.Rollback()
			return ConvertError(err, req)
		}
		num, err := affected.RowsAffected()
		if err != nil {
			tx.Rollback()
			return ConvertError(err, req)
		}
		totalAffected += num
	}
	err = tx.Commit()
	if err != nil {
		fmt.Errorf("Commiting transaction failed, rolling back...")
		return ConvertError(err, nil)
	}
	stream.SendAndClose(&{{.GetNumRowsMessageName}} { {{.GetNumRowsFieldName}}: totalAffected })
	return nil
}
`
	bidir = `
func (s *{{.GetServiceImplName}}) {{.GetMethod}}(stream {{.GetStreamType}}), error {
	stmt, err := s.DB.Prepare("{{.GetQuery}}")
	if err != nil {
		return ConvertError(err, nil)
	}

	defer stmt.Close()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return ConvertError(err, nil)
		}
		var( {{range $field := .GetSafeResponseFields}}
			{{$field.K}} {{$field.V}} {{end}}
		)

		err = stmt.QueryRow({{range $local := .GetQueryParams}}
			ToSafeType(req.{{$local}}),
		{{end}}).
		Scan({{range $key := .GetQueryParams}}
			&{{$key}},
		{{end}})

		result := &{{.GetOutputType}}{}

		{{range $local, $go := .GetResponseFieldsMap}}
		AssignTo(result.{{$go}}, {{$local}}){{end}}

		if err := stream.Send(result); err != nil {
			return ConvertError(err, req)
		}
	}
	return nil
}
`
)

func init() {
	var err error
	unaryTemplate, err = template.New("unaryTemplate").Parse(unary)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing unary template!")
	}
	serverStreamTemplate, err = template.New("serverTemplate").Parse(server)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing server stream template!")
	}

	clientStreamTemplate, err = template.New("clientTemplate").Parse(client)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing client stream template!")
	}

	bidirStreamTemplate, err = template.New("bidirTemplate").Parse(bidir)
	if err != nil {
		logrus.WithError(err).Fatal("Error parsing bidirectional stream template!")
	}
	//printTemplates()
}

func printTemplates() {
	logrus.Info("UNARY")
	unaryTemplate.Execute(os.Stdout, &Method{})
	fmt.Printf("\n")
	logrus.Info("SERVER")
	serverStreamTemplate.Execute(os.Stdout, &Method{})
	fmt.Printf("\n")
	logrus.Info("CLIENT")
	clientStreamTemplate.Execute(os.Stdout, &Method{})
	fmt.Printf("\n")
	logrus.Info("BIDI")
	bidirStreamTemplate.Execute(os.Stdout, &Method{})
	fmt.Printf("\n")
	logrus.Info("////////////////////////")
}

type Service struct {
	Desc              *descriptor.ServiceDescriptorProto
	File              *descriptor.FileDescriptorProto
	AllStruct         *StructList
	ImplementedStruct *StructList
	Files             *FileList
	Methods           []*Method
}

func NewService(service *descriptor.ServiceDescriptorProto,
	file *descriptor.FileDescriptorProto, allStruct *StructList, implStruct *StructList, files *FileList) *Service {
	s := &Service{
		Desc:              service,
		File:              file,
		AllStruct:         allStruct,
		ImplementedStruct: implStruct,
		Files:             files,
	}
	return s

}

func (s *Service) AddMethod(m *Method) {
	if s.Methods == nil {
		s.Methods = []*Method{
			m,
		}
	} else {
		s.Methods = append(s.Methods, m)
	}
}

func (s *Service) Generate() string {
	ret := fmt.Sprintf(`type %sImpl struct {
		DB *sql.DB
	}`, s.Desc.GetName())

	for _, method := range s.Desc.GetMethod() {
		if opt := GetMethodOption(method); opt != nil {
			m := NewMethod(method, s.Desc, opt, s.File, s.AllStruct, s.ImplementedStruct, s.Files)
			ret += m.Generate()
			s.AddMethod(m)
		}
	}
	return ret
}

type Method struct {
	Desc              *descriptor.MethodDescriptorProto
	Service           *descriptor.ServiceDescriptorProto
	Options           *persist.QLImpl
	File              *descriptor.FileDescriptorProto
	AllStruct         *StructList
	ImplementedStruct *StructList
	ImportList        map[string]string
	Files             *FileList
}

func NewMethod(method *descriptor.MethodDescriptorProto,
	service *descriptor.ServiceDescriptorProto,
	opt *persist.QLImpl,
	file *descriptor.FileDescriptorProto,
	allStruct *StructList,
	implStruct *StructList,
	files *FileList) *Method {

	m := &Method{
		Desc:              method,
		Service:           service,
		Options:           opt,
		File:              file,
		AllStruct:         allStruct,
		ImplementedStruct: implStruct,
		Files:             files,
	}
	m.ImportList = make(map[string]string)
	return m
}
func (m *Method) GetServiceImplName() string {
	return _gen.CamelCase(m.Service.GetName() + "Impl")
}

func (m *Method) GetMethod() string {
	return _gen.CamelCase(m.Desc.GetName())
}

func (m *Method) GetType(typ string) string {
	if struc := m.AllStruct.GetEntry(typ); struc != nil {
		if m.File.GetName() == struc.File.GetName() {
			// same package
			return struc.GetGoName()
		} else {
			// we have to determine the import path
			var name string
			if struc.File.GetOptions() != nil && struc.File.GetOptions().GoPackage != nil {
				if idx := strings.LastIndex(struc.File.GetOptions().GetGoPackage(), ";"); idx >= 0 {
					name = struc.File.GetOptions().GetGoPackage()[idx+1:]
					pkg := struc.File.GetOptions().GetGoPackage()[0:idx]
					name = m.Files.NewOrGetFile(m.File).AddImport(name, pkg, struc.File)
				} else if idx := strings.LastIndex(struc.File.GetOptions().GetGoPackage(), "/"); idx >= 0 {
					pkg := struc.File.GetOptions().GetGoPackage()[0:idx]
					name = struc.File.GetOptions().GetGoPackage()[idx+1:]
					name = m.Files.NewOrGetFile(m.File).AddImport(name, pkg, struc.File)
				} else {
					name = struc.File.GetOptions().GetGoPackage()
					pkg := struc.File.GetOptions().GetGoPackage()
					name = m.Files.NewOrGetFile(m.File).AddImport(name, pkg, struc.File)
				}
				return name + "." + struc.GetGoName()
			} else {
				// TODO add this to import paths
				return strings.Replace(struc.File.GetPackage(), ".", "_", -1) + "." + struc.GetGoName()
			}
		}

	}
	return ""
}

func (m *Method) GetInputType() string {
	return m.GetType(m.Desc.GetInputType())
}

func (m *Method) GetStreamType() string {
	return "Service_TestMethodServer"
}

func (m *Method) GetNumRowsMessageName() string {
	return "NumRows"
}

func (m *Method) GetNumRowsFieldName() string {
	return "Count"
}

func (m *Method) GetOutputType() string {
	return m.GetType(m.Desc.GetOutputType())
}

func (m *Method) GetServiceTypeMapping() *persist.TypeMapping {
	if proto.HasExtension(m.Service.Options, persist.E_Mapping) {
		e, err := proto.GetExtension(m.Service.Options, persist.E_Mapping)
		if err == nil {
			return e.(*persist.TypeMapping)
		}
	}
	return nil
}

func (m *Method) GetServiceFile() *FileStruct {
	return m.Files.GetFileByDesc(m.File)
}

func (m *Method) GetMethodTypeMapping() *persist.TypeMapping {
	return m.Options.Mapping
}

// GetUserSafeType is processing a field against the method or service defined options
// and process and register the necessary imports
func (m *Method) GetUserSafeType(field *descriptor.FieldDescriptorProto) string {
	logrus.WithField("field", field).Debug("checking field")
	list := func() *persist.TypeMapping {
		if mp := m.GetMethodTypeMapping(); mp != nil {
			return mp
		}
		if mp := m.GetServiceTypeMapping(); mp != nil {
			return mp
		}
		return nil
	}()
	logrus.WithField("list", list).Debug("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	if list != nil {
		for _, mp := range list.Types {
			logrus.WithField("mapping", mp).Debug("User defined mapping")
			// we have a message mapping where we map a message like
			// .google.protobuf.Timestamp to a go structure
			if mp.ProtoTypeName != nil &&
				(mp.ProtoType == nil ||
					mp.GetProtoType() == descriptor.FieldDescriptorProto_TYPE_ENUM ||
					mp.GetProtoType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE) {

				if field.GetTypeName() == mp.GetProtoTypeName() &&
					mp.GetProtoType() == field.GetType() &&
					((mp.ProtoLabel == nil && field.Label == nil) ||
						(mp.GetProtoLabel() == field.GetLabel())) {
					pkg, url := GetGoPackageAndPathFromURL(mp.GetGoPackage())
					// register package with the method file
					pkg = m.GetServiceFile().AddImport(pkg, url, m.File)
					// return go_type
					return pkg + "." + mp.GetGoType()
				}
			}
			// we map a builtin type int33, int64, string with a given
			// label (REPEATED mainly) to a db type
			if mp.ProtoTypeName == nil &&
				(mp.ProtoType != nil ||
					mp.GetProtoType() != descriptor.FieldDescriptorProto_TYPE_ENUM ||
					mp.GetProtoType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE ||
					mp.GetProtoType() != descriptor.FieldDescriptorProto_TYPE_GROUP) {
				if mp.GetProtoLabel() == field.GetLabel() &&
					mp.GetProtoType() == field.GetType() {
					pkg, url := GetGoPackageAndPathFromURL(mp.GetGoPackage())
					// register package with the method file
					pkg = m.GetServiceFile().AddImport(pkg, url, m.File)
					// return go_type
					return pkg + "." + mp.GetGoType()
				}
			}
		}
	}

	return ""
}

func (m *Method) GetSafeType(field *descriptor.FieldDescriptorProto) string {
	logrus.WithField("field", field).Debug("info")

	if ret := m.GetUserSafeType(field); ret != "" {
		return ret
	}

	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "float64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "float32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return "uint32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return "uint64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return "uint32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		logrus.Fatalf("Groups/Oneof are not supported yet")
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if stru := m.AllStruct.GetEntry(field.GetTypeName()); stru != nil {
			// find if is one of the structures that we've implemented Value and Scan Methods
			if s := m.ImplementedStruct.GetEntry(field.GetTypeName()); s != nil {
				if s.File.GetPackage() == m.File.GetPackage() {
					return field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]
				} else {
					pkg, url := GetGoPackageAndPathFromURL(s.GetGoPath())
					// register package with the method file
					pkg = m.GetServiceFile().AddImport(pkg, url, m.File)
					return pkg + "." + field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]
				}
			}
		} else {
			logrus.Fatalf("Can't find message structure for %+v", field)
		}
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		if stru := m.AllStruct.GetEntry(field.GetTypeName()); stru != nil {
			// find if is one of the structures that we've implemented Value and Scan Methods
			if s := m.ImplementedStruct.GetEntry(field.GetTypeName()); s != nil {
				if s.File.GetPackage() == m.File.GetPackage() {
					return field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]
				} else {
					pkg, url := GetGoPackageAndPathFromURL(s.GetGoPath())
					// register package with the method file
					pkg = m.GetServiceFile().AddImport(pkg, url, m.File)
					return pkg + "." + field.GetTypeName()[strings.LastIndex(field.GetTypeName(), ".")+1:]
				}
			}
		} else {
			logrus.Fatalf("Can't find enum structure for %+v", field)
		}
	// 	desc := g.ObjectNamed(field.GetTypeName())
	// 	typ, wire = g.TypeName(desc), "varint"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "int64"
	default:
		logrus.Fatalf("Unknown mapping for %+v", field)

	}
	return ""
}

type Tuple struct {
	K string
	V string
}

func (m *Method) GetSafeResponseFields() []*Tuple {
	var ret []*Tuple = nil
	outType := m.AllStruct.GetEntry(m.Desc.GetOutputType())
	for _, field := range outType.MessageDescriptor.Field {
		if ret == nil {
			ret = []*Tuple{
				&Tuple{
					K: field.GetName(),
					V: m.GetSafeType(field),
				},
			}
		} else {
			ret = append(ret, &Tuple{
				K: field.GetName(),
				V: m.GetSafeType(field),
			})
		}
	}
	return ret
}

func (m *Method) GetResponseFieldsMap() map[string]string {
	var ret map[string]string
	ret = make(map[string]string)
	outType := m.AllStruct.GetEntry(m.Desc.GetOutputType())
	for _, field := range outType.MessageDescriptor.Field {
		ret[field.GetName()] = _gen.CamelCase(field.GetName())
	}
	return ret
}

func (m *Method) GetQuery() string {
	return strings.Replace(strings.Replace(m.Options.GetQuery(), "\\", "\\\\", -1), "\"", "\\\"", -1)
}

func (m *Method) GetQueryParams() []string {
	return m.Options.Arguments
}

func (m *Method) Generate() string {
	var tpl bytes.Buffer
	if err := unaryTemplate.Execute(&tpl, m); err != nil {
		logrus.Fatal("Error processing value function temaplate")
	}
	return tpl.String()
}
