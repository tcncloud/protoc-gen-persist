package persist_lib

import "cloud.google.com/go/spanner"

func ExtraSrvExtraUnaryQuery(req ExtraSrvExtraUnaryQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM extra_unary",
		Params: map[string]interface{}{},
	}
}
func MySpannerUniaryInsertQuery(req MySpannerUniaryInsertQueryParams) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"id":         req.GetId(),
		"start_time": req.GetStartTime(),
		"name":       "bananas",
	})
}
func MySpannerUniarySelectQuery(req MySpannerUniarySelectQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@id AND name=@name",
		Params: map[string]interface{}{
			"id":   req.GetId(),
			"name": req.GetName(),
		},
	}
}
func MySpannerTestNestQuery(req MySpannerTestNestQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@thing",
		Params: map[string]interface{}{
			"thing": req.GetThing(),
		},
	}
}
func MySpannerTestEverythingQuery(req MySpannerTestEverythingQueryParams) spanner.Statement {
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
func MySpannerUniarySelectWithDirectivesQuery(req MySpannerUniarySelectWithDirectivesQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table@{FORCE_INDEX=index} Where id=@id AND name=@name",
		Params: map[string]interface{}{
			"id":   req.GetId(),
			"name": req.GetName(),
		},
	}
}
func MySpannerUniaryUpdateQuery(req MySpannerUniaryUpdateQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.GetStartTime(),
		"name":       "oranges",
		"id":         req.GetId(),
	})
}
func MySpannerUniaryDeleteRangeQuery(req MySpannerUniaryDeleteRangeQueryParams) *spanner.Mutation {
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
func MySpannerUniaryDeleteSingleQuery(req MySpannerUniaryDeleteSingleQueryParams) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.Key{
		"abc",
		123,
		req.GetId(),
	})
}
func MySpannerNoArgsQuery(req MySpannerNoArgsQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "select * from example_table limit 1",
		Params: map[string]interface{}{},
	}
}
func MySpannerServerStreamQuery(req MySpannerServerStreamQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM example_table",
		Params: map[string]interface{}{},
	}
}
func MySpannerClientStreamInsertQuery(req MySpannerClientStreamInsertQueryParams) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"id":         req.GetId(),
		"start_time": req.GetStartTime(),
		"name":       3,
	})
}
func MySpannerClientStreamDeleteQuery(req MySpannerClientStreamDeleteQueryParams) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.Key{
		req.GetId(),
	})
}
func MySpannerClientStreamUpdateQuery(req MySpannerClientStreamUpdateQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.GetStartTime(),
		"name":       req.GetName(),
		"id":         req.GetId(),
	})
}
func MySpannerUniaryInsertWithHooksQuery(req MySpannerUniaryInsertWithHooksQueryParams) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"id":         req.GetId(),
		"start_time": req.GetStartTime(),
		"name":       "bananas",
	})
}
func MySpannerUniarySelectWithHooksQuery(req MySpannerUniarySelectWithHooksQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@id",
		Params: map[string]interface{}{
			"id": req.GetId(),
		},
	}
}
func MySpannerUniaryUpdateWithHooksQuery(req MySpannerUniaryUpdateWithHooksQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.GetStartTime(),
		"name":       "oranges",
		"id":         req.GetId(),
	})
}
func MySpannerUniaryDeleteWithHooksQuery(req MySpannerUniaryDeleteWithHooksQueryParams) *spanner.Mutation {
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
func MySpannerServerStreamWithHooksQuery(req MySpannerServerStreamWithHooksQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM example_table",
		Params: map[string]interface{}{},
	}
}
func MySpannerClientStreamUpdateWithHooksQuery(req MySpannerClientStreamUpdateWithHooksQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"name": "asdf",
		"id":   req.GetId(),
	})
}

type ExtraSrvExtraUnaryQueryParams interface {
}
type MySpannerUniaryInsertQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type MySpannerUniarySelectQueryParams interface {
	GetId() int64
	GetName() string
}
type MySpannerTestNestQueryParams interface {
	GetThing() []byte
}
type MySpannerTestEverythingQueryParams interface {
	GetStr() string
	GetTable() []byte
	GetTimes() [][]byte
	GetSomes() [][]byte
	GetStrs() []string
	GetTables() [][]byte
	GetTime() interface{}
	GetSome() []byte
}
type MySpannerUniarySelectWithDirectivesQueryParams interface {
	GetId() int64
	GetName() string
}
type MySpannerUniaryUpdateQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type MySpannerUniaryDeleteRangeQueryParams interface {
	GetStartId() int64
	GetEndId() int64
}
type MySpannerUniaryDeleteSingleQueryParams interface {
	GetId() int64
}
type MySpannerNoArgsQueryParams interface {
}
type MySpannerServerStreamQueryParams interface {
}
type MySpannerClientStreamInsertQueryParams interface {
	GetName() string
	GetId() int64
	GetStartTime() interface{}
}
type MySpannerClientStreamDeleteQueryParams interface {
	GetId() int64
}
type MySpannerClientStreamUpdateQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type MySpannerUniaryInsertWithHooksQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type MySpannerUniarySelectWithHooksQueryParams interface {
	GetId() int64
}
type MySpannerUniaryUpdateWithHooksQueryParams interface {
	GetId() int64
	GetStartTime() interface{}
	GetName() string
}
type MySpannerUniaryDeleteWithHooksQueryParams interface {
	GetEndId() int64
	GetStartId() int64
}
type MySpannerServerStreamWithHooksQueryParams interface {
}
type MySpannerClientStreamUpdateWithHooksQueryParams interface {
	GetName() string
	GetId() int64
}
