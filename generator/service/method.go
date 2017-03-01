package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/tcncloud/protoc-gen-persist/persist"
)

type Method struct {
	Desc    *descriptor.MethodDescriptorProto
	Service *Service
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

func (m *Method) GetName() string {
	return m.Desc.GetName()
}

func (m *Method) IsEnabled() bool {
	if m.GetMethodOption() != nil {
		return true
	}
	return false
}

func (m *Method) IsSQL() bool {
	if opt := m.GetMethodOption(); opt != nil {
		return opt.GetPersist() == persist.PersistenceOptions_SQL
	}
	return false
}

func (m *Method) IsMongo() bool {
	if opt := m.GetMethodOption(); opt != nil {
		return opt.GetPersist() == persist.PersistenceOptions_MONGO
	}
	return false
}

func (m *Method) IsSpanner() bool {
	if opt := m.GetMethodOption(); opt != nil {
		return opt.GetPersist() == persist.PersistenceOptions_SPANNER
	}
	return false
}

func (m *Method) IsUnary() bool {
	return !m.Desc.GetClientStreaming() && !m.Desc.GetServerStreaming()
}

func (m *Method) IsClientStreaming() bool {
	return m.Desc.GetClientStreaming() && !m.Desc.GetServerStreaming()
}

func (m *Method) IsServerStreaming() bool {
	return !m.Desc.GetClientStreaming() && m.Desc.GetServerStreaming()
}
func (m *Method) IsBidiStreaming() bool {
	return m.Desc.GetClientStreaming() && m.Desc.GetServerStreaming()
}

type Methods []*Method

func (m *Methods) AddMethod(desc *descriptor.MethodDescriptorProto, service *Service) {
	*m = append(*m, &Method{Desc: desc, Service: service})
}
