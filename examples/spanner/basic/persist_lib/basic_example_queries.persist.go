package persist_lib

import "cloud.google.com/go/spanner"

func NumRowsFromExtraUnaryQuery(req NumRowsFromExtraUnaryQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM extra_unary",
		Params: map[string]interface{}{},
	}
}
func ExampleTableFromUniaryInsertQuery(req ExampleTableFromUniaryInsertQueryParams) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"name":       "bananas",
		"id":         req.GetId(),
		"start_time": req.GetStartTime(),
	})
}
func ExampleTableFromUniarySelectQuery(req ExampleTableFromUniarySelectQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@id AND name=@name",
		Params: map[string]interface{}{
			"id":   req.GetId(),
			"name": req.GetName(),
		},
	}
}
func SomethingFromTestNestQuery(req SomethingFromTestNestQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@thing",
		Params: map[string]interface{}{
			"thing": req.GetThing(),
		},
	}
}
func HasTimestampFromTestEverythingQuery(req HasTimestampFromTestEverythingQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@time AND some=@some AND str=@str AND table=@table AND times = @times AND somes = @somes AND strs = @strs AND tables = @tables",
		Params: map[string]interface{}{
			"time":   req.GetTime(),
			"some":   req.GetSome(),
			"str":    req.GetStr(),
			"table":  req.GetTable(),
			"times":  req.GetTimes(),
			"somes":  req.GetSomes(),
			"strs":   req.GetStrs(),
			"tables": req.GetTables(),
		},
	}
}
func ExampleTableFromUniarySelectWithDirectivesQuery(req ExampleTableFromUniarySelectWithDirectivesQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table@{FORCE_INDEX=index} Where id=@id AND name=@name",
		Params: map[string]interface{}{
			"id":   req.GetId(),
			"name": req.GetName(),
		},
	}
}
func ExampleTableFromUniaryUpdateQuery(req ExampleTableFromUniaryUpdateQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.GetStartTime(),
		"name":       "oranges",
		"id":         req.GetId(),
	})
}
func ExampleTableRangeFromUniaryDeleteRangeQuery(req ExampleTableRangeFromUniaryDeleteRangeQueryParams) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.KeyRange{
		Start: spanner.Key{
			req.GetStartId(),
		},
		End: spanner.Key{
			req.GetEndId(),
		},
		Kind: spanner.ClosedOpen,
	})
}
func ExampleTableFromUniaryDeleteSingleQuery(req ExampleTableFromUniaryDeleteSingleQueryParams) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.Key{
		"abc",
		123,
		req.GetId(),
	})
}
func ExampleTableFromNoArgsQuery(req ExampleTableFromNoArgsQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "select * from example_table limit 1",
		Params: map[string]interface{}{},
	}
}
func NameFromServerStreamQuery(req NameFromServerStreamQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM example_table",
		Params: map[string]interface{}{},
	}
}
func ExampleTableFromClientStreamInsertQuery(req ExampleTableFromClientStreamInsertQueryParams) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"id":         req.GetId(),
		"start_time": req.GetStartTime(),
		"name":       3,
	})
}
func ExampleTableFromClientStreamDeleteQuery(req ExampleTableFromClientStreamDeleteQueryParams) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.Key{
		req.GetId(),
	})
}
func ExampleTableFromClientStreamUpdateQuery(req ExampleTableFromClientStreamUpdateQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.GetStartTime(),
		"name":       req.GetName(),
		"id":         req.GetId(),
	})
}
func ExampleTableFromUniaryInsertWithHooksQuery(req ExampleTableFromUniaryInsertWithHooksQueryParams) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"id":         req.GetId(),
		"start_time": req.GetStartTime(),
		"name":       "bananas",
	})
}
func ExampleTableFromUniarySelectWithHooksQuery(req ExampleTableFromUniarySelectWithHooksQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@id",
		Params: map[string]interface{}{
			"id": req.GetId(),
		},
	}
}
func ExampleTableFromUniaryUpdateWithHooksQuery(req ExampleTableFromUniaryUpdateWithHooksQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.GetStartTime(),
		"name":       "oranges",
		"id":         req.GetId(),
	})
}
func ExampleTableRangeFromUniaryDeleteWithHooksQuery(req ExampleTableRangeFromUniaryDeleteWithHooksQueryParams) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.KeyRange{
		Start: spanner.Key{
			req.GetStartId(),
		},
		End: spanner.Key{
			req.GetEndId(),
		},
		Kind: spanner.ClosedOpen,
	})
}
func NameFromServerStreamWithHooksQuery(req NameFromServerStreamWithHooksQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM example_table",
		Params: map[string]interface{}{},
	}
}
func ExampleTableFromClientStreamUpdateWithHooksQuery(req ExampleTableFromClientStreamUpdateWithHooksQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"name": "asdf",
		"id":   req.GetId(),
	})
}

type NumRowsFromExtraUnaryQueryParams interface {
}
type ExampleTableFromUniaryInsertQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type ExampleTableFromUniarySelectQueryParams interface {
	GetId() int64
	GetName() string
}
type SomethingFromTestNestQueryParams interface {
	GetThing() []byte
}
type HasTimestampFromTestEverythingQueryParams interface {
	GetTable() []byte
	GetTimes() [][]byte
	GetSomes() [][]byte
	GetStrs() []string
	GetTables() [][]byte
	GetTime() interface{}
	GetSome() []byte
	GetStr() string
}
type ExampleTableFromUniarySelectWithDirectivesQueryParams interface {
	GetId() int64
	GetName() string
}
type ExampleTableFromUniaryUpdateQueryParams interface {
	GetName() string
	GetId() int64
	GetStartTime() interface{}
}
type ExampleTableRangeFromUniaryDeleteRangeQueryParams interface {
	GetStartId() int64
	GetEndId() int64
}
type ExampleTableFromUniaryDeleteSingleQueryParams interface {
	GetId() int64
}
type ExampleTableFromNoArgsQueryParams interface {
}
type NameFromServerStreamQueryParams interface {
}
type ExampleTableFromClientStreamInsertQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type ExampleTableFromClientStreamDeleteQueryParams interface {
	GetId() int64
}
type ExampleTableFromClientStreamUpdateQueryParams interface {
	GetStartTime() interface{}
	GetName() string
	GetId() int64
}
type ExampleTableFromUniaryInsertWithHooksQueryParams interface {
	GetStartTime() interface{}
	GetName() string
	GetId() int64
}
type ExampleTableFromUniarySelectWithHooksQueryParams interface {
	GetId() int64
}
type ExampleTableFromUniaryUpdateWithHooksQueryParams interface {
	GetStartTime() interface{}
	GetName() string
	GetId() int64
}
type ExampleTableRangeFromUniaryDeleteWithHooksQueryParams interface {
	GetEndId() int64
	GetStartId() int64
}
type NameFromServerStreamWithHooksQueryParams interface {
}
type ExampleTableFromClientStreamUpdateWithHooksQueryParams interface {
	GetId() int64
	GetName() string
}
