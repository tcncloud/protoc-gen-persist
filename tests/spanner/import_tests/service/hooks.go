package service

import (
	"cloud.google.com/go/spanner"
	it "github.com/tcncloud/protoc-gen-persist/examples/spanner/import_tests"
	t "github.com/tcncloud/protoc-gen-persist/examples/test"
)

func TestBeforeHook(req *t.ExampleTable) (*it.ExampleTable, error) {
	return nil, nil
}
func TestAfterHook(req *t.ExampleTable, res *it.ExampleTable) error {
	return nil
}

func AggRequests(req *t.ExampleTable, res *it.AggExampleTables) error {
	return nil
}

type MySpannerImpl struct {
	SpannerDB *spanner.Client
}
