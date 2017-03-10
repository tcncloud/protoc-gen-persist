package time

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type Time struct {
	Seconds int64
	Nanos   int32
}

func (s Time) ToSql(src *timestamp.Timestamp) time.Time {
	return time.Unix(src.Seconds, int64(src.Nanos))
}
func (s Time) FromSql() *timestamp.Timestamp {
	return &timestamp.Timestamp{
		Nanos:   s.Nanos,
		Seconds: s.Seconds,
	}

}
