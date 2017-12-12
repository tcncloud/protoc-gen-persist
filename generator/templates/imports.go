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

const ImportTemplate = `{{define "import_template"}}
{{$file := .}}

import({{range $import := .ImportList}}
{{if $file.NotSameAsMyPackage .GoImportPath}}{{$import.GoPackageName}} "{{$import.GoImportPath}}"{{end}}{{end}}
)
{{/*
var _ = "{{$file.GetPackageName}}" // file.GetPackageName()
var _ = "{{$file.GetOrigName}}" // file.GetOrigName()
var _ = "{{$file.GetImplFileName}}" // file.GetImplFileName()
var _ = "{{$file.GetImplDir}}" // file.GetImplDir()
var _ = "{{$file.GetImplPackage}}" // file.GetImplPackage()
var _ = "{{$file.GetFileName}}" //file.GetFileName()
var _ = "{{$file.GetGoPackage}}" // file.GetGoPackage()
var _ = "{{$file.GetFullGoPackage}}" // file.GetFullGoPackage()
var _ = "{{$file.GetGoPath}}" // file.GetGoPath()
var _ = "{{$file.NeedsPersistLibDir}}" // file.NeedsPersistLibDir()
var _ = "{{$file.DifferentImpl}}" // file.DifferentImpl()
*/}}

{{end}}`
