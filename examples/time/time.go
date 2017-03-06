package time

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type Time struct {
	Seconds int64
	Nanos   int32
}

func (s *Time) Scan(src interface{}) error {
	switch src.(type) {
	case time.Time:
		t := src.(time.Time)
		s.Nanos = int32(t.Nanosecond())
		s.Seconds = t.Unix()
	case *time.Time:
		t := src.(*time.Time)
		s.Nanos = int32(t.Nanosecond())
		s.Seconds = t.UnixNano()
	default:
		return fmt.Errorf("Can't convert from %+v!", src)
	}
	return nil
}

func (s Time) Value() (driver.Value, error) {
	return driver.Value(time.Unix(s.Seconds, int64(s.Nanos))), nil

}

func (s Time) Get() *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Nanos:   s.Nanos,
		Seconds: s.Seconds,
	}
}
