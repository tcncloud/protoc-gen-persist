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

package pkgimport_test

import (
	. "github.com/tcncloud/protoc-gen-persist/generator/pkgimport"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Import", func() {

})

var _ = Describe("Imports", func() {
	il := NewImports()
	Describe("function Exist()", func() {
		Describe("for fmt import path", func() {
			It("should return true", func() {
				Expect(il.Exist("fmt")).To(Equal(true))
			})
		})
		Describe("for github.com/tcncloud", func() {
			It("should return false", func() {
				Expect(il.Exist("github.com/tcncloud")).To(Equal(false))
			})
		})
	})
	Describe("function Append", func() {
		Describe("append github.com/namtzigla import", func() {
			It("should succeed", func() {
				il.Append(&Import{GoPackageName: "nam", GoImportPath: "github.com/nam"})
				Expect(il.Exist("github.com/nam")).To(Equal(true))
			})
		})

	})
	Describe("funcion Generate()", func() {
		Describe("for a new Imports structure", func() {
			It("should return a valid import statement", func() {
				i := NewImports()
				expectString := "import(\n \"fmt\" fmt \n \"sql\" database/sql \n \"driver\" database/sql/driver \n \"jsonpb\" github.com/golang/protobuf/jsonpb \n\n)"
				Expect(i.Generate()).To(Equal(expectString))
			})
		})
	})
})
