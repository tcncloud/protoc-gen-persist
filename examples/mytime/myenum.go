package mytime

import (
	"cloud.google.com/go/spanner"
	"github.com/tcncloud/protoc-gen-persist/examples/test"
)

type MyEnum struct {
	Status int32
}

func (t MyEnum) ToSpanner(src test.TestEnum) *MyEnum{
	t.Status = int32(src)

	return &t
}

func (t MyEnum) ToProto() test.TestEnum {
	return test.TestEnum(t.Status)
}

func (t *MyEnum) SpannerScan(src *spanner.GenericColumnValue) error {
	var lt int64

	err := src.Decode(&lt)
	if err != nil {
		return err
	}
	t.Status = int32(lt)

	return nil
}

func (t *MyEnum) SpannerValue() (interface{}, error) {
	return int64(t.Status), nil
}


