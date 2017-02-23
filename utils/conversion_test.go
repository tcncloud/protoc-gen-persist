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

package utils_test

import (
	"testing"
	"time"

	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	utils "github.com/tcncloud/protoc-gen-persist/utils"
)

func TestConversion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Conversion methods Suite")
}

var _ = Describe("protobuf message value conversion functions", func() {
	var _ = Describe("ToSafeType", func() {
		It("can convert []string => *pq.StringArray", func() {
			in := []string{"hello", "world"}
			out := utils.ToSafeType(in)

			Expect(out).To(BeAssignableToTypeOf(&pq.StringArray{}))

			back, ok := out.(*pq.StringArray)

			Expect(ok).To(Equal(true))

			Expect([]string(*back)).To(BeEquivalentTo([]string{"hello", "world"}))
		})

		It("can convert []int64 => *pq.Int64Array", func() {
			in := []int64{1, 2}
			out := utils.ToSafeType(in)

			Expect(out).To(BeAssignableToTypeOf(&pq.Int64Array{}))

			back, ok := out.(*pq.Int64Array)

			Expect(ok).To(Equal(true))
			Expect([]int64(*back)).To(BeEquivalentTo([]int64{1, 2}))
		})

		It("can convert []float64 => *pq.Float64Array", func() {
			in := []float64{1.123, 2.234}
			out := utils.ToSafeType(in)

			Expect(out).To(BeAssignableToTypeOf(&pq.Float64Array{}))

			back, ok := out.(*pq.Float64Array)

			Expect(ok).To(Equal(true))
			Expect([]float64(*back)).To(BeEquivalentTo([]float64{1.123, 2.234}))
		})

		It("can convert *google protobuf Timestamp => *time.Time", func() {
			in := &google_protobuf.Timestamp{
				Seconds: int64(1234567),
				Nanos:   12345,
			}
			out := utils.ToSafeType(in)

			con, ok := out.(*time.Time)

			Expect(ok).To(Equal(true))

			back := utils.ToProtobufTime(con)

			Expect(back).To(BeEquivalentTo(&google_protobuf.Timestamp{
				Seconds: int64(1234567),
				Nanos:   12345,
			}))
		})

		It("can convert default Valuer types", func() {
			if _, ok := utils.ToSafeType("hello").(*string); !ok {
				Fail("failed on string")
			}
			if _, ok := utils.ToSafeType(int64(1234)).(*int64); !ok {
				Fail("failed on int64")
			}
			if _, ok := utils.ToSafeType(int32(1234)).(*int32); !ok {
				Fail("failed on int32")
			}
			if _, ok := utils.ToSafeType(float64(1234.1234)).(*float64); !ok {
				Fail("failed on float64")
			}
		})
	})

	var _ = Describe("AssignTo", func() {
		It("can assign a pq.StringArray value to a *[]string key", func() {
			key := new([]string)
			val := pq.StringArray{"first", "second"}

			utils.AssignTo(key, val)
			Expect(*key).To(Equal([]string{"first", "second"}))
		})

		It(" can assign a pq.Int64Array value to a *[]int64 key", func() {
			key := new([]int64)
			val := pq.Int64Array{int64(1), int64(2)}

			utils.AssignTo(key, val)
			Expect(*key).To(Equal([]int64{int64(1), int64(2)}))
		})

		It("can assign a pq.Float64Array value to a *[]float64 key", func() {
			key := new([]float64)
			val := pq.Float64Array{1.1234, 2.2345}

			utils.AssignTo(key, val)
			Expect(*key).To(Equal([]float64{1.1234, 2.2345}))
		})

		It("can assign a int64 value to a *int64 key", func() {
			key := new(int64)
			val := int64(1234)

			utils.AssignTo(key, val)
			Expect(*key).To(Equal(int64(1234)))
		})

		It("can assign a int32 value to a *int32 key", func() {
			key := new(int32)
			val := int32(1234)

			utils.AssignTo(key, val)
			Expect(*key).To(Equal(int32(1234)))
		})

		It("can assign a float64 value to a *float64 key", func() {
			key := new(float64)
			val := 100.1234

			utils.AssignTo(key, val)
			Expect(*key).To(Equal(100.1234))
		})

		It("can assign a string value to a *string key", func() {
			key := new(string)
			val := "hello world"

			utils.AssignTo(key, val)
			Expect(*key).To(Equal("hello world"))
		})

		It("can assign a *time.Time value to a **google Timestamp key", func() {
			key := new(*google_protobuf.Timestamp)
			val := new(time.Time)
			*val = time.Now().UTC().Truncate(time.Microsecond)

			utils.AssignTo(key, val)

			convertedBack := ((*utils.ToTime(*key)).Truncate(time.Microsecond))
			Expect(val.String()).To(Equal(convertedBack.String()))
		})

		It("returns true if value is assigned to key", func() {
			key := new(string)
			val := ""

			res := utils.AssignTo(key, val)
			Expect(res).To(Equal(true))
		})

		var _ = Context("when no supported value is placed as val", func() {
			It("returns false", func() {
				key := new([]string)
				val := []string{"not", "supported"}

				res := utils.AssignTo(key, val)

				Expect(res).To(Equal(false))
			})

			It("does not assign to key", func() {
				key := new([]string)
				val := []string{"not", "supported"}

				utils.AssignTo(key, val)
				k := new([]string)
				Expect(key).To(BeEquivalentTo(k))
			})
		})

		var _ = Context("when key cannot convert to correct type", func() {
			It("returns false", func() {
				key := new(int64)
				val := []string{"not", "integers"}

				res := utils.AssignTo(key, val)
				Expect(res).To(Equal(false))
			})

			It("does not assign to key", func() {
				key := new(int64)
				val := []string{"not", "integers"}

				utils.AssignTo(key, val)
				Expect(key).To(BeEquivalentTo(new(int64)))
			})
		})

		var _ = Context("when key is nil", func() {
			It("returns false", func() {
				var key *int64
				val := int64(12345)

				res := utils.AssignTo(key, val)
				Expect(res).To(Equal(false))
			})

			It("does not assign to key", func() {
				var key *int64
				val := int64(12345)

				utils.AssignTo(key, val)
				Expect(key).To(BeZero())
			})
		})
	})
})
