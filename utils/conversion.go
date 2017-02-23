package utils

import(
	"time"
	"database/sql"
	"strings"
	"github.com/lib/pq"
	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/ptypes"
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func ToSafeType(in interface{}) interface{} {
	switch t := in.(type) {
	case []string:
		return (*pq.StringArray)(&t)
	case []int64:
		return (*pq.Int64Array)(&t)
	case []float64:
		return (*pq.Float64Array)(&t)
	case *google_protobuf.Timestamp:
		return (*time.Time)(ToTime(t))
	case string:
		return (*string)(&t)
	case int64:
		return (*int64)(&t)
	case int32:
		return (*int32)(&t)
	case float64:
		return (*float64)(&t)
	default:
		return in
	}
}

// dereferences key  and assigns val into key, based on val's  type.
// used to assign into protobuf fields without having to do a lot of type checking
// key should be a pointer, val should be the actual value to be assigned
// returns true if key was assigned,  otherwise returns false
func AssignTo(key interface{}, val interface{}) bool {
	switch t := val.(type) {
	case pq.StringArray:
		if k, ok := key.(*[]string); ok {
			v := ([]string)(t)
			if k != nil {
				*k = v
				return true
			}
		}
	case pq.Int64Array:
		if k, ok := key.(*[]int64); ok {
			v := ([]int64)(t)
			if k != nil {
				*k = v
				return true
			}
		}
	case pq.Float64Array:
		if k, ok := key.(*[]float64); ok {
			v := ([]float64)(t)
			if k != nil {
				*k = v
				return true
			}
		}
	case int64:
		if k, ok := key.(*int64); ok {
			v := (int64)(t)
			if k != nil {
				*k = v
				return true
			}
		}
	case int32:
		if k, ok := key.(*int32); ok {
			v := (int32)(t)
			if k != nil {
				*k = v
				return true
			}
		}
	case float64:
		if k, ok := key.(*float64); ok {
			v := (float64)(t)
			if k != nil {
				*k = v
				return true
			}
		}
	case string:
		if k, ok := key.(*string); ok {
			v := (string)(t)
			if k != nil {
				*k = v
				return true
			}
		}
	case *time.Time:
		if k, ok := key.(**google_protobuf.Timestamp); ok {
			v := ToProtobufTime(t)
			if k != nil {
				*k = v
				return true
			}
		}
	default:
		logrus.WithField("key", key).WithField("val", val).Warn("val contained an unknown type")
	}
	logrus.WithField("key", key).Warn("could not assign to key, doing nothing!")
	return false
}

func ToTime(entry *google_protobuf.Timestamp) *time.Time {
	if entry == nil {
		return nil
	}
	lTime, err := ptypes.Timestamp(entry)
	if err != nil {
		logrus.WithError(err).Error("something went wrong on timestamp conversion!")
		return nil
	}
	return &lTime
}

func ToProtobufTime(lTime *time.Time) *google_protobuf.Timestamp {
	res, err := ptypes.TimestampProto(*lTime)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"time.Time": lTime,
			"error":     err,
		}).Info("threw error when converting to protobuf timestamp")

		return nil
	}
	return res
}

func ConvertError(err error, req interface{}) error {
	if err == sql.ErrNoRows {
		return grpc.Errorf(codes.NotFound, "%+v doesnt exist", req)
	} else if strings.Contains(err.Error(), "duplicate key") {
		return grpc.Errorf(codes.AlreadyExists, "%+v already exists")
	}
	return grpc.Errorf(codes.Unknown, err.Error())
}
