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

const SpannerUnaryMethodTemplate = `{{define "spanner_unary_method"}}// spanner unary select {{.GetName}}
func (s* {{.GetServiceName}}Impl) {{.GetName}} (ctx context.Context, req *{{.GetInputType}}) (*{{.GetOutputType}}, error) {
{{if .Spanner.IsSelect}}{{template "spanner_unary_select" .}}{{end}}
{{if .Spanner.IsUpdate}}{{template "spanner_unary_update" .}}{{end}}
{{if .Spanner.IsInsert}}{{template "spanner_unary_insert" .}}{{end}}
{{if .Spanner.IsDelete}}{{template "spanner_unary_delete" .}}{{end}}
{{end}}`

const SpannerUnarySelectTemplate = `{{define "spanner_unary_select"}}
	var err error
	var (
{{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
		{{$field}} {{$type}}{{end}}
	)

	{{template "before_hook" .}}
	{{template "declare_spanner_arg_map" .}}

	//stmt := spanner.Statement{SQL: "{ {.Spanner.Query} }", Params: params}
	stmt := spanner.Statement{SQL: "{{.Spanner.Query}}", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err == iterator.Done {
		return nil, grpc.Errorf(codes.NotFound, "no rows found")
	} else if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	// scan our values out of the row
	{{range $index, $t := .GetTypeDescArrayForStruct .GetOutputTypeStruct}}
	{{if $t.IsMapped}}
	gcv{{$index}} := new(spanner.GenericColumnValue)
	err = row.ColumnByName("{{$t.ProtoName}}", gcv{{$index}})
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	err = {{$t.Name}}.SpannerScan(gcv{{$index}})
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	{{else}}
	err = row.ColumnByName("{{$t.ProtoName}}", &{{$t.Name}})
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	{{end}}{{end}}

	_, err = iter.Next()
	if err != iterator.Done {
		fmt.Println("Unary select that returns more than one row..")
	}
	res := {{.GetOutputType}}{
	{{range $field, $type := .GetTypeDescForFieldsInStruct .GetOutputTypeStruct}}
	{{$field}}: {{template "addr" $type}}{{template "base" $type}}{{template "mapping" $type}},{{end}}
	}

	{{template "after_hook" .}}

	return &res, nil
}
{{end}}`

const SpannerUnaryInsertTemplate = `{{define "spanner_unary_insert"}}
	var err error
	{{template "before_hook" .}}
	{{template "declare_spanner_arg_slice" .}}

	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Insert("{{.Spanner.TableName}}", {{.Spanner.InsertColsAsString}}, params)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := {{.GetOutputType}}{}

	{{template "after_hook" .}}

	return &res, nil
}
{{end}}`

const SpannerUnaryUpdateTemplate = `{{define "spanner_unary_update"}}
	var err error

	{{template "before_hook" .}}
	{{template "declare_spanner_arg_map" .}}

	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.UpdateMap("{{.Spanner.TableName}}", params)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil, grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := {{.GetOutputType}}{}

	{{template "after_hook" .}}

	return &res, nil
}
{{end}}`

const SpannerUnaryDeleteTemplate = `{{define "spanner_unary_delete"}}
	var err error

	{{template "before_hook" .}}
	{{template "declare_spanner_delete_key" .}}

	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.DeleteKeyRange({{.Spanner.TableName}}, key)
	_, err = s.SpannerDB.Apply(ctx, muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, grpc.Errorf(codes.NotFound, err.Error())
		}
	}
	res := {{.GetOutputType}}{}

	{{template "after_hook" .}}

	return &res, nil
}
{{end}}`

const SpannerClientStreamingMethodTemplate = `{{define "spanner_client_streaming_method"}}// spanner client streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(stream {{.GetServiceName}}_{{.GetName}}Server) error {
	var err error
	res := {{.GetOutputType}}{}
	{{$aft := .GetMethodOption.GetAfter}}
	{{if $aft}}
		reqs := make([]*{{.GetInputType}}, 0)
	{{end}}
	muts := make([]*spanner.Mutation, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err  != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		{{template "before_hook" .}}
		{{if $aft}}
			reqs = append(reqs, req)
		{{end}}

		{{if .Spanner.IsInsert}}{{template "spanner_client_streaming_insert" .}}{{end}}
		{{if .Spanner.IsUpdate}}{{template "spanner_client_streaming_update" .}}{{end}}
		{{if .Spanner.IsDelete}}{{template "spanner_client_streaming_delete" .}}{{end}}

		////////////////////////////// NOTE //////////////////////////////////////
		// In the future, we might do apply if muts gets really big,  but for now,
		// we only do one apply on the database with all the records stored in muts
		//////////////////////////////////////////////////////////////////////////
	}
	_, err = s.SpannerDB.Apply(context.Background(), muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	{{if $aft}}
		for _, req := range reqs {
			{{template "after_hook" .}}
		}
	{{end}}
	stream.SendAndClose(&res)
	return nil
}
{{end}}`

const SpannerClientStreamingUpdateTemplate = `{{define "spanner_client_streaming_update"}}//spanner client streaming update
{{template "declare_spanner_arg_map" .}}

muts = append(muts, spanner.UpdateMap("{{.Spanner.TableName}}", params))
{{end}}`

const SpannerClientStreamingInsertTemplate = `{{define "spanner_client_streaming_insert"}}//spanner client streaming insert
{{template "declare_spanner_arg_slice" .}}

	muts = append(muts, spanner.Insert("{{.Spanner.TableName}}", {{.Spanner.InsertColsAsString}}, params))
{{end}}`

