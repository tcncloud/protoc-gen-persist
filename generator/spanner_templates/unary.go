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

package spanner_templates

const SpannerUnaryTemplate = `{{define "spanner_unary_method"}}
func (s* {{.GetServiceName}}Impl) {{.GetName}} (ctx context.Context, req *{{.GetInputType}}) (*{{.GetOutputType}}, error) {
	{{$method := .}}
	{{template "before_hook" .}}

	params := map[string]interface{}{}
	{{range $key, $val := .GetTypeDescForFieldsInStructSnakeCase .GetInputTypeStruct}}
		{{if $val.IsMapped}}
			if params["@{{$key}}"], err = ({{$val.GoName}}{}).ToSpanner(req.{{.Name}}).SpannerValue(); err != nil {
				return gstatus.Errorf(codes.Unknown, "could not convert type: %v", err)
			}
		{{else}}
			params["@{{$key}}"] =  req.{{$val.Name}}
		{{end}}
	{{end}}

	var res *{{.GetOutputType}}
	var iterErr error
	s.PERSIST.{{.GetName}}(ctx, params, func(row map[string]interface{}) {
		var ok bool
		{{range $index, $field := .GetTypeDescArrayForStruct .GetOutputTypeStruct}}
			{{if $field.IsMapped}}
				var {{$field.Name}} *spanner.GenericColumnValue
				{{$field.Name}}, ok = row["$field.ProtoName"].(*spanner.GenericColumValue)
				if !ok {
					iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{else}}
				var {{$field.Name}} {{$method.DefaultMapping $field.FieldDescriptor}}
				{{$field.Name}}, ok = row["$field.ProtoName"].({{$method.DefaultMapping $field.FieldDescriptor}})
				if !ok {
					iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				}
			{{end}}
			{{if $field.IsMapped}}
				res.{{$field.Name}}, err = ({{$field.GoName}}{}).SpannerScan({{$field.Name}}).ToProto()
				if err != nil {
					iterErr = gstatus.Errorf(codes.Unknown, "could not scan out custom type: %s", err)
				}
			{{else}}
				res.{{$field.Name}} = {{$field.Name}}
			{{end}}
		{{end}}
	})
	{{template "after_hook" . }}

	return res, nil
}
{{end}}`

