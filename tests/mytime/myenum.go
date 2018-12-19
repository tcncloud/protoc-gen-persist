package mytime

import (
	"cloud.google.com/go/spanner"
	"database/sql/driver"
	"fmt"
	"github.com/tcncloud/protoc-gen-persist/tests/test"
)

type MyEnum struct {
	Status int32
}

func (t MyEnum) ToSql(src test.TestEnum) *MyEnum {
	return t.ToSpanner(src)
}
func (t MyEnum) ToSpanner(src test.TestEnum) *MyEnum {
	t.Status = int32(src)

	return &t
}

func (t MyEnum) ToProto() test.TestEnum {
	return test.TestEnum(t.Status)
}

func (t *MyEnum) Scan(src interface{}) error {
	ti, ok := src.(int32)
	if !ok {
		return fmt.Errorf("could not scan out enum")
	}
	t.Status = ti
	return nil
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
func (t *MyEnum) Value() (driver.Value, error) {
	return t.Status, nil
}
func (t *MyEnum) SpannerValue() (interface{}, error) {
	return int64(t.Status), nil
}
