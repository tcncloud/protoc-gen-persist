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

package utils

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

// GetGoPackageAndPathFromURL returns go package name & path from a go path url
// url formats supported
//  - github.com/package;package_go -> package_go, github.com/package
//  - github.com/package -> package, github.com/package
//  - pacakge/test -> test, package/test
//  - package -> package, package
func GetGoPackageAndPathFromURL(goURL string) (string, string) {
	if idx := strings.LastIndex(goURL, ";"); idx >= 0 {
		pkg := goURL[idx+1:]
		path := goURL[0 : idx-1]
		return pkg, path
	} else if idx := strings.LastIndex(goURL, "/"); idx >= 0 {
		pkg := goURL[idx+1:]
		path := goURL
		return pkg, path
	} else {
		return goURL, goURL
	}
}

// check if a service has at least one method that has the persist.ql extension defined
func IsServicePersistEnabled(service *descriptor.ServiceDescriptorProto) bool {
	if service.Method != nil {
		for _, method := range service.Method {
			if IsMethodEnabled(method) {
				return true
			}
		}
	}
	return false
}

func IsMethodEnabled(method *descriptor.MethodDescriptorProto) bool {
	if method != nil && method.GetOptions() != nil && proto.HasExtension(method.Options, persist.E_Ql) {
		return true
	}
	return false
}

func GetMethodOption(method *descriptor.MethodDescriptorProto) *persist.QLImpl {
	if IsMethodEnabled(method) {
		if ret, err := proto.GetExtension(method.Options, persist.E_Ql); err == nil {
			return ret.(*persist.QLImpl)
		}
	}
	return nil
}