const SpannerClientStreamingTemplate = `{{define "spanner_client_streaming_method"}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(stream {{.GetFilePackage}}{{.GetServiceName}}_{{.GetName}}Server) error {
	{{$method := .}}
	var err error
	feed, stop := s.PERSIST.{{.GetName}}(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, err.Error())
		}
		{{template "before_hook" .}}

		params := map[string]interface{}{}

		{{range $key, $val := .GetTypeDescForFieldsInStructSnakeCase .GetInputTypeStruct}}
			{{if $val.IsMapped}}
				if params["@{{$key}}"], err = ({{$val.GoName}}{}).ToSpanner(req.{{.Name}}).SpannerValue(); err != nil {
					return gstatus.Errorf(codes.Unknown, "could not convert type: %v", err)
				}
			{{else}}
				params["@{{$key}}"] =  req.{{$val.Name}}
			{{end}}
		{{end}}
		if err := feed(params); err != nil {
			return err
		}
	}

	row := stop()
	res := {{.GetOutputType}}{}

	var ok bool
	{{range $index, $field := .GetTypeDescArrayForStruct .GetOutputTypeStruct}}
		{{if $field.IsMapped}}
			var {{$field.Name}} *spanner.GenericColumnValue
			{{$field.Name}}, ok = row["$field.ProtoName"].(*spanner.GenericColumValue)
			if !ok {
				iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
			}
		{{else}}
			var {{$field.Name}} {{$method.DefaultMapping $field.FieldDescriptor}}
			{{$field.Name}}, ok = row["$field.ProtoName"].({{$method.DefaultMapping $field.FieldDescriptor}})
			if !ok {
				iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
			}
		{{end}}
		{{if $field.IsMapped}}
			res.{{$field.Name}}, err = ({{$field.GoName}}{}).SpannerScan({{$field.Name}}).ToProto()
			if err != nil {
				iterErr = gstatus.Errorf(codes.Unknown, "could not scan out custom type: %s", err)
			}
		{{else}}
			res.{{$field.Name}} = {{$field.Name}}
		{{end}}
	{{end}}

	{{template "after_hook" .}}

	if err := stream.SendAndClose(res); err != nil {
		return err
	}
	return nil
}
{{end}}
`
const SpannerServerStreamingTemplate = `{{define "spanner_server_streaming_method"}}// spanner server streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(req *{{.GetInputType}}, stream {{.GetFilePackage}}{{.GetServiceName}}_{{.GetName}}Server) error {
	{{$method := .}}
	var err error
	_ = err

	{{template "before_hook" .}}

	params := map[string]interface{}{}

	{{range $key, $val := .GetTypeDescForFieldsInStructSnakeCase .GetInputTypeStruct}}
		{{if $val.IsMapped}}
			if params["@{{$key}}"], err = ({{$val.GoName}}{}).ToSpanner(req.{{.Name}}).SpannerValue(); err != nil {
				return gstatus.Errorf(codes.Unknown, "could not convert type: %v", err)
			}
		{{else}}
			params["@{{$key}}"] =  req.{{$val.Name}}
		{{end}}
	{{end}}

	var iterErr error
	persist_lib.{{.GetName}}(ctx, params, func(row map[string]interface{}) {
		if iterErr != nil {
			return
		}
		var res *{{.GetOutputType}}
		var ok bool
		{{range $index, $field := .GetTypeDescArrayForStruct .GetOutputTypeStruct}}
			{{if $field.IsMapped}}
				var {{$field.Name}} *spanner.GenericColumnValue
			{{else}}
				var {{$field.Name}} {{.DefaultMapping $field.FieldDescriptor}}
			{{end}}
			{{$field.Name}}, ok = row["$field.ProtoName"].({{.DefaultMapping $field.FieldDescriptor}})
			if !ok {
				iterErr = gstatus.Errorf(codes.Unknown, "could not convert type %v", err)
				return
			}
			{{if $field.IsMapped}}
				res.{{$field.Name}} = "something.."
			{{else}}
				res.{{$field.Name}} = {{$field.Name}}
			{{end}}
		{{end}}

		{{template "after_hook" .}}

		if err := stream.Send(res); err != nil {
			iterError = err
			return
		}
	})
	return nil
}
{{end}}
`

// an example of a unary template
// import (
// 	mytime "github.com/tcncloud/protoc-gen-persist/examples/mytime"
// 	pb "github.com/tcncloud/protoc-gen-persist/examples/spanner/basic"
// 	hooks "github.com/tcncloud/protoc-gen-persist/examples/spanner/hooks"
// 	test "github.com/tcncloud/protoc-gen-persist/examples/test"
// 	context "golang.org/x/net/context"
// 	iterator "google.golang.org/api/iterator"
// 	grpc "google.golang.org/grpc"
// 	codes "google.golang.org/grpc/codes"
// 	gstatus "google.golang.org/grpc/gstatus"
// )

// func (s *MySpannerImpl) ClientStream(stream pb.MySpanner_ClientStreamInsertServer) error {
// 	var err error

// 	feed, stop := mtime.ClientStream(stream.context())

// 	for {
// 		req, err := stream.Recv()
// 		if err == io.EOF {
// 			break
// 		} else if err != nil {
// 			return grpc.Errorf(codes.Unknown, err.Error())
// 		}
// 		res, err := mytime.BeforeHook(req)
// 		if err != nil {
// 			return err
// 		} else if res != nil {
// 			// skip
// 			continue
// 		}
// 		params := map[string]interface{}{
// 			"@col1": req.Field1,
// 			"@col2": mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue(),
// 		}
// 		if err := feed(params); err != nil {
// 			return err
// 		}
// 	}
// 	row := stop()
// 	var field1 string
// 	var start_time *spanner.GenericColumnValue
// 	var ok bool
// 	field1, ok = row["field1"].(string)
// 	if !ok {
// 		return nil, gstatus.Errorf("could not convert %+v, to type %T", row["field1"], field1)
// 	}
// 	start_time, ok = row["start_time"].(*spanner.GenericColumnValue)
// 	if !ok {
// 		return nil, gstatus.Errorf("could not convert %+v, to type %T needed for response field: %s",
// 			row["start_time"], start_time, "StartTime")
// 	}

