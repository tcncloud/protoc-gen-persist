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

const SqlUnaryMethodTemplate = `{{define "sql_unary_method"}}// sql unary {{.GetName}} unimplemented
func (s* {{.GetServiceName}}Impl) {{.GetName}} (ctx context.Context, req *{{.GetInputType}}) (*{{.GetOutputType}}, error) {
	var (
{{range $field, $type := .GetFieldsWithLocalTypesFor .GetOutputTypeStruct}}
{{$field}} {{$type}}
{{end}}
	)
	err := s.SqlDB.QueryRow({{.GetQuery}} {{.GetQueryParamString}}).
		Scan(&id, &start_time, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, grpc.Errorf(codes.NotFound, "%+v doesn't exist", req)
		} else if strings.Contains(err.Error(), "duplicate key") {
			return nil, grpc.Errorf(codes.AlreadyExists, "%+v already exists")
		}
		return nil, grpc.Errorf(codes.Unknown, err.Error())
	}
	res := &pb.Table2{
		Id: id,
		StartTime: start_time.UnpackSql(),
		Name: name,
	}
	return res, nil
}

{{end}}`
const SqlClientStreamingMethodTemplate = `{{define "sql_client_streaming_method"}}// sql client streaming {{.GetName}} unimplemented{{end}}`
const SqlServerStreamingMethodTemplate = `{{define "sql_server_streaming_method"}}// sql server streaming {{.GetName}} unimplemented{{end}}`
const SqlBidiStreamingMethodTemplate = `{{define "sql_bidi_streaming_method"}}// sql bidi streaming {{.GetName}} unimplemented{{end}}`
