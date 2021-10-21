package persist

import (
	"context"
	"database/sql"
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
type NotFound struct {
	Msg string
}

func (n NotFound) Error() string {
	return n.Msg
}
