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
		persist/*.proto

proto-examples:
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		examples/sql/basic/*.proto
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		examples/sql/little_of_everything/*.proto
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		examples/spanner/basic/*.proto
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		examples/spanner/delete/*.proto
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		examples/spanner/bob_example/*.proto
	protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
		examples/test/*.proto

build: generate
	glide install
	go build

install: build
	go install

test: deps build
	ginkgo -r

test-compile:
	go build
	DEBUG=true protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--plugin=./protoc-gen-persist \
		--persist_out=.  examples/sql/little_of_everything/*.proto
	DEBUG=true protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--plugin=./protoc-gen-persist \
		--persist_out=.  examples/sql/basic/*.proto
	DEBUG=true protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--plugin=./protoc-gen-persist \
		--persist_out=.  examples/spanner/basic/*.proto
	DEBUG=true protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--plugin=./protoc-gen-persist \
		--persist_out=.  examples/spanner/delete/*.proto
	DEBUG=true protoc -I/usr/local/include -I. -I$$GOPATH/src \
		--plugin=./protoc-gen-persist \
		--persist_out=.  examples/spanner/bob_example/*.proto


test-sql-impl: build
	env GOOS=linux go build -o ./test-impl/server.main ./test-impl/server/sql
	env GOOS=linux go build -o ./test-impl/client.main ./test-impl/client/sql

test-spanner-impl: build
	go build -o ./test-impl/server.main ./test-impl/server/spanner/basic
	go build -o ./test-impl/client.main ./test-impl/client/spanner/basic

test-bobs: build
	go build -o ./test-impl/server.main ./test-impl/server/spanner/bobs
	go build -o ./test-impl/client.main ./test-impl/client/spanner/bobs

deps: $(GOPATH)/bin/protoc-gen-go $(GOPATH)/bin/ginkgo  $(GOPATH)/bin/glide


$(GOPATH)/bin/protoc-gen-go:
	go get -u github.com/golang/protobuf/protoc-gen-go

$(GOPATH)/bin/ginkgo:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega

$(GOPATH)/bin/glide:
	go get -u github.com/Masterminds/glide

clean:
	rm -f examples/*.pb.go examples/*.persist.go examples/test/*.pb.go
