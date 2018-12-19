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

package main

import (
	"io/ioutil"
	"os"

	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/sirupsen/logrus"
	"github.com/tcncloud/protoc-gen-persist/generator"
)

func init() {
	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	if len(os.Args) > 1 {
		fmt.Println("This executable is meant to be used by protoc!\nGo to http://github.com/tcncloud/protoc-gen-persist for more info")
		os.Exit(-1)
	}
	var req plugin_go.CodeGeneratorRequest

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logrus.Fatal("Can't read the stdin!")
	}
	if err := proto.Unmarshal(data, &req); err != nil {
		logrus.Fatal("Error parsing data!")
	}
	g := generator.NewGenerator(&req)
	err = g.Process()
	if err != nil {
		logrus.Fatalf("error processing generator: %s", err)
		return
	}

	// Send back the results.
	resp, err := g.GetResponse()
	if err != nil {
		logrus.Fatalf("recieved err getting the file response: %s", err)
	}
	logrus.Debugf("file length: %d\n", len(resp.File))
	data, err = proto.Marshal(resp)
	if err != nil {
		logrus.Fatal("I can't serialize response")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		logrus.Fatal("Can't send data to stdout!")
	}

}
