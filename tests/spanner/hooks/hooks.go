// test before and after hooks

package hooks

import (
	"fmt"
	pb "github.com/tcncloud/protoc-gen-persist/tests/test"
)

var cache map[int64]*pb.ExampleTable

func init() {
	cache = make(map[int64]*pb.ExampleTable)
}

func UniaryInsertBeforeHook(req *pb.ExampleTable) (*pb.ExampleTable, error) {
	fmt.Printf("UniaryInsertBeforeHook: %+v\n", *req)
	return nil, nil
}

func UniaryInsertAfterHook(req *pb.ExampleTable, res *pb.ExampleTable) error {
	fmt.Printf("UniaryInsertAfterHook req:%+v res:%+v\n", *req, *res)
	return nil
}

func UniaryUpdateBeforeHook(req *pb.ExampleTable) (*pb.PartialTable, error) {
	fmt.Printf("UniaryUpdateBeforeHook: %+v\n", *req)
	return nil, nil
}

func UniaryUpdateAfterHook(req *pb.ExampleTable, res *pb.PartialTable) error {
	fmt.Printf("UniaryUpdateAfterHook req:%+v res:%+v\n", *req, *res)
	return nil
}

func UniaryDeleteBeforeHook(req *pb.ExampleTableRange) (*pb.ExampleTable, error) {
	fmt.Printf("UniaryDeleteBeforeHook: %+v\n", *req)
	return nil, nil
}

func UniaryDeleteAfterHook(req *pb.ExampleTableRange, res *pb.ExampleTable) error {
	fmt.Printf("UniaryDeleteAfterHook req:%+v res:%+v\n", *req, *res)
	return nil
}

func ServerStreamBeforeHook(req *pb.Name) ([]*pb.ExampleTable, error) {
	fmt.Printf("ServerStreamBeforeHook: %+v\n", *req)
	return nil, nil
}

func ServerStreamAfterHook(req *pb.Name, res *pb.ExampleTable) error {
	fmt.Printf("ServerStreamAfterHook req:%+v res:%+v\n", *req, *res)
	return nil
}

func ClientStreamUpdateBeforeHook(req *pb.ExampleTable) (*pb.NumRows, error) {
	fmt.Printf("ClientStreamUpdateBeforeHook: %+v\n", *req)
	return nil, nil
}

func ClientStreamUpdateAfterHook(req *pb.ExampleTable, res *pb.NumRows) error {
	fmt.Printf("ClientStreamAfterHook req:%+v res:%+v\n", *req, *res)
	return nil
}
