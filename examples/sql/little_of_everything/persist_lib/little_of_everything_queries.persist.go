package persist_lib

func ExampleService1UnaryExample1Query(tx Runable, req ExampleService1UnaryExample1QueryParams) *Result {
	row := tx.QueryRow(
		"SELECT id AS 'table_key', id, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
		req.GetStartTime(),
	)
	return newResultFromRow(row)
}
func ExampleService1UnaryExample2Query(tx Runable, req ExampleService1UnaryExample2QueryParams) *Result {
	row := tx.QueryRow(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetId(),
	)
	return newResultFromRow(row)
}
func ExampleService1ServerStreamSelectQuery(tx Runable, req ExampleService1ServerStreamSelectQueryParams) *Result {
	res, err := tx.Query(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func ExampleService1ClientStreamingExampleQuery(tx Runable, req ExampleService1ClientStreamingExampleQueryParams) *Result {
	res, err := tx.Exec(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}

type ExampleService1UnaryExample1QueryParams interface {
	GetTableId() int32
	GetStartTime() interface{}
}
type ExampleService1UnaryExample2QueryParams interface {
	GetId() int32
}
type ExampleService1ServerStreamSelectQueryParams interface {
	GetTableId() int32
}
type ExampleService1ClientStreamingExampleQueryParams interface {
	GetTableId() int32
}
