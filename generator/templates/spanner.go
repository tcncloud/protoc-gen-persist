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

const PersistLibInput = `{{define "persist_lib_input"}}{{$method := . -}}
type {{template "persist_lib_input_name" $method}} struct{
{{range $key, $val := $method.GetTypeDescForQueryFields -}}
	{{$val.Name}} {{if $val.IsMapped -}} interface{} {{else}} {{$val.GoName}}{{end}}
{{end}}
}
{{end}}

{{define "persist_lib_input_name"}}{{$method := . -}}
{{$method.Service.GetName}}{{$method.GetName}}Input
{{- end}}

{{define "persist_lib_method_input_name"}}{{$method := . -}}
{{$method.GetInputTypeName}}For{{$method.Desc.GetName}}
{{- end}}

{{define "persist_lib_default_handler"}}{{$method := . -}}
	{{if $method.IsClientStreaming -}}
		func Default{{$method.GetName}}Handler(cli *spanner.Client) func(context.Context) (func(*{{template "persist_lib_input_name" $method}}), func() (*spanner.Row, error)) {
			return func(ctx context.Context) (feed func(*{{template "persist_lib_input_name" $method}}), done func() (*spanner.Row, error)) {
				var muts []*spanner.Mutation
				feed = func(req *{{template "persist_lib_input_name" $method}}) {
					muts = append(muts, {{template "persist_lib_method_input_name" $method}}(req))
				}
				done = func() (*spanner.Row, error) {
					if _, err := cli.Apply(ctx, muts); err != nil {
						return nil, err
					}
					return nil, nil // we dont have a row, because we are an apply
				}
				return feed, done
			}
		}
	{{else}}
		func Default{{$method.GetName}}Handler(cli *spanner.Client) func(context.Context, *{{template "persist_lib_input_name" $method}}, func(*spanner.Row)) error {
			return func(ctx context.Context, req *{{template "persist_lib_input_name" $method}}, next func(*spanner.Row)) error {
				{{- if $method.IsSelect}}
				iter := cli.Single().Query(ctx, {{template "persist_lib_method_input_name" $method}}(req))
				if err := iter.Do(func(r *spanner.Row) error {
					next(r)
					return nil
				}); err != nil {
					return err
				}
				{{- else}}
				if _, err := cli.Apply(ctx, []*spanner.Mutation{ {{template "persist_lib_method_input_name" $method}}(req)}); err != nil {
					return err
				}
				next(nil) // this is an apply, it has no result

				{{- end}}

				return nil
			}
		}
	{{end}}
{{end}}
`

const SpannerUnaryTemplate = `{{define "spanner_unary_method" -}}
func (s* {{.GetServiceName}}Impl) {{.GetName}} (ctx context.Context, req *{{.GetInputType}}) (*{{.GetOutputType}}, error) {
	var err error
	_ = err
	{{- $method := . -}}

	{{template "before_hook" .}}

	params := &persist_lib.{{template "persist_lib_input_name" $method}}{}
	{{range $key, $val := .GetTypeDescForQueryFields -}}
		{{if $val.IsMapped -}}
			if params.{{$val.Name}}, err = ({{$val.GoName}}{}).ToSpanner(req.{{.Name}}).SpannerValue(); err != nil {
				return nil, gstatus.Errorf(codes.Unknown, "could not convert type: %v", err)
			}
		{{else -}}
			params.{{$val.Name}} =  req.{{$val.Name}}
		{{end -}}
	{{end}}

	var res = {{.GetOutputType}}{}
	var iterErr error
	err = s.PERSIST.{{.GetName}}(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		{{range $index, $field := .GetTypeDescArrayForStruct .GetOutputTypeStruct -}}
			{{if $field.IsMapped -}}
				var {{$field.Name}} *spanner.GenericColumnValue
				if err := row.ColumnByName("{{$field.ProtoName}}", {{$field.Name}}); err != nil {
					iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{else -}}
				var {{$field.Name}} {{$method.DefaultMapping $field.FieldDescriptor}}
				if err := row.ColumnByName("{{$field.ProtoName}}", &{{$field.Name}}); err != nil {
					iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{end}}
			{{if $field.IsMapped -}}
				{
					local := &{{$field.GoName}}{}
					if err := local.SpannerScan({{$field.Name}}); err != nil {
						iterErr = gstatus.Errorf(codes.Unknown, "could not scan out custom type: %s", err)
						return
					}
					res.{{$field.Name}} = local.ToProto()
				}
			{{else -}}
				res.{{$field.Name}} = {{$field.Name}}
			{{end -}}
		{{end -}}
	})
	if err != nil {
		return nil, err
	}
	{{template "after_hook" . }}

	return &res, nil
}
{{end}}`

