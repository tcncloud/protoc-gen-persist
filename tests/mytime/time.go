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

package mytime

import (
	"cloud.google.com/go/spanner"
	"database/sql/driver"
	"github.com/golang/protobuf/ptypes/timestamp"
	"strconv"
	"strings"
)

type MyTime struct {
	Seconds int64
	Nanos   int32
}

func (s MyTime) ToSql(src *timestamp.Timestamp) *MyTime {
	s.Seconds = src.Seconds
	s.Nanos = src.Nanos
	return &s
}

func (s MyTime) ToSpanner(src *timestamp.Timestamp) *MyTime {
	s.Seconds = src.Seconds
	s.Nanos = src.Nanos
	return &s
}

func (s MyTime) ToProto() *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Nanos:   s.Nanos,
		Seconds: s.Seconds,
	}

}

func (t *MyTime) Scan(src interface{}) error {
	ti, ok := src.(string)
	if !ok {
		t.Seconds = int64(0)
		t.Nanos = 0
	}
	tis := strings.Split(ti, ",")
	secs, err := strconv.ParseInt(tis[0], 10, 64)
	if err != nil {
		return err
	}
	nans, err := strconv.ParseInt(tis[1], 10, 32)
	if err != nil {
		return err
	}
	t.Seconds = secs
	t.Nanos = int32(nans)

	return nil
}

func (t *MyTime) Value() (driver.Value, error) {
	ti := strconv.FormatInt(t.Seconds, 10) + "," + strconv.FormatInt(int64(t.Nanos), 10)
	return ti, nil
}

func (t *MyTime) SpannerScan(src *spanner.GenericColumnValue) error {
	var strTime string
	err := src.Decode(&strTime)
	if err != nil {
		return err
	}
	return t.Scan(strTime)
}

func (t *MyTime) SpannerValue() (interface{}, error) {
	return t.Value()
}
