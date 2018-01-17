package persist_lib

func Testservice1UnaryExample1Query(tx Runable, req Testservice1UnaryExample1QueryParams) *Result {
	row := tx.QueryRow(
		"SELECT id AS 'table_key', id, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
		req.GetStartTime(),
	)
	return newResultFromRow(row)
}
func Testservice1UnaryExample2Query(tx Runable, req Testservice1UnaryExample2QueryParams) *Result {
	row := tx.QueryRow(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetId(),
	)
	return newResultFromRow(row)
}
func Testservice1ServerStreamSelectQuery(tx Runable, req Testservice1ServerStreamSelectQueryParams) *Result {
	res, err := tx.Query(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func Testservice1ClientStreamingExampleQuery(tx Runable, req Testservice1ClientStreamingExampleQueryParams) *Result {
	res, err := tx.Exec(
		"SELECT id AS 'table_id', key, value, msg as inner_message, status as inner_enum FROM test_table WHERE id = $1 ",
		req.GetTableId(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}

type Testservice1UnaryExample1QueryParams interface {
	GetTableId() int32
	GetStartTime() interface{}
}
type Testservice1UnaryExample2QueryParams interface {
	GetId() int32
}
type Testservice1ServerStreamSelectQueryParams interface {
	GetTableId() int32
}
type Testservice1ClientStreamingExampleQueryParams interface {
	GetTableId() int32
}