const SpannerClientStreamingTemplate = `{{define "spanner_client_streaming_method"}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(stream {{.GetFilePackage}}{{.GetServiceName}}_{{.GetName}}Server) error {
	{{- $method := . -}}
	var err error
	_ = err
	feed, stop := s.PERSIST.{{.GetName}}(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, err.Error())
		}
		{{template "before_hook" .}}
		params := &persist_lib.{{template "persist_lib_input_name" $method}}{}
		{{range $key, $val := .GetTypeDescForQueryFields -}}
			{{if $val.IsMapped -}}
				if params.{{$val.Name}}, err = ({{$val.GoName}}{}).ToSpanner(req.{{.Name}}).SpannerValue(); err != nil {
					return gstatus.Errorf(codes.Unknown, "could not convert type: %v", err)
				}
			{{else -}}
				params.{{$val.Name}} =  req.{{$val.Name}}
			{{end -}}
		{{end}}

		feed(params)
	}

	row, err := stop()
	if err != nil {
		return err
	}
	res := {{.GetOutputType}}{}
	if row != nil {
		{{range $index, $field := .GetTypeDescArrayForStruct .GetOutputTypeStruct -}}
			{{if $field.IsMapped -}}
				var {{$field.Name}} *spanner.GenericColumnValue
				if err := row.ColumnByName("{{$field.ProtoName}}", {{$field.Name}}); err != nil {
					return gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{else}}
				var {{$field.Name}} {{$method.DefaultMapping $field.FieldDescriptor}}
				if err := row.ColumnByName("{{$field.ProtoName}}", &{{$field.Name}}); err != nil {
					return gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{end}}
			{{if $field.IsMapped -}}
				{
					local := &{{$field.GoName}}{}
					if err := local.SpannerScan({{$field.Name}}); err != nil {
						return gstatus.Errorf(codes.Unknown, "could not scan out custom type: %s", err)
					}
					res.{{$field.Name}} = local.ToProto()
				}
			{{else}}
				res.{{$field.Name}} = {{$field.Name}}
			{{end}}
		{{end}}
	}

	{{template "after_hook" .}}

	if err := stream.SendAndClose(&res); err != nil {
		return err
	}
	return nil
}
{{end}}
`
const SpannerServerStreamingTemplate = `{{define "spanner_server_streaming_method"}}// spanner server streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(req *{{.GetInputType}}, stream {{.GetFilePackage}}{{.GetServiceName}}_{{.GetName}}Server) error {
	{{- $method := . -}}
	var err error
	_ = err

	{{template "before_hook" .}}

	params := &persist_lib.{{template "persist_lib_input_name" $method}}{}
	{{range $key, $val := .GetTypeDescForQueryFields -}}
		{{if $val.IsMapped -}}
			if params.{{$val.Name}}, err = ({{$val.GoName}}{}).ToSpanner(req.{{.Name}}).SpannerValue(); err != nil {
				return gstatus.Errorf(codes.Unknown, "could not convert type: %v", err)
			}
		{{else -}}
			params.{{$val.Name}} =  req.{{$val.Name}}
		{{end -}}
	{{end}}

	var iterErr error
	err = s.PERSIST.{{.GetName}}(stream.Context(), params, func(row *spanner.Row) {
		if iterErr != nil || row == nil{
			return
		}
		var res {{.GetOutputType}}
		{{range $index, $field := .GetTypeDescArrayForStruct .GetOutputTypeStruct -}}
			{{if $field.IsMapped -}}
				var {{$field.Name}} *spanner.GenericColumnValue
				if err := row.ColumnByName("{{$field.ProtoName}}", {{$field.Name}}); err != nil {
					iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{else}}
				var {{$field.Name}} {{$method.DefaultMapping $field.FieldDescriptor}}
				if err := row.ColumnByName("{{$field.ProtoName}}", &{{$field.Name}}); err != nil {
					iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{end}}
			{{if $field.IsMapped -}}
				{
					local := &{{$field.GoName}}{}
					if err := local.SpannerScan({{$field.Name}}); err != nil {
						iterErr = gstatus.Errorf(codes.Unknown, "could not scan out custom type: %s", err)
						return
					}
					res.{{$field.Name}} = local.ToProto()
				}
			{{else}}
				res.{{$field.Name}} = {{$field.Name}}
			{{end -}}
		{{end -}}

		{{template "after_hook" .}}

		if err := stream.Send(&res); err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return err
	} else if iterErr != nil {
		return iterErr
	}

	return nil
}
{{end}}
`
