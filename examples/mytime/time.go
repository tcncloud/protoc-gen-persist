package mytime

import (
	"strings"
	"strconv"
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

func (s MyTime) ToSpanner(src *timestamp.Timestamp) *MyTime {
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
	ti, ok := src.(string)
	if !ok {
		t.Seconds = int64(0)
		t.Nanos = 0
	}
	tis := strings.Split(ti, ",")
	secs, err := strconv.ParseInt(tis[0], 10, 64)
	if err != nil {
		return err
	}
	nans, err := strconv.ParseInt(tis[1], 10, 32)
	if err != nil {
		return err
	}
	t.Seconds = secs
	t.Nanos = int32(nans)

	return nil
}

func (t *MyTime) Value() (driver.Value, error) {
	ti := strconv.FormatInt(t.Seconds, 10) + "," + strconv.FormatInt(int64(t.Nanos), 10)
	return ti, nil
}

func (t *MyTime) SpannerScan(src *spanner.GenericColumnValue) error {
	return nil
}

func (t *MyTime) SpannerValue() (interface{}, error) {
	return nil, nil
}
