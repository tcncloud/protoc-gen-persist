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

const MethodTemplate = `{{define "implement_method"}}
{{if .IsUnary}} {{template "unary_method" .}} {{end}}
{{if .IsClientStreaming}} {{template "client_streaming_method" .}} {{end}}
{{if .IsServerStreaming}} {{template "server_streaming_method" .}} {{end}}
{{if .IsBidiStreaming}} {{template "bidi_method" .}} {{end}}
{{end}}`

const UnaryMethodTemplate = `{{define "unary_method"}}
{{if .IsSQL}}{{template "sql_unary_method" .}}{{end}}
{{if .IsSpanner}}{{template "spanner_unary_method" .}}{{end}}
{{end}}`

const ClientStreamingMethodTemplate = `{{define "client_streaming_method"}}
{{if .IsSQL}}{{template "sql_client_streaming_method" .}}{{end}}
{{if .IsSpanner}}{{template "spanner_client_streaming_method" .}}{{end}}
{{end}}`

const ServerStreamingMethodTemplate = `{{define "server_streaming_method"}}
{{if .IsSQL}}{{template "sql_server_streaming_method" .}}{{end}}
{{if .IsSpanner}}{{template "spanner_server_streaming_method" .}}{{end}}
{{end}}`

const BidiStreamingMethodTemplate = `{{define "bidi_method"}}
{{if .IsSQL}}{{template "sql_bidi_streaming_method" .}}{{end}}
{{end}}`
