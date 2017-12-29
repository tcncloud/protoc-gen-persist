package persist_lib

import "cloud.google.com/go/spanner"

func BobFromDeleteBobsQuery(req *BobForBobs) *spanner.Mutation {
	return spanner.Delete("bob_table", spanner.KeyRange{
		Start: spanner.Key{
			"Bob",
		},
		End: spanner.Key{
			"Bob",
			req.StartTime,
		},
		Kind: spanner.ClosedOpen,
	})
}
func BobFromPutBobsQuery(req *BobForBobs) *spanner.Mutation {
	return spanner.InsertMap("bob_table", map[string]interface{}{
		"name":       req.Name,
		"start_time": req.StartTime,
		"id":         req.Id,
	})
}
func EmptyFromGetBobsQuery(req *EmptyForBobs) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * from bob_table",
		Params: map[string]interface{}{},
	}
}
func NamesFromGetPeopleFromNamesQuery(req *NamesForBobs) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * FROM bob_table WHERE name IN UNNEST(@names)",
		Params: map[string]interface{}{
			"@names": req.Names,
		},
	}
}
