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
	"strconv"
	"text/template"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/tcncloud/protoc-gen-persist/generator/templates"
)

var (
	fileTemplate *template.Template
	TemplateList = map[string]string{
		"import_template":                 templates.ImportTemplate,
		"implement_structs":               templates.StructsTemplate,
		"implement_services":              templates.ServicesTemplate,
		"implement_method":                templates.MethodTemplate,
		"return_convert_helpers":          templates.ReturnConvertHelpers,
		"before_hook":                     templates.BeforeHook,
		"after_hook":                      templates.AfterHook,
		"unary_method":                    templates.UnaryMethodTemplate,
		"client_streaming_method":         templates.ClientStreamingMethodTemplate,
		"server_streaming_method":         templates.ServerStreamingMethodTemplate,
		"bidi_method":                     templates.BidiStreamingMethodTemplate,
		"sql_unary_method":                templates.SqlUnaryMethodTemplate,
		"sql_client_streaming_method":     templates.SqlClientStreamingMethodTemplate,
		"sql_server_streaming_method":     templates.SqlServerStreamingMethodTemplate,
		"sql_bidi_streaming_method":       templates.SqlBidiStreamingMethodTemplate,
		"mongo_unary_method":              templates.MongoUnaryMethodTemplate,
		"mongo_client_streaming_method":   templates.MongoClientStreamingMethodTemplate,
		"mongo_server_streaming_method":   templates.MongoServerStreamingMethodTemplate,
		"mongo_bidi_streaming_method":     templates.MongoBidiStreamingMethodTemplate,
		"spanner_unary_method":            templates.SpannerUnaryTemplate,
		"spanner_client_streaming_method": templates.SpannerClientStreamingTemplate,
		"spanner_server_streaming_method": templates.SpannerServerStreamingTemplate,
		"persist_lib_input":               templates.PersistLibInput,
	}
)

func init() {
	logrus.Debug("files package init()")
	var err error
	fileTemplate, err = template.New("fileTemplate").Parse(templates.MainTemplate)
	if err != nil {
		logrus.WithError(err).Fatal("Fail to parse file template")
	}
	fileTemplate := fileTemplate.Funcs(template.FuncMap{
		"Quotes": strconv.Quote,
	})

	for n, tmpl := range TemplateList {
		_, err = fileTemplate.Parse(tmpl)
		if err != nil {
			logrus.WithError(err).Fatalf("Fatal error parsing template %s", n)
		}
	}
}

func ExecuteFileTemplate(fileStruct *FileStruct) []byte {
	var buffer bytes.Buffer
	err := fileTemplate.Execute(&buffer, fileStruct)
	if err != nil {
		logrus.WithError(err).Fatal("Fatal error executing file template")
	}
	return buffer.Bytes()
}

func ExecutePersistLibTemplate(fileStruct *FileStruct) ([]byte, error) {
	var buffer bytes.Buffer
	t, err := template.New("t").Parse(templates.PersistLibTemplate)
	if err != nil {
		return nil, fmt.Errorf("could not parse the persistLibTemplate: %s", err)
	}
	t = t.Funcs(template.FuncMap{
		"Quotes": strconv.Quote,
	})
	for n, tmpl := range TemplateList {
		if _, err := t.Parse(tmpl); err != nil {
			logrus.WithError(err).Fatalf("Fatal error parsing template for persist lib: %s", n)
		}
	}

	if err := t.Execute(&buffer, fileStruct); err != nil {
		return nil, fmt.Errorf("could not execute persist lib template: %s", err)
	}
	return buffer.Bytes(), nil
}