// 	StartTime = (mytime.MyTime{}).SpannerScan(start_time)

// 	if err := mytime.AfterHook(req, res); err != nil {
// 		return gstatus.Errorf(codes.Unknown, err.Error())
// 	}

// 	return &test.ExampleTable{
// 		Field1:    field1,
// 		StartTime: StartTime.ToProto(),
// 	}, nil
// }
// func (s *MySpannerImpl) UniaryInsert(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
// 	params := map[string]interface{}{
// 		"@col1": req.Field1,
// 		"@col2": mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue(),
// 	}
// 	var res *test.ExampleTable

// 	row, err := mytime.UniaryInsert(ctx, params, func(r *spanner.Row) {
// 		var field1 string
// 		var start_time *spanner.GenericColumnValue
// 		var ok bool
// 		field1, ok = row["field1"].(string)
// 		if !ok {
// 			return nil, gstatus.Errorf("could not convert %+v, to type %T", row["field1"], field1)
// 		}
// 		start_time, ok = row["start_time"].(*spanner.GenericColumnValue)
// 		if !ok {
// 			return nil, gstatus.Errorf("could not convert %+v, to type %T needed for response field: %s",
// 				row["start_time"], start_time, "StartTime")
// 		}

// 		StartTime = (mytime.MyTime{}).SpannerScan(start_time)

// 	})
// 	if err != nil {
// 		return gstatus.Errorf(codes.Unknown, err.Error())
// 	}

// 	return &test.ExampleTable{
// 		Field1:    field1,
// 		StartTime: StartTime.ToProto(),
// 	}, nil
// }

// // package pb/persist_lib
// func UniaryInsertMut(reqs map[string]interface{}) *spanner.Mutation {
// 	// from the query
// 	return spanner.InsertMap("table", map[string]interface{}{
// 		"field1":     req["@col1"],
// 		"start_time": req["@col2"],
// 		"field3":     3.3,
// 	})
// }

// // func(s *{{.ServiceName}}Impl) {{.MethodName}}(ctx context.Context, req *{{.RequestPackage}}.{{.RequestName}})
// // (*{{.ResponsePackage}}.{{ResponseName}}, error) {
// func (s *MySpannerImpl) UniaryInsert(ctx context.Context, req *test.ExampleTable) (*test.ExampleTable, error) {
// 	//params := map[string]interface{}{
// 	params := map[string]interface{}{
// 		//"{{.ArgKey}}": req.{{.FieldName}}
// 		"@col1": req.Field1,
// 		//"{{.ArgKey}}
// 		"@col2": mytime.MyTime{}.ToSpanner(req.StartTime).SpannerValue(),
// 	}
// 	// var res *{{.ResponsePackage}}.{{.ResponseName}}
// 	var res *test.ExampleTable

// 	//row, err := persist_lib.{{MethodName}}(ctx, parms, func(r *spanner.Row) {
// 	row, err := mytime.UniaryInsert(ctx, params, func(r *spanner.Row) {
// 		var field1 string
// 		var start_time *spanner.GenericColumnValue
// 		var ok bool
// 		field1, ok = row["field1"].(string)
// 		if !ok {
// 			return nil, gstatus.Errorf("could not convert %+v, to type %T", row["field1"], field1)
// 		}
// 		start_time, ok = row["start_time"].(*spanner.GenericColumnValue)
// 		if !ok {
// 			return nil, gstatus.Errorf("could not convert %+v, to type %T needed for response field: %s",
// 				row["start_time"], start_time, "StartTime")
// 		}

// 		StartTime = (mytime.MyTime{}).SpannerScan(start_time)

// 	})
// 	if err != nil {
// 		return gstatus.Errorf(codes.Unknown, err.Error())
// 	}

// 	return &test.ExampleTable{
// 		Field1:    field1,
// 		StartTime: StartTime.ToProto(),
// 	}, nil
// }
