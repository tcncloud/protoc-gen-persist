package little_of_everything

import (
	"database/sql/driver"
	"fmt"
)

type MyMappedEnum struct {
	Status int32
}

func (t MyMappedEnum) ToSql(src MappedEnum) *MyMappedEnum {
	t.Status = int32(src)
	return &t
}

func (t MyMappedEnum) ToProto() MappedEnum {
	return MappedEnum(t.Status)
}

func (t *MyMappedEnum) Scan(src interface{}) error {
	ti, ok := src.(int32)
	if !ok {
		return fmt.Errorf("could not scan out enum")
	}
	t.Status = ti
	return nil
}
func (t *MyMappedEnum) Value() (driver.Value, error) {
	return t.Status, nil
}
