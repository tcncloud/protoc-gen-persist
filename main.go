package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

func Return(response *plugin_go.CodeGeneratorResponse) {
	data, err := proto.Marshal(response)
	if err != nil {
		log.Fatal("That's wired ... ")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		log.Fatal(os.Stderr, "I can't send data to stdout !")
	}
}

func main() {
	var req plugin_go.CodeGeneratorRequest
	var res *plugin_go.CodeGeneratorResponse

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		res.Error = proto.String("Can't read the input")
		Return(res)
		return
	}

	if err := proto.Unmarshal(data, &req); err != nil {
		res.Error = proto.String("Error parsing stdin data")
		Return(res)
		return
	}
	// DO processing

	// Send back the results.
	Return(res)
}
