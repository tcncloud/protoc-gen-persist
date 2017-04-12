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

package templates

const ServicesTemplate = `{{define "implement_services"}}
{{range $srv := .}}
{{if $srv.IsServiceEnabled}}
{{if $srv.IsSQL}}
type {{$srv.GetName}}Impl struct {
	SqlDB *sql.DB
}

func New{{$srv.GetName}}Impl(driver, connString string) (*{{$srv.GetName}}Impl, error) {
	db, err := sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}
	return &{{$srv.GetName}}Impl{ SqlDB: db }, nil
}
{{end}}
{{if $srv.IsSpanner}}
type {{$srv.GetName}}Impl struct {
	SpannerDB *spanner.Client
}

func New{{$srv.GetName}}Impl(d string, conf *spanner.ClientConfig, opts ...option.ClientOption) (*{{$srv.GetName}}Impl, error) {
	var client *spanner.Client
	var err error
	if conf != nil {
		client, err = spanner.NewClientWithConfig(context.Background(), d, *conf, opts...)
	} else {
		client, err = spanner.NewClient(context.Background(), d, opts...)
	}
	if err != nil {
		return nil, err
	}
	return &{{$srv.GetName}}Impl{SpannerDB: client}, nil
}
// need to implement rows
//{ {template "spanner_row_handler" .} }

{{end}}
{{range $method := $srv.Methods}}
{{template "implement_method" $method}}
{{end}}
{{end}}
{{end}}
{{end}}`
