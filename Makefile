# Copyright 2017, TCN Inc.
# All rights reserved.

# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are
# met:

#     * Redistributions of source code must retain the above copyright
# notice, this list of conditions and the following disclaimer.
#     * Redistributions in binary form must reproduce the above
# copyright notice, this list of conditions and the following disclaimer
# in the documentation and/or other materials provided with the
# distribution.
#     * Neither the name of TCN Inc. nor the names of its
# contributors may be used to endorse or promote products derived from
# this software without specific prior written permission.

# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
# "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
# LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
# A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
# OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
# SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
# LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
# DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
# THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
# (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


PROTO_FILES:= persist/options.proto examples/example1.proto

all: build

generate: deps proto-persist proto-examples

proto-persist:
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		persist/options.proto

proto-examples:
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		examples/example1.proto

build: generate 
	glide install
	go build 

install: build
	go install

test: deps
	ginkgo -r 

test-compile: build
	DEBUG=true protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--plugin=./protoc-gen-persist \
		--persist_out=. \
		examples/*.proto

deps: $(GOPATH)/bin/protoc-gen-go $(GOPATH)/bin/ginkgo  $(GOPATH)/bin/glide


$(GOPATH)/bin/protoc-gen-go:
	go get -u github.com/golang/protobuf/protoc-gen-go

$(GOPATH)/bin/ginkgo:	
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega  

$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide