package structures_test

import (
	. "github.com/tcncloud/protoc-gen-persist/generator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/proto"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var _ = Describe("Struct", func() {
	Describe("IsInnerType()", func() {
		It("should return true when ParentDescriptor field is not null", func() {
			s := &Struct{
				MessageDescriptor: &desc.DescriptorProto{
					Name: proto.String("test"),
				},
				ParentDescriptor: &desc.DescriptorProto{
					Name: proto.String("test_test"),
				},
			}
			Expect(s.IsInnerType()).To(Equal(true))
		})
	})
	Describe("GetGoName()", func() {
		Describe("If Struct is a message with name test_test", func() {
			It("should return TestTest", func() {
				s := &Struct{
					MessageDescriptor: &desc.DescriptorProto{
						Name: proto.String("test_test"),
					},
				}
				Expect(s.GetGoName()).To(Equal("TestTest"))
			})
		})
		Describe("If Struct is an inner message named test_test to the message test", func() {
			It("should return Test_TestTest", func() {
				s := &Struct{
					MessageDescriptor: &desc.DescriptorProto{
						Name: proto.String("test_test"),
					},
					ParentDescriptor: &desc.DescriptorProto{
						Name: proto.String("test"),
					},
				}
				Expect(s.GetGoName()).To(Equal("Test_TestTest"))

			})
		})
		Describe("If Struct is an inner enum named test_test to the message test", func() {
			It("should return Test_TestTest", func() {
				s := &Struct{
					EnumDescriptor: &desc.EnumDescriptorProto{
						Name: proto.String("test_test"),
					},
					ParentDescriptor: &desc.DescriptorProto{
						Name: proto.String("test"),
					},
				}
				Expect(s.GetGoName()).To(Equal("Test_TestTest"))

			})
		})
	})
})

var _ = Describe("StructList", func() {})
