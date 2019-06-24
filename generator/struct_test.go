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
	"reflect"
	"testing"
)

func Test_compareProtoName(t *testing.T) {
	type args struct {
		name      string
		protoname string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: ".package.Name",
			args: args{
				name:      "package.Name",
				protoname: ".package.Name",
			},
			want: true,
		},
		{
			name: "+package.Name",
			args: args{
				name:      "package.Name",
				protoname: "+package.Name",
			},
			want: false,
		},
		{
			name: "fail",
			args: args{
				name:      "test",
				protoname: "t",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareProtoName(tt.args.name, tt.args.protoname); got != tt.want {
				t.Errorf("compareProtoName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructList_GetStructByProtoName(t *testing.T) {
	type args struct {
		name string
	}
	t1 := &Struct{
		ProtoName: "test.Test",
	}
	t2 := &Struct{
		GoName: "I have no idea",
	}
	t3 := &Struct{
		ProtoName: "xxxxx",
	}
	// predefined struxtures
	s := &StructList{t1, t2, t3}
	tests := []struct {
		name string
		s    *StructList
		args args
		want *Struct
	}{
		{
			name: "simple test",
			s:    s,
			args: args{
				name: "test.Test",
			},
			want: t1,
		},
		{
			name: "simple fail",
			s:    s,
			args: args{
				name: ".test.Test",
			},
			want: nil,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetStructByProtoName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructList.GetStructByProtoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
