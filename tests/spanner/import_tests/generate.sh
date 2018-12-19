#!/bin/sh

protoc --include_imports -I. -I$GOPATH/src \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist --include_source_info \
			 --descriptor_set_out persist_and_go_same.pb persist_and_go_same.proto
protoc --include_imports -I. -I$GOPATH/src \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist --include_source_info \
			 --descriptor_set_out persist_and_go.pb persist_and_go.proto
protoc --include_imports -I. -I$GOPATH/src \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist --include_source_info \
			 --descriptor_set_out persist_no_go.pb persist_no_go.proto
protoc --include_imports -I. -I$GOPATH/src \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist --include_source_info \
			 --descriptor_set_out go_no_persist.pb go_no_persist.proto
protoc --include_imports -I. -I$GOPATH/src \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist --include_source_info \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist/examples/test \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist/examples \
			 --descriptor_set_out neither_simple_package.pb neither_simple_package.proto
protoc --include_imports -I. -I$GOPATH/src \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist --include_source_info \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist/examples/test \
			 -I$GOPATH/src/github.com/tcncloud/protoc-gen-persist/examples \
			 --descriptor_set_out neither_nested_package.pb neither_nested_package.proto


