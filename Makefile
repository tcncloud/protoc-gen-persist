
all: build

generate: deps
	protoc -I/usr/local/include -I. -I$$GOPATH/src --go_out=Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:$$GOPATH/src ./persist/options.proto

build: generate
	go build 

install: build
	go install


deps: $(GOPATH)/bin/protoc-gen-go \
	$(GOPATH)/bin/ginkgo 


$(GOPATH)/bin/protoc-gen-grpc-gateway:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

$(GOPATH)/bin/protoc-gen-swagger:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

$(GOPATH)/bin/protoc-gen-go:
	go get -u github.com/golang/protobuf/protoc-gen-go

$(GOPATH)/bin/ginkgo:	
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega  
