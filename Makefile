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


PROTO_FILES:= persist/options.proto tests/example1.proto
PROTOC_DIR?=/usr/local

PROTOC_INCLUDE:=$(PROTOC_DIR)/include

PROTOC:=$(PROTOC_DIR)/bin/protoc

all: build

generate: deps proto-persist proto-examples

proto-persist:
	$(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
		--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. \
		persist/*.proto

proto-examples:
	 $(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
	 	--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
	 	tests/spanner/basic/*.proto
	 $(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
	 	--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
	 	tests/sql/little_of_everything/*.proto
	 $(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
	 	--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
	 	tests/sql/basic/*.proto
	 $(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
	 	--go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src \
	 	tests/test/*.proto

build: generate
	go mod download
	go build

install: build
	go install

test: deps build
	ginkgo -r

test-compile:
	go build
	DEBUG=false $(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
	 	--plugin=./protoc-gen-persist \
	 	--persist_out=persist_root=tests/sql/little_of_everything:.  tests/sql/little_of_everything/*.proto
	DEBUG=false $(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
	 	--plugin=./protoc-gen-persist \
	 	--persist_out=persist_root=tests/sql/basic:.  tests/sql/basic/*.proto
	DEBUG=false $(PROTOC) -I$(PROTOC_INCLUDE) -I. -I$$GOPATH/src \
		--plugin=./protoc-gen-persist \
		--persist_out=persist_root=tests/spanner/basic:.  tests/spanner/basic/*.proto
	cd ./tests/sql/little_of_everything && go build
	cd ./tests/sql/basic && go build
	cd ./tests/spanner/basic && go build


deps: $(GOPATH)/bin/protoc-gen-go $(GOPATH)/bin/ginkgo  


$(GOPATH)/bin/protoc-gen-go:
	go get -u github.com/golang/protobuf/protoc-gen-go

$(GOPATH)/bin/ginkgo:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega


clean:
	go clean -modcache
	rm -f tests/*.pb.go tests/*.persist.go tests/test/*.pb.go
	rm -f tests/spanner/bob_example/*.pb.go tests/spanner/bob_example/*.persist.go
	rm -f tests/spanner/basic/*.pb.go tests/spanner/basic/*.persist.go
	rm -rf tests/spanner/basic/persist_lib
