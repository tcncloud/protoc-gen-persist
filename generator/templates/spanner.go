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


const SpannerHelperTemplates = `
{{define "type_desc_to_def"}}
{{if .IsMapped}}
	//is mapped
	conv, err = GoName{}.ToSpanner(req.{{.Name}}).Value()
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
{{else}}
	// is not mapped
	conv = req.{{.Name}}
{{end}}{{end}}
`
const SpannerUnarySelectTemplate = `{{define "spanner_unary_select"}}
	var (
{{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
		{{$field}} {{$type}}{{end}}
	)
	params := make(map[string]interface{})

	var conv string
	var err error
	//.GetSpannerSelectArgs
{{range $key, $val := .GetSpannerSelectArgs}}
{{if $val.IsFieldValue}}
	//if is.IsFieldValue
	{{template "type_desc_to_def" $val.Field}}
	params["{{$val.Name}}"] = conv
{{else}}
	//else
	//conv = { {$val.Value} }
	conv = {{$val.Value}}
	//params[{ {$val.Name} }] = conv
	params["{{$val.Name}}"] = conv
	{{end}}{{end}}

	//stmt := spanner.Statement{SQL: "{ {.Spanner.Query} }", Params: params}
	stmt := spanner.Statement{SQL: "{{.Spanner.Query}}", Params: params}
	tx := s.Client.Single()
	defer tx.Close()
	iter := tx.Query(ctx, stmt)
	rows := s.SRH.NewRowsFromIter(iter)
	rows.Next()
	if err = rows.Err(); err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	//err = rows.Scan({ {range $index, $t := .GetTypeDescArrayForStruct .GetOutputTypeStruct} } &{ {$t.Name} },{ {end} })
	err = rows.Scan({{range $index, $t := .GetTypeDescArrayForStruct .GetOutputTypeStruct}} &{{$t.Name}},{{end}})
	if err == sql.ErrNoRows {
		return nil, grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
	} else if err != nil {
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}

	//res := &{ {.GetOutputType} }{
	res := &{{.GetOutputType}}{
	{{range $field, $type := .GetTypeDescForFieldsInStruct .GetOutputTypeStruct}}
	{{$field}}: {{template "addr" $type}}{{template "base" $type}}{{template "mapping" $type}},{{end}}
	}
	return res, nil
}
{{end}}`

const SpannerUnaryInsertTemplate = `{{define "spanner_unary_insert"}}
	params := []interface{}{
		{{range $index, $val := .GetSpannerInsertArgs "req"}}
		{{$val}},{{end}}
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.Insert("{{.Spanner.TableName}}", {{.Spanner.InsertCols}}, params)
	_, err := s.SpannerDB.Apply(muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			returnn nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := &{{.GetOutputType}}{}

	return res, nil
}
{{end}}`

const SpannerUnaryUpdateTemplate = `{{define "spanner_unary_update"}}
	params := map[string]interface{}{
		{{range $key, $val := .GetSpannerUpdateArgs "req"}}
		{{$key}}: {{$val}},\n{{end}}
	}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.UpdateMap("{{.Spanner.TableName}}", params)
	_, err := s.SpannerDB.Apply(muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	res := &{{.GetOutputType}}{}

	return res, nil
}
{{end}}`

const SpannerUnaryDeleteTemplate = `{{define "spanner_unary_delete"}}
	key := {{.GetDeleteKeyRange}}
	muts := make([]*spanner.Mutation, 1)
	muts[0] = spanner.DeleteKeyRange("{{.Spanner.TableName}}", key)
	_, err := s.SpannerDB.Apply(muts)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return grpc.Errorf(codes.NotFound, err.Error())
		}
	}
{{end}}`

