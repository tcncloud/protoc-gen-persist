package generator

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

type Generator struct {
	Files           []*string
	OriginalRequest *plugin_go.CodeGeneratorRequest
}

func NewGenerator(request *plugin_go.CodeGeneratorRequest) *Generator {
	ret := new(Generator)
	ret.OriginalRequest = request
	return ret
}

// check if a service has at least one method that has the persist.ql extension defined
func IsServicePersistEnabled(service *descriptor.ServiceDescriptorProto) bool {
	if service.Method != nil {
		for _, method := range service.Method {
			if method.GetOptions() != nil {
				if proto.HasExtension(method.Options, persist.E_Ql) {
					// at least one method implement persist.ql
					return true
				}
			}
		}
	}
	return false
}

// Process the request
func (g *Generator) ProcessRequest() {
	for _, file := range g.OriginalRequest.ProtoFile {
		for _, service := range file.Service {
			for _, method := range service.Method {
				method.GetOptions()
			}
		}
	}
}
