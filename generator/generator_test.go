package generator_test

import (
	"testing"

	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/tcncloud/protoc-gen-persist/generator"
)

func TestCommandLineParameter(t *testing.T) {
	var req plugin_go.CodeGeneratorRequest
	g := generator.NewGenerator(&req)

	g.CommandLineParameters("paths=source_relative")

	if (g.SourceRelative == false && g.ImportPaths == false) || (g.SourceRelative == false && g.ImportPaths == true) {
		t.Error("error parsing the paths=source_relative")
	}

	g.CommandLineParameters("paths=imports")
	if (g.SourceRelative == true && g.ImportPaths == true) || (g.SourceRelative == true && g.ImportPaths == false) {
		t.Error("error parsing the paths=source_relative")
	}

}
