package persist_lib

import "database/sql"

func PartialTableFromUniarySelectQuery(tx Runable, req PartialTableFromUniarySelectQueryParams) *sql.Row {
	return tx.QueryRow(
		"SELECT * from example_table Where id=$1 AND start_time>$2 ",
		req.GetId(),
		req.GetStartTime(),
	)
}
func PartialTableFromUniarySelectWithHooksQuery(tx Runable, req PartialTableFromUniarySelectWithHooksQueryParams) *sql.Row {
	return tx.QueryRow(
		"SELECT * from example_table Where id=$1 AND start_time>$2 ",
		req.GetId(),
		req.GetStartTime(),
	)
}
func NameFromServerStreamQuery(tx Runable, req NameFromServerStreamQueryParams) (*sql.Rows, error) {
	return tx.Query(
		"SELECT * FROM example_table WHERE name=$1 ",
		req.GetName(),
	)
}
func NameFromServerStreamWithHooksQuery(tx Runable, req NameFromServerStreamWithHooksQueryParams) (*sql.Rows, error) {
	return tx.Query(
		"SELECT * FROM example_table WHERE name=$1 ",
		req.GetName(),
	)
}
func ExampleTableFromBidirectionalQuery(tx Runable, req ExampleTableFromBidirectionalQueryParams) *sql.Row {
	return tx.QueryRow(
		"UPDATE example_table SET (start_time, name) = ($2, $3) WHERE id=$1 RETURNING * ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
}
func ExampleTableFromBidirectionalWithHooksQuery(tx Runable, req ExampleTableFromBidirectionalWithHooksQueryParams) *sql.Row {
	return tx.QueryRow(
		"UPDATE example_table SET (start_time, name) = ($2, $3) WHERE id=$1 RETURNING * ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
}
func ExampleTableFromClientStreamQuery(tx Runable, req ExampleTableFromClientStreamQueryParams) (sql.Result, error) {
	return tx.Exec(
		"INSERT INTO example_table (id, start_time, name) VALUES ($1, $2, $3) ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
}
func ExampleTableFromClientStreamWithHookQuery(tx Runable, req ExampleTableFromClientStreamWithHookQueryParams) (sql.Result, error) {
	return tx.Exec(
		"INSERT INTO example_table (id, start_time, name) VALUES ($1, $2, $3) ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
}

type PartialTableFromUniarySelectQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
}
type PartialTableFromUniarySelectWithHooksQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
}
type NameFromServerStreamQueryParams interface {
	GetName() string
}
type NameFromServerStreamWithHooksQueryParams interface {
	GetName() string
}
type ExampleTableFromBidirectionalQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type ExampleTableFromBidirectionalWithHooksQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type ExampleTableFromClientStreamQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type ExampleTableFromClientStreamWithHookQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
