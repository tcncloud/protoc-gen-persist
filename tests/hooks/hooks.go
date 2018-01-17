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

func UniarySelectBeforeHook(req *pb.PartialTable) (*pb.ExampleTable, error) {
	fmt.Printf("UniarySelectBeforeHook: cache: %#v\n", cache)
	if req != nil {
		res := cache[req.Id]
		fmt.Printf("UniarySelectBeforeHook: req:%+v , cache: %+v\n", *req, res)
		return res, nil
	} else {
		fmt.Println("UniarySelectBeforeHook: req was nil...")
	}
	return nil, nil
}

func UniarySelectAfterHook(req *pb.PartialTable, res *pb.ExampleTable) error {
	if req != nil {
		fmt.Printf("uniarySelectAfterHook: req:%+v , res:%+v\n", *req, *res)
		cache[req.Id] = res
	}
	fmt.Printf("UniarySelectAfterHook: cache now: %#v\n", cache)
	return nil
}

func ServerStreamBeforeHook(req *pb.Name) ([]*pb.ExampleTable, error) {
	fmt.Printf("ServerStreamBeforeHook: %+v\n", req)
	return nil, nil
}

func ServerStreamAfterHook(req *pb.Name, res *pb.ExampleTable) error {
	fmt.Printf("ServerStreamAfterHook: %+v\n", res)
	return nil
}

func BidirectionalBeforeHook(req *pb.ExampleTable) (*pb.ExampleTable, error) {
	fmt.Printf("BidirectionalBeforeHook: %+v\n", req)
	return nil, nil
}

func BidirectionalAfterHook(req *pb.ExampleTable, res *pb.ExampleTable) error {
	fmt.Printf("BidirectionalAfterHook: %+v\n", res)
	return nil
}

func ClientStreamBeforeHook(req *pb.ExampleTable) (*pb.NumRows, error) {
	fmt.Printf("ClientStreamBeforeHook: %+v\n", req)
	return nil, nil
}

func ClientStreamAfterHook(req *pb.ExampleTable, res *pb.Ids) error {
	if res != nil {
		fmt.Printf("ClientStreamAfterHook adding to Ids. so far collected: %+v\n", *res)
		res.Ids = append(res.Ids, req.Id)
	} else {
		fmt.Println("ClientStreamingAfterHook res was nil...")
	}
	return nil
}
