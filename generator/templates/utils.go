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

const ReturnConvertHelpers = `
	{{/* all three of these templates are used only in proto response construction */}}
	{{define "addr"}}{{if .IsMessage}}&{{end}}{{end}}
	{{define "base"}}{{if .IsEnum}}{{.EnumName}}({{.Name}}.ToProto()){{else}}{{.Name}}{{end}}{{end}}
	{{define "mapping"}}{{if and .IsMapped (not .IsEnum)}}.ToProto(){{end}}{{end}}
`

const BeforeHook = `
	{{define "before_hook"}}
	{{/* give it a Method as dot, assumes a "req" exists to give to the hook as parameter*/}}
	{{/* works for all method types, but assumes we are in a loop to continue on for */}}
	{{/* client streaming and bidi streaming methods */}}
	{{$before := .GetMethodOption.GetBefore -}}
		{{if $before -}}
		{{$pkg := .GetGoPackage $before.GetPackage -}}
		{{if eq $pkg "" -}}
			beforeRes, err := {{$before.GetName}}(req)
		{{else -}}
			beforeRes, err := {{$pkg}}.{{$before.GetName}}(req)
		{{end -}}
			if err != nil {
				{{if .IsUnary -}}
					return nil, grpc.Errorf(codes.Unknown, err.Error())
				{{else if (and (.IsClientStreaming) (not .IsSpanner))}}
					tx.Rollback()
					return grpc.Errorf(codes.Unknown, err.Error())
				{{else -}}
					return grpc.Errorf(codes.Unknown, err.Error())
				{{end -}}
			}
			if beforeRes != nil {
				{{if .IsClientStreaming -}}
					continue
				{{end}}
				{{if and .IsBidiStreaming (not .IsSpanner) -}}
					err = stream.Send(beforeRes)
					if err != nil {
						return grpc.Errorf(codes.Unknown, err.Error())
					}
					continue
				{{end -}}
				{{if or .IsUnary -}}
					return beforeRes, nil
				{{end -}}
				{{if .IsServerStreaming -}}
					for _, res := range beforeRes {
						err = stream.Send(res)
						if err != nil {
							return err
						}
					}
					return nil
				{{end -}}
			}
		{{end -}}
	{{end}}
`
const AfterHook = `
	{{define "after_hook"}}
	{{/* give it a Method as dot, assumes a "res" exists to give the hook as parameter */}}
	{{$after := .GetMethodOption.GetAfter -}}
		{{if $after -}}
			{{$pkg := .GetGoPackage $after.GetPackage -}}
			{{if eq $pkg "" -}}
				err = {{$after.GetName}}({{if (and (.IsClientStreaming) .IsSpanner)}}nil{{else}}req{{end}}, &res)
			{{else -}}
				err = {{.GetGoPackage $after.GetPackage}}.{{$after.GetName}}({{if (and (.IsClientStreaming) .IsSpanner)}}nil{{else}}req{{end}}, &res)
			{{end -}}
			if err != nil {
				{{if .IsUnary -}}
					return nil, grpc.Errorf(codes.Unknown, err.Error())
				{{else if (and (.IsClientStreaming) (not .IsSpanner))}}
					tx.Rollback()
					return grpc.Errorf(codes.Unknown, err.Error())
				{{else if (and (.IsServerStreaming) .IsSpanner) -}}
					iterErr = grpc.Errorf(codes.Unknown, err.Error())
				{{else -}}
					return grpc.Errorf(codes.Unknown, err.Error())
				{{end -}}
			}
		{{end -}}
	{{end}}
`
