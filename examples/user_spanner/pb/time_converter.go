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

package pb

import (
	"time"

	"cloud.google.com/go/spanner"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type TimeString struct {
	t *timestamp.Timestamp
}

func (ts TimeString) ToSpanner(t *timestamp.Timestamp) MappingImpl_UServ_TimestampTimestamp {
	ts.t = t
	return &ts
}
func (ts TimeString) ToProto(req **timestamp.Timestamp) error {
	*req = ts.t
	return nil
}

func (t *TimeString) SpannerScan(src *spanner.GenericColumnValue) error {
	var tStr string
	if err := src.Decode(&tStr); err != nil {
		return err
	}
	ti, err := time.Parse(time.RFC3339, tStr)
	if err != nil {
		return err
	}
	stamp, err := ptypes.TimestampProto(ti)
	if err != nil {
		return err
	}
	t.t = stamp
	return nil
}

func (t *TimeString) SpannerValue() (interface{}, error) {
	return ptypes.TimestampString(t.t), nil
}
func (t TimeString) Empty() MappingImpl_UServ_TimestampTimestamp {
	return new(TimeString)
}

type SliceStringConverter struct {
	v *SliceStringParam
}

func (s *SliceStringConverter) ToSpanner(v *SliceStringParam) MappingImpl_UServ_SliceStringParam {
	s.v = v
	return s
}
func (s *SliceStringConverter) ToProto(req **SliceStringParam) error {
	*req = s.v
	return nil
}

func (s *SliceStringConverter) SpannerScan(src *spanner.GenericColumnValue) error {
	values := make([]string, 0)
	if err := src.Decode(values); err != nil {
		return err
	}

	s.v = &SliceStringParam{Slice: []string(values)}

	return nil
}

func (s *SliceStringConverter) SpannerValue() (interface{}, error) {
	return s.v.GetSlice(), nil
}

func (s *SliceStringConverter) Empty() MappingImpl_UServ_SliceStringParam {
	return new(SliceStringConverter)
}

var inc int64

func IncId(u *User) ([]*User, error) {
	u.Id = inc
	inc++
	return nil, nil
}
