package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

type Method struct {
	Desc *descriptor.MethodDescriptorProto
}

func (m *Method) GetMethodOption() *persist.QLImpl {
	if m.Desc.Options != nil && proto.HasExtension(m.Desc.Options, persist.E_Ql) {
		ext, err := proto.GetExtension(m.Desc.Options, persist.E_Ql)
		if err == nil {
			return ext.(*persist.QLImpl)
		}
	}
	return nil
}

type Methods []*Method

func (m *Methods) AddMethod(desc *descriptor.MethodDescriptorProto) {
	*m = append(*m, &Method{Desc: desc})
}
