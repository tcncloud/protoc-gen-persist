package persist_lib

func AmazingUniarySelectQuery(tx Runable, req AmazingUniarySelectQueryParams) *Result {
	row := tx.QueryRow(
		"SELECT * from example_table Where id=$1 AND start_time>$2 ",
		req.GetId(),
		req.GetStartTime(),
	)
	return newResultFromRow(row)
}
func AmazingUniarySelectWithHooksQuery(tx Runable, req AmazingUniarySelectWithHooksQueryParams) *Result {
	row := tx.QueryRow(
		"SELECT * from example_table Where id=$1 AND start_time>$2 ",
		req.GetId(),
		req.GetStartTime(),
	)
	return newResultFromRow(row)
}
func AmazingServerStreamQuery(tx Runable, req AmazingServerStreamQueryParams) *Result {
	res, err := tx.Query(
		"SELECT * FROM example_table WHERE name=$1 ",
		req.GetName(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func AmazingServerStreamWithHooksQuery(tx Runable, req AmazingServerStreamWithHooksQueryParams) *Result {
	res, err := tx.Query(
		"SELECT * FROM example_table WHERE name=$1 ",
		req.GetName(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func AmazingBidirectionalQuery(tx Runable, req AmazingBidirectionalQueryParams) *Result {
	row := tx.QueryRow(
		"UPDATE example_table SET (start_time, name) = ($2, $3) WHERE id=$1 RETURNING * ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
	return newResultFromRow(row)
}
func AmazingBidirectionalWithHooksQuery(tx Runable, req AmazingBidirectionalWithHooksQueryParams) *Result {
	row := tx.QueryRow(
		"UPDATE example_table SET (start_time, name) = ($2, $3) WHERE id=$1 RETURNING * ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
	return newResultFromRow(row)
}
func AmazingClientStreamQuery(tx Runable, req AmazingClientStreamQueryParams) *Result {
	res, err := tx.Exec(
		"INSERT INTO example_table (id, start_time, name) VALUES ($1, $2, $3) ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}
func AmazingClientStreamWithHookQuery(tx Runable, req AmazingClientStreamWithHookQueryParams) *Result {
	res, err := tx.Exec(
		"INSERT INTO example_table (id, start_time, name) VALUES ($1, $2, $3) ",
		req.GetId(),
		req.GetStartTime(),
		req.GetName(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}

type AmazingUniarySelectQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
}
type AmazingUniarySelectWithHooksQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
}
type AmazingServerStreamQueryParams interface {
	GetName() string
}
type AmazingServerStreamWithHooksQueryParams interface {
	GetName() string
}
type AmazingBidirectionalQueryParams interface {
	GetName() string
	GetId() int64
	GetStartTime() interface{}
}
type AmazingBidirectionalWithHooksQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type AmazingClientStreamQueryParams interface {
	GetStartTime() interface{}
	GetName() string
	GetId() int64
}
type AmazingClientStreamWithHookQueryParams interface {
	GetStartTime() interface{}
	GetName() string
	GetId() int64
}
