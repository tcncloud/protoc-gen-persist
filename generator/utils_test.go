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

package generator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tcncloud/protoc-gen-persist/generator"
)

var _ = Describe("IsServicePersistEnabled", func() {
	Describe("for a service that implement custom extension persist.ql", func() {
		It("should return true", func() {
			Expect(generator.IsServicePersistEnabled(descr.Service[0])).To(Equal(true))
		})
	})
})

var _ = Describe("IsMethodEnabled", func() {
	Describe("For a method that implement persist.ql", func() {
		It("should return true", func() {
			Expect(generator.IsMethodEnabled(descr.Service[0].Method[0])).To(Equal(true))
		})
	})
	Describe("For a method that does not implement persist.ql", func() {
		It("should return false", func() {
			Expect(generator.IsMethodEnabled(descr.Service[0].Method[1])).To(Equal(false))
		})
	})
	Describe("For a nil parameter", func() {
		It("should return false", func() {
			Expect(generator.IsMethodEnabled(descr.Service[0].Method[1])).To(Equal(false))
		})
	})
})

var _ = Describe("GetMethodExtensionData", func() {
	Describe("For UnaryExample1 method", func() {
		It("should return a structure", func() {
			Expect(generator.GetMethodExtensionData(descr.Service[0].Method[0])).ToNot(BeNil())
		})
	})
	Describe("For RandomMethod method", func() {
		It("should return nil", func() {
			Expect(generator.GetMethodExtensionData(descr.Service[0].Method[1])).To(BeNil())
		})
	})
})
