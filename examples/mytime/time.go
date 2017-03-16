package mytime

import (
	"time"
	"database/sql/driver"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type MyTime struct {
	Seconds int64
	Nanos   int32
}

func (s MyTime) ToSql(src *timestamp.Timestamp) *MyTime {
	s.Seconds = src.Seconds
	s.Nanos = src.Nanos
	return &s
}
func (s MyTime) ToProto() *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Nanos:   s.Nanos,
		Seconds: s.Seconds,
	}

}
func (t *MyTime) Scan(src interface{}) error {
	ti, ok := src.(time.Time)
	if !ok {
		t.Seconds = int64(0)
		t.Nanos = 0
	} else {
		t.Seconds = ti.Unix()
		t.Nanos = int32(ti.UnixNano())
	}
	return nil
}

func (t *MyTime) Value() (driver.Value, error) {
	ti := time.Unix(t.Seconds, int64(t.Nanos))
	return ti, nil
}
