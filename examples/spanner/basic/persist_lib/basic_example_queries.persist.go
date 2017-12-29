package persist_lib

import "cloud.google.com/go/spanner"

func NumRowsFromExtraUnaryQuery(req *Test_NumRowsForExtraSrv) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM extra_unary",
		Params: map[string]interface{}{},
	}
}
func ExampleTableFromUniaryInsertQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"name":       "bananas",
		"id":         req.Id,
		"start_time": req.StartTime,
	})
}
func ExampleTableFromUniarySelectQuery(req *Test_ExampleTableForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@id AND name=@name",
		Params: map[string]interface{}{
			"@id":   req.Id,
			"@name": req.Name,
		},
	}
}
func SomethingFromTestNestQuery(req *SomethingForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@thing",
		Params: map[string]interface{}{
			"@thing": req.Thing,
		},
	}
}
func HasTimestampFromTestEverythingQuery(req *HasTimestampForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@time AND some=@some AND str=@str AND table=@table AND times = @times AND somes = @somes AND strs = @strs AND tables = @tables",
		Params: map[string]interface{}{
			"@time":   req.Time,
			"@some":   req.Some,
			"@str":    req.Str,
			"@table":  req.Table,
			"@times":  req.Times,
			"@somes":  req.Somes,
			"@strs":   req.Strs,
			"@tables": req.Tables,
		},
	}
}
func ExampleTableFromUniarySelectWithDirectivesQuery(req *Test_ExampleTableForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table@{FORCE_INDEX=index} Where id=@id AND name=@name",
		Params: map[string]interface{}{
			"@id":   req.Id,
			"@name": req.Name,
		},
	}
}
func ExampleTableFromUniaryUpdateQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.StartTime,
		"name":       "oranges",
		"id":         req.Id,
	})
}
func ExampleTableRangeFromUniaryDeleteRangeQuery(req *Test_ExampleTableRangeForMySpanner) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.KeyRange{
		Start: spanner.Key{
			req.StartId,
		},
		End: spanner.Key{
			req.EndId,
		},
		Kind: spanner.ClosedOpen,
	})
}
func ExampleTableFromUniaryDeleteSingleQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.Key{
		"abc",
		123,
		req.Id,
	})
}
func ExampleTableFromNoArgsQuery(req *Test_ExampleTableForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL:    "select * from example_table limit 1",
		Params: map[string]interface{}{},
	}
}
func NameFromServerStreamQuery(req *Test_NameForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM example_table",
		Params: map[string]interface{}{},
	}
}
func ExampleTableFromClientStreamInsertQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"id":         req.Id,
		"start_time": req.StartTime,
		"name":       3,
	})
}
func ExampleTableFromClientStreamDeleteQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.Key{
		req.Id,
	})
}
func ExampleTableFromClientStreamUpdateQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.StartTime,
		"name":       req.Name,
		"id":         req.Id,
	})
}
func ExampleTableFromUniaryInsertWithHooksQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.InsertMap("example_table", map[string]interface{}{
		"start_time": req.StartTime,
		"name":       "bananas",
		"id":         req.Id,
	})
}
func ExampleTableFromUniarySelectWithHooksQuery(req *Test_ExampleTableForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * from example_table Where id=@id",
		Params: map[string]interface{}{
			"@id": req.Id,
		},
	}
}
func ExampleTableFromUniaryUpdateWithHooksQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"start_time": req.StartTime,
		"name":       "oranges",
		"id":         req.Id,
	})
}
func ExampleTableRangeFromUniaryDeleteWithHooksQuery(req *Test_ExampleTableRangeForMySpanner) *spanner.Mutation {
	return spanner.Delete("example_table", spanner.KeyRange{
		Start: spanner.Key{
			req.StartId,
		},
		End: spanner.Key{
			req.EndId,
		},
		Kind: spanner.ClosedOpen,
	})
}
func NameFromServerStreamWithHooksQuery(req *Test_NameForMySpanner) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * FROM example_table",
		Params: map[string]interface{}{},
	}
}
func ExampleTableFromClientStreamUpdateWithHooksQuery(req *Test_ExampleTableForMySpanner) *spanner.Mutation {
	return spanner.UpdateMap("example_table", map[string]interface{}{
		"name": "asdf",
		"id":   req.Id,
	})
}