const SpannerClientStreamingMethodTemplate = `{{define "spanner_client_streaming_method"}}// spanner client streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(stream {{.GetServiceName}}_{{.GetName}}Server) error {
	var totalAffected int64
	muts := make([]*spanner.Mutation, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		totalAffected += 1

		{{if .Spanner.IsInsert}}{{template "spanner_client_streaming_insert" .}}{{end}}
		{{if .Spanner.IsUpdate}}{{template "spanner_client_streaming_update" .}}{{end}}
		{{if .Spanner.IsDelete}}{{template "spanner_client_streaming_delete" .}}{{end}}
		//In the future, we might do apply if muts gets really big,  but for now,
		// we only do one apply on the database with all the records stored in muts
	}
	_, err := s.SpannerDB.Apply(muts)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return grpc.Errorf(codes.AlreadyExists, err.Error())
		} else {
			return nil, grpc.Errorf(codes.Unknown, err.Error())
		}
	}
	stream.SendAndClose(&{{.GetOutputType}}{Count: totalAffeted})

{{end}}`

const SpannerClientStreamingUpdateTemplate = `{{define "spanner_client_streaming_update"}}//spanner client streaming update
params := map[string]interface{}{
{{range $key, $val := .GetSpannerUpdateArgs}}
	{{$key}}: {{$val}},\n{{end}}
}
muts = append(muts, spanner.UpdateMap("{{.Spanner.TableName}}", params))
{{end}}`

const SpannerClientStreamingInsertTemplate = `{{define "spanner_client_streaming_insert"}}//spanner client streaming update
params := []interface{}{
{{range $index, $val := .GetSpannerInsertArgs}}
	{{$val}},\n{{end}}
}
muts = append(muts, spanner.Insert("{{.Spanner.TableName}}", {{.Spanner.InsertColsAsString}}, params))
{{end}}`

const SpannerClientStreamingDeleteTemplate = `{{define "spanner_client_streaming_delete"}}//spanner client streaming update
key := {{.GetDeleteKeyRange "req"}}
muts = append(muts, spanner.DeleteKeyRange("{{.Spanner.TableName}}", key)
{{end}}`

const SpannerServerStreamingMethodTemplate = `{{define "spanner_server_streaming_method"}}// spanner server streaming {{.GetName}}
func (s *{{.GetServiceName}}Impl) {{.GetName}}(req *{{.GetInputType}}, stream {{.GetServiceName}}_{{.GetName}}Server) error {
	var (
	{{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
		{{$field}} {{$type}}{{end}}
	)
	params := make(map[string]interface{})

	var conv string
	var err error
	//.GetSpannerSelectArgs
{{range $key, $val := .GetSpannerSelectArgs}}
{{if $val.IsFieldValue}}
	//if is.IsFieldValue
	{{template "type_desc_to_def" $val.Field}}
	params["{{$val.Name}}"] = conv
{{else}}
	//else
	//conv = { {$val.Value} }
	conv = {{$val.Value}}
	//params["{ {$val.Name} }"] = conv
	params["{{$val.Name}}"] = conv
{{end}}{{end}}

	stmt := spanner.Statement{SQL: "{{.Spanner.Query}}", Params: params}
	tx := s.Client.Single()
	defer tx.Close()
	iter := tx.Query(ctx, stmt)
	rows := s.SRH.NewRowsFromIter(iter)
	for rows.Next() {
		if err := rows.Err(); err != nil {
			if err == sql.ErrNowRows {
				return grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
			}
			return grpc.Errorf(codes.nknown, err.Error())
		}
		err := rows.Scan({{range $index, $t := .GetTypeDescArrayForStruct .GetOutputTypeStruct}} &{{$t.Name}},{{end}})
		if err != nil {
			return grpc.Errorf(codes.Unknown, err.Error())
		}
		res := &{{.GetOutputType}}{
		{{range $field, $type := .GetTypeDescForFieldsInStruct .GetOutputTypeStruct}}
		{{$field}}: {{template "addr" $type}}{{template "base" $type}}{{template "mapping" $type}},{{end}}
		}
		stream.Send(res)
	}
	return  nil
}
{{end}}`

const SpannerBidiStreamingMethodTemplate = `{{define "spanner_bidi_streaming_method"}}// spanner bidi streaming {{.GetName}} unimplemented{{end}}`
