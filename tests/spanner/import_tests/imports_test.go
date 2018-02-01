package import_tests_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	//"github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tcncloud/protoc-gen-persist/generator"
)

func TestNestedService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Import Test Suite")
}

var _ = Describe("when processing imports", func() {
	_ = Describe("When go_package and persist.package are defined and the same", func() {
		It("generates service_impl in same package as service.pb.go definition", func() {
			out := ProcessFile("persist_and_go_same.pb")
			fmt.Printf("out? %+v\n", *out)
			fmt.Printf("len of files: %d\n", len(out.File))
		})
		It("does not import the service.pb.go file", func() {
		})
	})
	_ = Describe("When go_package and persist.package are defined but different", func() {
		_ = ProcessFile("persist_and_go.pb")
		It("generates service_impl in correct spot", func() {
		})
		It("imports the service.pb.go definitions", func() {
		})
	})
	_ = Describe("When only persist is defined", func() {
		_ = ProcessFile("persist_no_go.pb")
		It("imports the service.pb.go file", func() {
		})
	})
	_ = Describe("When only go_package is defined", func() {
		_ = ProcessFile("go_no_persist.pb")
		It("generates service impl in same package as service.pb.go definition", func() {
		})
		It("does not import the service.pb.go file", func() {
		})
	})
	_ = Describe("When neither are defined, and package is simple", func() {
		_ = ProcessFile("neither_simple_package.pb")
		It("generates correct package name for impl", func() {
		})
		It("does not import the service.pb.go file", func() {
		})
	})
	_ = Describe("When neither are defined and package is nested", func() {
		_ = ProcessFile("neither_nested_package.pb")
		It("generates correct package name for impl", func() {
		})
		It("does not import the service.pb.go file", func() {
		})
	})

})

func ProcessFile(loc string) *plugin.CodeGeneratorResponse {
	//logrus.SetLevel(logrus.DebugLevel)
	var req plugin.CodeGeneratorRequest
	file, err := os.Open(loc)
	if err != nil {
		panic(fmt.Sprintf("could not open: %s, err: %s", loc, err))
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("could not read file at: %s, err: %s", loc, err))
	}
	fmt.Printf("len of data: %d\n", len(data))
	if len(data) == 0 {
		panic("no data read from file")
	}
	var f google_protobuf.FileDescriptorSet
	if err := proto.Unmarshal(data, &f); err != nil {
		panic(fmt.Sprintf("could not Unmarshal file: %s", err))
	}
	req.FileToGenerate = append(req.FileToGenerate, loc)
	req.ProtoFile = f.File
	//fmt.Printf("files to generate: %+v\n", req.FileToGenerate)
	fmt.Printf("Protofile len: %d\n", len(req.ProtoFile))
	g := generator.NewGenerator(&req)
	if err := g.Process(); err != nil {
		panic(fmt.Sprintf("could not process the generator for %s, err: %s", loc, err))
	}
	resp, err := g.GetResponse()
	if err != nil {
		panic(fmt.Sprintf("could not get the protobuf response for %s, err: %s", loc, err))
	}
	if resp == nil {
		panic("response was nil")
	}
	fmt.Printf("finished generating response for: %s\n", loc)
	return resp
}
