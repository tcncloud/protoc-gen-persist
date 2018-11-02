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
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tcncloud/protoc-gen-persist/persist"

	"golang.org/x/tools/imports"
)

var reduceEmptyLines = regexp.MustCompile("(\n)+")

// GetGoPath get a go import url under the following formats
// github.com/path/project/dir;package
// github.com/path/project/dir
// project/dir;package
// project/dir
// and will return the path portion from url:
// github.com/path/project/dir
// project/dir
func GetGoPath(url string) string {
	idx := strings.LastIndex(url, ";")
	switch {
	case idx >= 0:
		return url[0:idx]
	default:
		return url
	}
}

// GetGoPackage get a go import url under the following formats
// github.com/path/project/dir;package
// github.com/path/project/dir
// project/dir;package
// project/dir
// and will return the package name from url
// package
// dir
// package
// dir
func GetGoPackage(url string) string {
	switch {
	case strings.Contains(url, ";"):
		idx := strings.LastIndex(url, ";")
		return url[idx+1:]
	case strings.Contains(url, "/"):
		idx := strings.LastIndex(url, "/")
		return url[idx+1:]
	default:
		return url
	}
}

func FormatCode(filename string, buffer []byte) []byte {
	// reduce the empty lines
	tmp := reduceEmptyLines.ReplaceAll(buffer, []byte{'\n'})
	buf, err := imports.Process(filename, tmp, nil)
	if err != nil {
		logrus.WithError(err).Errorf("Error processing file %s", filename)
		return tmp
	}
	return buf
}
func getGoNamesForTypeMapping(tm *persist.TypeMapping_TypeDescriptor, file *FileStruct) (string, string) {
	name := file.GetGoTypeName(tm.GetProtoTypeName())
	nameParts := strings.Split(name, ".")
	for i, v := range nameParts {
		nameParts[i] = strings.Title(v)
	}
	titled := strings.Join(nameParts, "")
	return name, titled
}
