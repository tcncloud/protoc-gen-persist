package basic

import (
	"cloud.google.com/go/spanner"
)

type MyMappedEnum struct {
	Status int32
}

func (t MyMappedEnum) ToSpanner(src MyEnum) *MyMappedEnum {
	t.Status = int32(src)
	return &t
}

func (t MyMappedEnum) ToProto() MyEnum {
	return MyEnum(t.Status)
}

func (t *MyMappedEnum) SpannerScan(src *spanner.GenericColumnValue) error {
	var lt int64

	err := src.Decode(&lt)
	if err != nil {
		return err
	}
	t.Status = int32(lt)

	return nil
}
func (t *MyMappedEnum) SpannerValue() (interface{}, error) {
	return int64(t.Status), nil
}