const SpannerClientStreamingDeleteTemplate = `{{define "spanner_client_streaming_delete"}}//spanner client streaming delete
{{template "declare_spanner_delete_key" .}}

	muts = append(muts, spanner.DeleteKeyRange({{.Spanner.TableName}}, key))
{{end}}`

const SpannerServerStreamingMethodTemplate = `{{define "spanner_server_streaming_method"}}// spanner server streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(req *{{.GetInputType}}, stream {{.GetServiceName}}_{{.GetName}}Server) error {
	var (
	{{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
		{{$field}} {{$type}}{{end}}
	)

	{{template "before_hook" .}}

	{{if ne (len .Spanner.QueryArgs) 0}}
	var err error
	{{end}}

	{{template "declare_spanner_arg_map" .}}

	stmt := spanner.Statement{SQL: "{{.Spanner.Query}}", Params: params}
	tx := s.SpannerDB.Single()
	defer tx.Close()
	iter := tx.Query(context.Background(), stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}

		// scan our values out of the row
		{{range $index, $t := .GetTypeDescArrayForStruct .GetOutputTypeStruct}}
		{{if $t.IsMapped}}
		gcv{{$index}} := new(spanner.GenericColumnValue)
		err = row.ColumnByName("{{$t.ProtoName}}", gcv{{$index}})
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		err = {{$t.Name}}.SpannerScan(gcv{{$index}})
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		{{else}}
		err = row.ColumnByName("{{$t.ProtoName}}", &{{$t.Name}})
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		{{end}}{{end}}
		res := {{.GetOutputType}}{
		{{range $field, $type := .GetTypeDescForFieldsInStruct .GetOutputTypeStruct}}
		{{$field}}: {{template "addr" $type}}{{template "base" $type}}{{template "mapping" $type}},{{end}}
		}

		{{template "after_hook" .}}

		stream.Send(&res)
	}
	return  nil
}
{{end}}`

const SpannerBidiStreamingMethodTemplate = `{{define "spanner_bidi_streaming_method"}}// spanner bidi streaming {{.GetName}}
unimplemented
{{end}}`

const SpannerHelperTemplates = `
{{define "type_desc_to_def_map"}}
{{if .IsMapped}}
	conv, err = {{.GoName}}{}.ToSpanner(req.{{.Name}}).SpannerValue()
{{else}}
	conv = req.{{.Name}}
{{end}}{{end}}


{{define "type_desc_to_def_slice"}}
{{if .IsMapped}}
	conv, err = {{.GoName}}{}.ToSpanner(req.{{.Name}}).SpannerValue()
{{else}}
	conv = req.{{.Name}}
{{end}}{{end}}


{{define "return_err_on_method"}}
	if err != nil {
{{if .IsUnary}}
		return nil, grpc.Errorf(codes.Unknown, err.Error())
{{else}}
		return grpc.Errorf(codes.Unknown, err.Error())
{{end}}
	}
{{end}}


{{define "declare_spanner_arg_map"}}
{{$method := .}}
	params := make(map[string]interface{})
{{if gt (len .Spanner.OptionArguments) 0}}
	var conv interface{}
{{end}}
{{range $key, $val := .Spanner.QueryArgs}}
{{if $val.IsFieldValue}}
	{{template "type_desc_to_def_map" $val.Field}}
	{{template "return_err_on_method" $method}}
	params["{{$val.Name}}"] = conv
{{else}}
{{if $val.IsValue}}
	conv = {{$val.Value}}
	params["{{$val.Name}}"] = conv
{{end}}
{{end}}{{end}}
{{end}}


{{define "declare_spanner_arg_slice"}}
{{$method := .}}
	params := make([]interface{}, 0)
{{if gt (len .Spanner.OptionArguments) 0}}
	var conv interface{}
{{end}}

{{range $index, $val := .Spanner.QueryArgs}}
{{if $val.IsFieldValue}}
	{{template "type_desc_to_def_slice" $val.Field}}
	{{template "return_err_on_method" $method}}
	params = append(params, conv)
{{else}}
{{if $val.IsValue}}
	params = append(params, {{$val.Value}})
{{end}}
{{end}}{{end}}
{{end}}


{{define "declare_spanner_delete_key"}}
{{$method := .}}
	start := make([]interface{}, 0)
	end := make([]interface{}, 0)
{{if gt (len .Spanner.OptionArguments) 0}}
	var conv interface{}
{{end}}
{{range $index, $arg := .Spanner.KeyRangeDesc.Start}}
{{if $arg.IsFieldValue}}
{{template "type_desc_to_def_slice" $arg.Field}}
{{template "return_err_on_method" $method}}
	start = append(start, conv)
{{else}}
	start = append(start, {{$arg.Value}})
{{end}}{{end}}
{{range $index, $arg := .Spanner.KeyRangeDesc.End}}
{{if $arg.IsFieldValue}}
{{template "type_desc_to_def_slice" $arg.Field}}
{{template "return_err_on_method" $method}}
	end = append(end, conv)
{{else}}
{{if $arg.IsValue}}
	end = append(end, {{$arg.Value}})
{{end}}
{{end}}{{end}}
	key := spanner.KeyRange{
		Start: start,
		End: end,
		Kind: spanner.{{.Spanner.KeyRangeDesc.Kind}},
	}
{{end}}

`
