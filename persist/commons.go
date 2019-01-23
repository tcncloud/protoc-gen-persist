package persist

import (
	"context"
	"database/sql"

	"cloud.google.com/go/spanner"
)

type PersistTx interface {
	Commit() error
	Rollback() error
	Runnable
}
type Runnable interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}
type SpannerRunnable interface {
	QueryWithStats(context.Context, spanner.Statement) *spanner.RowIterator
}
type SpannerScanner interface {
	SpannerScan(*spanner.GenericColumnValue) error
}
type SpannerValuer interface {
	SpannerValue() (interface{}, error)
}
type SpannerScanValuer interface {
	SpannerScanner
	SpannerValuer
}
