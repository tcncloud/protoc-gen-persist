package persist_lib

import "database/sql"

func ExampleTable1FromUnaryExample1Query(tx Runable, req ExampleTable1FromUnaryExample1QueryParams) *sql.Row {
	return tx.QueryRow(
		"SELECT id AS 'table_key', id, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
		req.GetStartTime(),
	)
}
func TestFromUnaryExample2Query(tx Runable, req TestFromUnaryExample2QueryParams) *sql.Row {
	return tx.QueryRow(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetId(),
	)
}
func ExampleTable1FromServerStreamSelectQuery(tx Runable, req ExampleTable1FromServerStreamSelectQueryParams) (*sql.Rows, error) {
	return tx.Query(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
	)
}
func ExampleTable1FromClientStreamingExampleQuery(tx Runable, req ExampleTable1FromClientStreamingExampleQueryParams) (sql.Result, error) {
	return tx.Exec(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
	)
}

type ExampleTable1FromUnaryExample1QueryParams interface {
	GetTableId() int32
	GetStartTime() interface{}
}
type TestFromUnaryExample2QueryParams interface {
	GetId() int32
}
type ExampleTable1FromServerStreamSelectQueryParams interface {
	GetTableId() int32
}
type ExampleTable1FromClientStreamingExampleQueryParams interface {
	GetTableId() int32
}
