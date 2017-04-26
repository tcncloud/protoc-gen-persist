// test before and after hooks

package mytime

import (
	"fmt"
	pb "github.com/tcncloud/protoc-gen-persist/examples/test"
)

var cache map[string]interface{}

func init() {
	cache = make(map[string]interface{})
}

func UniarySelectBeforeHook(req *pb.PartialTable) (*pb.ExampleTable, error) {
	fmt.Printf("UniarySelectBeforeHook: %+v\n", req)
	return nil, nil
}

func UniarySelectAfterHook(res *pb.ExampleTable) error {
	fmt.Printf("uniarySelectAfterHook: %+v\n", res)
	return nil
}

func ServerStreamBeforeHook(req *pb.Name) (*pb.ExampleTable, error) {
	fmt.Printf("ServerStreamBeforeHook: %+v\n", req)
	return nil, nil
}

func ServerStreamAfterHook(res *pb.ExampleTable) error {
	fmt.Printf("ServerStreamAfterHook: %+v\n", res)
	return nil
}

func BidirectionalBeforeHook(req *pb.ExampleTable) (*pb.ExampleTable, error) {
	fmt.Printf("BidirectionalBeforeHook: %+v\n", req)
	return nil, nil
}

func BidirectionalAfterHook(res *pb.ExampleTable) error {
	fmt.Printf("BidirectionalAfterHook: %+v\n", res)
	return nil
}

func ClientStreamBeforeHook(req *pb.ExampleTable) (*pb.NumRows, error) {
	fmt.Printf("ClientStreamBeforeHook: %+v\n", req)
	return nil, nil
}
